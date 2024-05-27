package blob

import (
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/blob/solidity"

	"github.com/ethereum/go-ethereum/common"
	mt "github.com/txaty/go-merkletree"
)

type dataBlock struct {
	value common.Hash
}

func (d *dataBlock) Serialize() ([]byte, error) {
	if len(d.value) == 0 {
		return nil, errors.New("datablock value is empty")
	}
	return d.value[:], nil
}

func NewMerkleTree(mode mt.TypeConfigMode, hashs ...common.Hash) (*mt.MerkleTree, error) {
	if len(hashs) == 0 {
		return nil, errors.New("hashs count is 0")
	}
	if len(hashs) == 1 {
		hashs = append(hashs, hashs[0])
	}

	dataBlocs := []mt.DataBlock{}
	for _, hash := range hashs {
		dataBlocs = append(dataBlocs, &dataBlock{
			value: hash,
		})
	}
	cfg := &mt.Config{
		Mode:               mode,
		DisableLeafHashing: true,
		SortSiblingPairs:   true,
		HashFunc:           solidity.Keccak256,
	}

	tree, err := mt.New(cfg, dataBlocs)
	if err != nil {
		return nil, err
	}

	return tree, nil
}
