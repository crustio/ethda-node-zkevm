package jsonrpc

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/blob"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	APIETHDA = "ethda"
)

// ETHDAEndpoints contains implementations for the "ethda" RPC endpoints
type ETHDAEndpoints struct {
	sk *ecdsa.PrivateKey
}

// NewETHDAEndpoints returns ETHDAEndpoints
func NewETHDAEndpoints(skPath, skPassword string) *ETHDAEndpoints {
	keystoreEncrypted, err := os.ReadFile(filepath.Clean(skPath))
	if err != nil {
		panic(err)
	}
	key, err := keystore.DecryptKey(keystoreEncrypted, skPassword)
	if err != nil {
		panic(err)
	}

	return &ETHDAEndpoints{
		sk: key.PrivateKey,
	}
}

type ProofResponse struct {
	BatchNumber uint64
	Proof       []common.Hash
}

func (e *ETHDAEndpoints) GetProofByHash(hash types.ArgHash) (interface{}, types.Error) {
	batchNumber := 0 // TODO

	hashes := []common.Hash{} // TODO

	tree, err := blob.NewMerkleTree(blob.ModeProofGen, hashes...)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed create merkle tree by hashes", err, true)
	}

	res := &ProofResponse{
		BatchNumber: uint64(batchNumber),
		Proof:       []common.Hash{},
	}

	hashIndex := -1
	leafs := tree.Leaves
	for i := 0; i < len(leafs); i++ {
		if hash.Hash() == common.BytesToHash(leafs[i]) {
			hashIndex = i
		}
	}

	if hashIndex == -1 { // hash not found
		errStr := "hash not found: " + hash.Hash().Hex()
		return RPCErrorResponse(types.DefaultErrorCode, errStr, errors.New(errStr), true)
	}

	for _, h := range tree.Proofs[hashIndex].Siblings {
		res.Proof = append(res.Proof, common.BytesToHash(h))
	}

	return res, nil
}

func (e *ETHDAEndpoints) SignBatchHash(hash types.ArgHash) (interface{}, types.Error) {
	log.Infof("Sign batch hash: %v", hash.Hash().Hex())
	sig, err := crypto.Sign(hash[:], e.sk)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to sign batch hash", err, true)
	}

	rBytes := sig[:32]
	sBytes := sig[32:64]
	vByte := sig[64]

	if strings.ToUpper(common.Bytes2Hex(sBytes)) > "7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF5D576E7357A4501DDFE92F46681B20A0" {
		magicNumber := common.Hex2Bytes("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141")
		sBig := big.NewInt(0).SetBytes(sBytes)
		magicBig := big.NewInt(0).SetBytes(magicNumber)
		s1 := magicBig.Sub(magicBig, sBig)
		sBytes = s1.Bytes()
		if vByte == 0 {
			vByte = 1
		} else {
			vByte = 0
		}
	}
	vByte += 27

	actualSignature := []byte{}
	actualSignature = append(actualSignature, rBytes...)
	actualSignature = append(actualSignature, sBytes...)
	actualSignature = append(actualSignature, vByte)

	log.Infof("Sign batch hash result: %v", sig)

	return common.Bytes2Hex(actualSignature), nil
}
