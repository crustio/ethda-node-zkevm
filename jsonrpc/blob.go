package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/0xPolygonHermez/zkevm-node/blob"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	types2 "github.com/ethereum/go-ethereum/core/types"
	"github.com/gorilla/mux"
	"github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
	"strconv"
)

func (s *Server) HandleGetBlocks(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	blockId := vars["block_id"]
	var height uint64
	var err error

	switch blockId {
	case "head":
		height, err = s.st.GetLastL2BlockNumber(context.Background(), nil)
		if err != nil {
			handleError(w, fmt.Errorf("could not decode block ID: %s", blockId))
			return
		}
	case "finalized":
		syncInfo, err := s.st.GetSyncingInfo(context.Background(), nil)
		if err != nil {
			handleError(w, fmt.Errorf("could not decode block ID: %s", blockId))
			return
		}
		height = syncInfo.CurrentBlockNumber
	default:
		if bytesutil.IsHex([]byte(blockId)) {
			height, err = strconv.ParseUint(blockId, 16, 64)
			if err != nil {
				handleError(w, fmt.Errorf("could not decode block ID into hex: %s", blockId))
				return
			}
		} else {
			height, err = strconv.ParseUint(blockId, 10, 64)
			if err != nil {
				handleError(w, fmt.Errorf("could not decode block ID: %s", blockId))
				return
			}
		}
	}
	block, err := s.st.GetL2BlockByNumber(context.Background(), height, nil)
	if err != nil {
		handleError(w, fmt.Errorf("could not retrieve block for height: %s,%d", blockId, height))
		return
	}

	blobH := ""
	for i := 0; i < len(block.Transactions()); i++ {
		tx := block.Transactions()[i]

		existed, err := s.p.IsBlob(context.Background(), tx.Hash())
		if err != nil {
			handleError(w, fmt.Errorf("failed to check blob tx from leveldb: %s=>%d", blockId, height))
			return
		}
		if existed {
			data, err := s.p.GetBlobTx(context.Background(), tx.Hash())
			if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
				handleError(w, fmt.Errorf("failed to load blob tx from leveldb", blockId, height))
				return
			}
			newTx := new(types2.Transaction)
			if err := newTx.UnmarshalBinary(data); err != nil {
				handleError(w, fmt.Errorf("failed to unmarshal blob data"))
				return
			}
			blobH = common.Bytes2Hex(newTx.BlobTxSidecar().Commitments[0][:])
		}
	}

	var blobCs []string
	if blobH != "" {
		blobCs = append(blobCs, blobH)
	}

	res := map[string]any{
		"version": "capella",
		"data": map[string]any{
			"message": map[string]any{
				"slot": block.Header().Number.String(),
				"body": map[string]any{
					"execution_payload": map[string]any{
						"block_hash": block.Hash().Hex(),
					},
					"blob_kzg_commitments": blobCs,
				},
			},
		},
	}

	respBytes, _ := json.Marshal(res)
	_, err = w.Write(respBytes)
	if err != nil {
		log.Error(err)
		return
	}
}

func (s *Server) HandleGetBlobSidecars(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	blockId := vars["block_id"]

	var height uint64
	var err error
	if bytesutil.IsHex([]byte(blockId)) {
		height, err = strconv.ParseUint(blockId, 16, 64)
		if err != nil {
			handleError(w, fmt.Errorf("could not decode block ID into hex: %s", blockId))
			return
		}
	} else {
		height, err = strconv.ParseUint(blockId, 10, 64)
		if err != nil {
			handleError(w, fmt.Errorf("could not decode block ID: %s", blockId))
			return
		}
	}
	block, err := s.st.GetL2BlockByNumber(context.Background(), height, nil)
	if err != nil {
		handleError(w, fmt.Errorf("could not retrieve blobs for height: %s", blockId))
		return
	}

	var bss []*ethpb.BlobSidecar
	txs := block.Transactions()
	for _, tx := range txs {
		existed, err := s.p.IsBlob(context.Background(), tx.Hash())
		if err != nil {
			handleError(w, fmt.Errorf("failed to check blob tx from leveldb: %s=>%d", blockId, height))
			return
		}
		if !existed {
			continue
		}
		data, err := s.p.GetBlobTx(context.Background(), tx.Hash())
		if err != nil {
			handleError(w, fmt.Errorf("failed to load blob tx from leveldb: %s=>%d", blockId, height))
			return
		}

		newTx := new(types2.Transaction)
		if err := newTx.UnmarshalBinary(data); err != nil {
			handleError(w, fmt.Errorf("could not unmarshal blob tx: %s", tx.Hash().Hex()))
			return
		}

		for i, _ := range newTx.BlobTxSidecar().Blobs {
			bsc := &ethpb.BlobSidecar{
				Index:         uint64(i),
				Blob:          newTx.BlobTxSidecar().Blobs[i][:],
				KzgCommitment: newTx.BlobTxSidecar().Commitments[i][:],
				KzgProof:      newTx.BlobTxSidecar().Proofs[i][:],
			}

			bss = append(bss, bsc)
		}
	}

	res := blob.BuildSidecardsResponse(bss)

	respBytes, _ := json.Marshal(res)
	_, err = w.Write(respBytes)
	if err != nil {
		log.Error(err)
		return
	}
}
