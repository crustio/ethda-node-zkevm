package jsonrpc

import (
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"os"
	"path/filepath"
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

func (e *ETHDAEndpoints) SignBatchHash(hash types.ArgHash) (interface{}, types.Error) {
	log.Infof("Sign batch hash: %v", hash.Hash().Hex())
	sig, err := e.sk.Sign(rand.Reader, hash[:], nil)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to sign batch hash", err, true)
	}

	log.Infof("Sign batch hash result: %v", sig)

	return common.Bytes2Hex(sig), nil
}
