package blob

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"strconv"
)

type SidecarsResponse struct {
	Data []*Sidecar `json:"data"`
}

type Sidecar struct {
	Index         string `json:"index"`
	Slot          string `json:"slot"`
	ProposerIndex string `json:"proposer_index"`
	Blob          string `json:"blob"`
	KZGCommitment string `json:"kzg_commitment"`
	KZGProof      string `json:"kzg_proof"`
}

func BuildSidecardsResponse(sidecars []*ethpb.BlobSidecar) *SidecarsResponse {
	resp := &SidecarsResponse{Data: make([]*Sidecar, len(sidecars))}
	for i, sc := range sidecars {
		resp.Data[i] = &Sidecar{
			Index:         strconv.FormatUint(sc.Index, 10),
			Blob:          hexutil.Encode(sc.Blob),
			KZGCommitment: hexutil.Encode(sc.KzgCommitment),
			KZGProof:      hexutil.Encode(sc.KzgProof),
		}
	}

	return resp
}
