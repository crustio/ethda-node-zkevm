package pool

import (
	"context"
	"fmt"
	"github.com/0xPolygonHermez/zkevm-node/blob"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

func GetBlobTxSender(tx types.Transaction) (common.Address, error) {
	var signer types.Signer
	signer = types.NewCancunSigner(tx.ChainId())

	legacyTx := blob.BlobTxToLegacyTx(tx)

	return types.Sender(signer, legacyTx)
}

func (p *Pool) validateBlobTx(ctx context.Context, tx types.Transaction) error {
	// gets tx sender for validations
	from, err := GetBlobTxSender(tx)
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

	if tx.Value().Cmp(new(big.Int).Sub(tx.Cost(), tx.Value())) < 0 {
		return fmt.Errorf("value is less than blob cost, %s < %s", tx.Value().String(), tx.Cost())
	}

	if tx.BlobGasFeeCapIntCmp(big.NewInt(1000000000)) == -1 {
		return fmt.Errorf("blob gas fee is less than base gas fee")
	}

	return nil
}

func (p *Pool) GetBlobTx(ctx context.Context, hash common.Hash) ([]byte, error) {
	return p.BlobDB.Get([]byte(fmt.Sprintf("blob-%s", hash.Hex())), nil)
}
