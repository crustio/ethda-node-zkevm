package types

import (
	"github.com/0xPolygonHermez/zkevm-node/blob"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

type BlobTxSidecar struct {
	Blobs       []kzg4844.Blob       `json:"blobs"`       // Blobs needed by the blob pool
	Commitments []kzg4844.Commitment `json:"commitments"` // Commitments needed by the blob pool
	Proofs      []kzg4844.Proof      `json:"proofs"`      // Proofs needed by the blob pool
}

// BlobTransaction structure
type BlobTransaction struct {
	Nonce                ArgUint64       `json:"nonce"`
	GasPrice             ArgBig          `json:"gasPrice"`
	Gas                  ArgUint64       `json:"gas"`
	To                   *common.Address `json:"to"`
	Value                ArgBig          `json:"value"`
	Input                ArgBytes        `json:"input"`
	V                    ArgBig          `json:"v"`
	R                    ArgBig          `json:"r"`
	S                    ArgBig          `json:"s"`
	Hash                 common.Hash     `json:"hash"`
	From                 common.Address  `json:"from"`
	BlockHash            *common.Hash    `json:"blockHash"`
	BlockNumber          *ArgUint64      `json:"blockNumber"`
	TxIndex              *ArgUint64      `json:"transactionIndex"`
	ChainID              ArgBig          `json:"chainId"`
	Type                 ArgUint64       `json:"type"`
	Receipt              *Receipt        `json:"receipt,omitempty"`
	L2Hash               *common.Hash    `json:"l2Hash,omitempty"`
	BlobHashes           []common.Hash   `json:"blobVersionedHashes"`
	Sidecar              *BlobTxSidecar  `json:"sidecar"`
	MaxFeePerBlobGas     ArgBig          `json:"maxFeePerBlobGas"`
	MaxPriorityFeePerGas ArgBig          `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         ArgBig          `json:"maxFeePerGas"`
}

func NewBlobTransaction(tx types.Transaction,
	receipt *types.Receipt,
	includeReceipt bool,
	withBlob bool,
	l2Hash *common.Hash,
) (*Transaction, error) {
	ltx := blob.BlobTxToLegacyTx(tx)
	v, r, s := ltx.RawSignatureValues()
	from, _ := state.GetSender(tx)

	res := &Transaction{
		Nonce:            ArgUint64(tx.Nonce()),
		GasPrice:         ArgBig(*tx.GasPrice()),
		Gas:              ArgUint64(tx.Gas()),
		To:               tx.To(),
		Value:            ArgBig(*tx.Value()),
		Input:            tx.Data(),
		V:                ArgBig(*v),
		R:                ArgBig(*r),
		S:                ArgBig(*s),
		Hash:             ltx.Hash(),
		From:             from,
		ChainID:          ArgBig(*tx.ChainId()),
		Type:             ArgUint64(tx.Type()),
		BlobHashes:       tx.BlobHashes(),
		MaxFeePerBlobGas: ArgBig(*tx.BlobGasFeeCap()),
	}

	if withBlob {
		res.Sidecar = &BlobTxSidecar{
			Blobs:       tx.BlobTxSidecar().Blobs,
			Commitments: tx.BlobTxSidecar().Commitments,
			Proofs:      tx.BlobTxSidecar().Proofs,
		}
	}

	if receipt != nil {
		bn := ArgUint64(receipt.BlockNumber.Uint64())
		res.BlockNumber = &bn
		res.BlockHash = &receipt.BlockHash
		ti := ArgUint64(receipt.TransactionIndex)
		res.TxIndex = &ti
		rpcReceipt, err := NewReceipt(*blob.BlobTxToLegacyTx(tx), receipt, l2Hash)
		if err != nil {
			return nil, err
		}
		if includeReceipt {
			res.Receipt = &rpcReceipt
		}
	}

	return res, nil
}
