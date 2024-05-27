package solidity

import (
	smt "github.com/FantasyJony/openzeppelin-merkle-tree-go/standard_merkle_tree"
)

func AbiPack(types []string, values ...interface{}) ([]byte, error) {
	return smt.AbiPack(types, values...)
}
