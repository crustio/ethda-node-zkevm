package zkblob

import (
	"context"
	"math/big"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

// etherman contains the methods required to interact with ethereum.
type etherman interface {
	BuildPostZkBlobTxData(sender common.Address, batchNumber int64, hashes []common.Hash) (to *common.Address, data []byte, err error)
	BuildSequenceBatchesTxData(sender common.Address, sequences []ethmanTypes.Sequence, maxSequenceTimestamp uint64, initSequenceBatchNumber uint64, l2Coinbase common.Address) (to *common.Address, data []byte, err error)
	EstimateGasSequenceBatches(sender common.Address, sequences []ethmanTypes.Sequence, maxSequenceTimestamp uint64, initSequenceBatchNumber uint64, l2Coinbase common.Address) (*types.Transaction, error)
	GetLatestBlockHeader(ctx context.Context) (*types.Header, error)
	GetLatestBatchNumber() (uint64, error)
}

type ethTxManager interface {
	Add(ctx context.Context, owner, id string, from common.Address, to *common.Address, value *big.Int, data []byte, gasOffset uint64, dbTx pgx.Tx) error
	ProcessPendingMonitoredTxs(ctx context.Context, owner string, failedResultHandler ethtxmanager.ResultHandler, dbTx pgx.Tx)
}
