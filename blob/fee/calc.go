package fee

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/params"
)

const (
	MinBlobBaseFee = 1000000000
)

func CalcBlobGasUsed(ctx context.Context, s stateInterface, pool txPool, batchNumber uint64) uint64 {
	batch, err := s.GetBatchByNumber(ctx, batchNumber, nil)
	if err != nil {
		return 0 // TODO handle error
	}
	brb, er := state.DecodeBatchV2(batch.BatchL2Data)
	if er != nil {
		return 0 // TODO handle error
	}

	var count uint64 = 0
	for _, btx := range brb.Blocks {
		for _, tx := range btx.Transactions {
			if ok, _ := pool.IsBlob(ctx, tx.Tx.Hash()); ok {
				count++
			}
		}
	}

	return uint64(count * params.BlobTxBlobGasPerBlob)
}
