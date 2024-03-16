package pool

import (
	"context"
	"fmt"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

func (p *Pool) validateBlobTx(ctx context.Context, tx types.Transaction) error {
	// gets tx sender for validations
	from, err := state.GetSender(tx)
	if err != nil {
		return ErrInvalidSender
	}

	log.Infof("from: ", from.Hex())

	lastL2Block, err := p.state.GetLastL2Block(ctx, nil)
	if err != nil {
		return err
	}

	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	balance, err := p.state.GetBalance(ctx, from, lastL2Block.Root())
	if err != nil {
		return err
	}

	if balance.Cmp(tx.Cost()) < 0 {
		return ErrInsufficientFunds
	}

	if tx.Value().Cmp(tx.Cost()) < 0 {
		return fmt.Errorf("value is less than blob cost")
	}

	if tx.BlobGasFeeCapIntCmp(big.NewInt(1000000000)) == -1 {
		return fmt.Errorf("blob gas fee is less than base gas fee")
	}

	return nil
}
