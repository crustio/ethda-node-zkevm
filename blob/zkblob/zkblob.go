package zkblob

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

const (
	monitoredBatchIDFormat = "batch-%v"
	disableZkBlob          = false
)

func HandleSequenceToZkBlob(ctx context.Context, sequences []types.Sequence, ethTxManagerOwner string, senderAddress common.Address, etherman etherman, ethTxManager ethTxManager, GasOffset uint64) {
	if disableZkBlob {
		return
	}

	for _, seq := range sequences {
		hashes := []common.Hash{}
		brb, err := state.DecodeBatchV2(seq.BatchL2Data)
		if err != nil {
			log.Errorf("error decode BatchL2Data for batch %d, err: %v", seq.BatchNumber, err)
		}

		for _, btx := range brb.Blocks {
			for _, tx := range btx.Transactions {
				hashes = append(hashes, tx.Tx.Hash())
			}
		}

		if len(hashes) == 0 {
			continue
		}

		to, data, err := etherman.BuildPostZkBlobTxData(senderAddress, int64(seq.BatchNumber), hashes)
		if err != nil {
			log.Errorf("error build postZkBlob tx data, batch number: %d: %v", seq.BatchNumber, err)
			continue
		}

		monitoredTxID := fmt.Sprintf(monitoredBatchIDFormat, seq.BatchNumber)
		err = ethTxManager.Add(ctx, ethTxManagerOwner, monitoredTxID, senderAddress, to, nil, data, GasOffset, nil)
		if err != nil {
			mTxLogger := ethtxmanager.CreateLogger(ethTxManagerOwner, monitoredTxID, senderAddress, to)
			mTxLogger.Errorf("error to add postZkBlob tx to eth tx manager: ", err)
			continue
		}
	}
}
