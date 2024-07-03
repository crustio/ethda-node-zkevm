package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/0xPolygonHermez/zkevm-node"
	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/event/nileventstorage"
	"github.com/0xPolygonHermez/zkevm-node/event/pgeventstorage"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/pgstatestorage"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/urfave/cli/v2"
)

const appName = "zkblob-sender"

var (
	configFileFlag = cli.StringFlag{
		Name:     config.FlagCfg,
		Aliases:  []string{"c"},
		Usage:    "Configuration `FILE`",
		Required: true,
	}
	networkFlag = cli.StringFlag{
		Name:     config.FlagNetwork,
		Aliases:  []string{"net"},
		Usage:    "Load default network configuration. Supported values: [`mainnet`, `testnet`, `cardona`, `custom`]",
		Required: true,
	}
	customNetworkFlag = cli.StringFlag{
		Name:     config.FlagCustomNetwork,
		Aliases:  []string{"net-file"},
		Usage:    "Load the network configuration file if --network=custom",
		Required: false,
	}
)

func main() {

	app := cli.NewApp()
	app.Name = appName
	app.Version = zkevm.Version
	flags := []cli.Flag{
		&configFileFlag,
	}
	app.Commands = []*cli.Command{
		{
			Name:    "run",
			Aliases: []string{},
			Usage:   "Run the zkblob-sender",
			Action:  start,
			Flags:   append(flags, &networkFlag, &customNetworkFlag),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}

func start(cliCtx *cli.Context) error {
	c, err := config.Load(cliCtx, true)
	if err != nil {
		return err
	}
	setupLog(c.Log)

	log.Info("Starting application")

	var (
		eventLog     *event.EventLog
		eventStorage event.Storage
	)

	ev := &event.Event{
		ReceivedAt: time.Now(),
		Source:     event.Source_Node,
		Level:      event.Level_Info,
		EventID:    event.EventID_NodeComponentStarted,
	}

	if c.EventLog.DB.Name != "" {
		eventStorage, err = pgeventstorage.NewPostgresEventStorage(c.EventLog.DB)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		eventStorage, err = nileventstorage.NewNilEventStorage()
		if err != nil {
			log.Fatal(err)
		}
	}
	// event log
	eventLog = event.NewEventLog(c.EventLog, eventStorage)

	etherman, err := newEtherman(*c)
	if err != nil {
		log.Fatal(err)
	}

	// READ CHAIN ID FROM POE SC
	l2ChainID, err := etherman.GetL2ChainID()
	if err != nil {
		log.Fatal(err)
	}

	// Core State DB
	stateSqlDB, err := db.NewSQLDB(c.State.DB)
	if err != nil {
		log.Fatal(err)
	}

	// st
	needsExecutor, needsStateTree := false, false
	st, _ := newState(cliCtx.Context, c, etherman, l2ChainID, stateSqlDB, eventLog, needsExecutor, needsStateTree, false)

	ethTxManagerStorage, err := ethtxmanager.NewPostgresStorage(c.State.DB)
	if err != nil {
		log.Fatal(err)
	}

	ev.Component = event.Component_Sequence_Sender
	ev.Description = "Running zkblob sender"
	err = eventLog.LogEvent(cliCtx.Context, ev)
	if err != nil {
		log.Fatal(err)
	}

	ethTxManager := ethtxmanager.New(c.EthTxManager, etherman, ethTxManagerStorage, st)

	// zkblob client
	zkBlobClient, err := NewClient(etherman, c.NetworkConfig.L1Config.L1ChainID, c.Etherman.URL, c.Blob, c.Blob.PrivateKey.Path, c.Blob.PrivateKey.Password)
	if err != nil {
		log.Fatal(err)
	}

	// zkblob service
	seqSender, err := New(c.Blob, st, etherman, ethTxManager, eventLog, zkBlobClient)
	if err != nil {
		log.Fatal(err)
	}

	seqSender.Start(cliCtx.Context)

	return nil
}

func newState(ctx context.Context, c *config.Config, etherman *etherman.Client, l2ChainID uint64, sqlDB *pgxpool.Pool, eventLog *event.EventLog, needsExecutor, needsStateTree, avoidForkIDInMemory bool) (*state.State, uint64) {
	// Executor
	var executorClient executor.ExecutorServiceClient
	if needsExecutor {
		executorClient, _, _ = executor.NewExecutorClient(ctx, c.Executor)
	}

	// State Tree
	var stateTree *merkletree.StateTree
	if needsStateTree {
		stateDBClient, _, _ := merkletree.NewMTDBServiceClient(ctx, c.MTClient)
		stateTree = merkletree.NewStateTree(stateDBClient)
	}

	stateCfg := state.Config{
		MaxCumulativeGasUsed:         c.State.Batch.Constraints.MaxCumulativeGasUsed,
		ChainID:                      l2ChainID,
		ForkIDIntervals:              []state.ForkIDInterval{},
		MaxResourceExhaustedAttempts: c.Executor.MaxResourceExhaustedAttempts,
		WaitOnResourceExhaustion:     c.Executor.WaitOnResourceExhaustion,
		ForkUpgradeBatchNumber:       c.ForkUpgradeBatchNumber,
		ForkUpgradeNewForkId:         c.ForkUpgradeNewForkId,
		MaxLogsCount:                 c.RPC.MaxLogsCount,
		MaxLogsBlockRange:            c.RPC.MaxLogsBlockRange,
		MaxNativeBlockHashBlockRange: c.RPC.MaxNativeBlockHashBlockRange,
		AvoidForkIDInMemory:          avoidForkIDInMemory,
	}
	stateDb := pgstatestorage.NewPostgresStorage(stateCfg, sqlDB)
	st := state.NewState(stateCfg, stateDb, executorClient, stateTree, eventLog, nil)
	// This is to force to build cache, and check that DB is ok before starting the application
	l1inforoot, err := st.GetCurrentL1InfoRoot(ctx, nil)
	if err != nil {
		log.Fatal("error getting current L1InfoRoot. Error: ", err)
	}
	log.Infof("Starting L1InfoRoot: %v", l1inforoot.String())

	forkIDIntervals, err := forkIDIntervals(ctx, st, etherman, c.NetworkConfig.Genesis.BlockNumber)
	if err != nil {
		log.Fatal("error getting forkIDs. Error: ", err)
	}
	st.UpdateForkIDIntervalsInMemory(forkIDIntervals)

	currentForkID := forkIDIntervals[len(forkIDIntervals)-1].ForkId
	log.Infof("Fork ID read from POE SC = %v", forkIDIntervals[len(forkIDIntervals)-1].ForkId)

	return st, currentForkID
}
func forkIDIntervals(ctx context.Context, st *state.State, etherman *etherman.Client, genesisBlockNumber uint64) ([]state.ForkIDInterval, error) {
	log.Debug("getting forkIDs from db")
	forkIDIntervals, err := st.GetForkIDs(ctx, nil)
	if err != nil && !errors.Is(err, state.ErrStateNotSynchronized) {
		return []state.ForkIDInterval{}, fmt.Errorf("error getting forkIDs from db. Error: %v", err)
	}
	numberForkIDs := len(forkIDIntervals)
	log.Debug("numberForkIDs: ", numberForkIDs)
	// var forkIDIntervals []state.ForkIDInterval
	if numberForkIDs == 0 {
		// Get last L1block Synced
		lastBlock, err := st.GetLastBlock(ctx, nil)
		if err != nil && !errors.Is(err, state.ErrStateNotSynchronized) {
			return []state.ForkIDInterval{}, fmt.Errorf("error checking lastL1BlockSynced. Error: %v", err)
		}
		if lastBlock != nil {
			log.Info("Getting forkIDs intervals. Please wait...")
			// Read Fork ID FROM POE SC
			forkIntervals, err := etherman.GetForks(ctx, genesisBlockNumber, lastBlock.BlockNumber)
			if err != nil {
				return []state.ForkIDInterval{}, fmt.Errorf("error getting forks. Please check the configuration. Error: %v", err)
			} else if len(forkIntervals) == 0 {
				return []state.ForkIDInterval{}, fmt.Errorf("error: no forkID received. It should receive at least one, please check the configuration")
			}

			dbTx, err := st.BeginStateTransaction(ctx)
			if err != nil {
				return []state.ForkIDInterval{}, fmt.Errorf("error creating dbTx. Error: %v", err)
			}
			log.Info("Storing forkID intervals into db")
			// Store forkIDs
			for _, f := range forkIntervals {
				err := st.AddForkID(ctx, f, dbTx)
				if err != nil {
					log.Errorf("error adding forkID to db. Error: %v", err)
					rollbackErr := dbTx.Rollback(ctx)
					if rollbackErr != nil {
						log.Errorf("error rolling back dbTx. RollbackErr: %s. Error : %v", rollbackErr.Error(), err)
						return []state.ForkIDInterval{}, rollbackErr
					}
					return []state.ForkIDInterval{}, fmt.Errorf("error adding forkID to db. Error: %v", err)
				}
			}
			err = dbTx.Commit(ctx)
			if err != nil {
				log.Errorf("error committing dbTx. Error: %v", err)
				rollbackErr := dbTx.Rollback(ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back dbTx. RollbackErr: %s. Error : %v", rollbackErr.Error(), err)
					return []state.ForkIDInterval{}, rollbackErr
				}
				return []state.ForkIDInterval{}, fmt.Errorf("error committing dbTx. Error: %v", err)
			}
			forkIDIntervals = forkIntervals
		} else {
			log.Debug("Getting initial forkID")
			forkIntervals, err := etherman.GetForks(ctx, genesisBlockNumber, genesisBlockNumber)
			if err != nil {
				return []state.ForkIDInterval{}, fmt.Errorf("error getting forks. Please check the configuration. Error: %v", err)
			} else if len(forkIntervals) == 0 {
				return []state.ForkIDInterval{}, fmt.Errorf("error: no forkID received. It should receive at least one, please check the configuration")
			}
			forkIDIntervals = forkIntervals
		}
	}
	return forkIDIntervals, nil
}

func newEtherman(c config.Config) (*etherman.Client, error) {
	return etherman.NewClient(c.Etherman, c.NetworkConfig.L1Config)
}

func setupLog(c log.Config) {
	log.Init(c)
}
