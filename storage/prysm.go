package storage

import (
	"bytes"
	"math"

	"github.com/attestantio/go-eth2-client/spec/deneb"
	fieldparams "github.com/prysmaticlabs/prysm/v5/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"github.com/rs/zerolog"
)

func NewPrysmBlobStorage(log zerolog.Logger, path string) (*PrysmBlobStorage, error) {
	blobStorage, err := NewBlobStorage(
		WithLogger(log),
		WithBasePath(path),
		WithBlobRetentionEpochs(math.MaxUint64),
		WithSaveFsync(true),
	)
	if err != nil {
		return nil, err
	}
	return &PrysmBlobStorage{blobStorage: blobStorage}, nil
}

var _ BlobStore = &PrysmBlobStorage{}

type PrysmBlobStorage struct {
	blobStorage *BlobStorage
}

func (p *PrysmBlobStorage) Exist(root [32]byte) bool {
	data, err := p.blobStorage.Indices(root)
	if err != nil {
		return false
	}
	for _, d := range data {
		if d == true {
			return true
		}
	}
	return false
}

func (p *PrysmBlobStorage) Save(root [32]byte, denebSidecar *deneb.BlobSidecar) error {
	sidecar := ConvSideCar(denebSidecar)
	return p.blobStorage.Save(root, sidecar)
}

func (p *PrysmBlobStorage) Get(root [32]byte, index uint64) (*ethpb.BlobSidecar, error) {
	blob, err := p.blobStorage.Get(root, index)
	if err != nil {
		return nil, err
	}

	return blob, nil
}

func (p *PrysmBlobStorage) Valid(root [32]byte, denebSidecar *deneb.BlobSidecar) (bool, error) {
	sidecar1 := ConvSideCar(denebSidecar)
	marshal1, err := sidecar1.MarshalSSZ()
	if err != nil {
		return false, err
	}
	sidecar2, err := p.Get(root, sidecar1.Index)
	if err != nil {
		return false, err
	}
	marshal2, err := sidecar2.MarshalSSZ()
	if err != nil {
		return false, err
	}
	return bytes.Equal(marshal1, marshal2), nil
}

func ConvSideCar(denebSidecar *deneb.BlobSidecar) *ethpb.BlobSidecar {
	var sidecar *ethpb.BlobSidecar = HydrateBlobSidecar(nil)
	sidecar.Index = uint64(denebSidecar.Index)
	sidecar.Blob = denebSidecar.Blob[:]
	sidecar.KzgCommitment = denebSidecar.KZGCommitment[:]
	sidecar.KzgProof = denebSidecar.KZGProof[:]
	for i, proof := range denebSidecar.KZGCommitmentInclusionProof {
		sidecar.CommitmentInclusionProof[i] = proof[:]
	}
	sidecar.SignedBlockHeader.Signature = denebSidecar.SignedBlockHeader.Signature[:]
	sidecar.SignedBlockHeader.Header.Slot = primitives.Slot(denebSidecar.SignedBlockHeader.Message.Slot)
	sidecar.SignedBlockHeader.Header.ProposerIndex = primitives.ValidatorIndex(denebSidecar.SignedBlockHeader.Message.ProposerIndex)
	sidecar.SignedBlockHeader.Header.ParentRoot = denebSidecar.SignedBlockHeader.Message.ParentRoot[:]
	sidecar.SignedBlockHeader.Header.BodyRoot = denebSidecar.SignedBlockHeader.Message.BodyRoot[:]
	sidecar.SignedBlockHeader.Header.StateRoot = denebSidecar.SignedBlockHeader.Message.StateRoot[:]
	return sidecar
}

// HydrateBlobSidecar hydrates a blob sidecar with correct field length sizes
// to comply with SSZ marshalling and unmarshalling rules.
func HydrateBlobSidecar(b *ethpb.BlobSidecar) *ethpb.BlobSidecar {
	if b == nil {
		b = &ethpb.BlobSidecar{}
	}
	if b.SignedBlockHeader == nil {
		b.SignedBlockHeader = HydrateSignedBeaconHeader(&ethpb.SignedBeaconBlockHeader{
			Header: &ethpb.BeaconBlockHeader{},
		})
	}
	if b.Blob == nil {
		b.Blob = make([]byte, fieldparams.BlobLength)
	}
	if b.KzgCommitment == nil {
		b.KzgCommitment = make([]byte, fieldparams.BLSPubkeyLength)
	}
	if b.KzgProof == nil {
		b.KzgProof = make([]byte, fieldparams.BLSPubkeyLength)
	}

	if b.CommitmentInclusionProof == nil {
		b.CommitmentInclusionProof = HydrateCommitmentInclusionProofs()
	}
	return b
}

// HydrateCommitmentInclusionProofs returns 2d byte slice of Commitment Inclusion Proofs
func HydrateCommitmentInclusionProofs() [][]byte {
	r := make([][]byte, fieldparams.KzgCommitmentInclusionProofDepth)
	for i := range r {
		r[i] = make([]byte, fieldparams.RootLength)
	}
	return r
}

// HydrateSignedBeaconHeader hydrates a signed beacon block header with correct field length sizes
// to comply with fssz marshalling and unmarshalling rules.
func HydrateSignedBeaconHeader(h *ethpb.SignedBeaconBlockHeader) *ethpb.SignedBeaconBlockHeader {
	if h.Signature == nil {
		h.Signature = make([]byte, fieldparams.BLSSignatureLength)
	}
	h.Header = HydrateBeaconHeader(h.Header)
	return h
}

// HydrateBeaconHeader hydrates a beacon block header with correct field length sizes
// to comply with fssz marshalling and unmarshalling rules.
func HydrateBeaconHeader(h *ethpb.BeaconBlockHeader) *ethpb.BeaconBlockHeader {
	if h == nil {
		h = &ethpb.BeaconBlockHeader{}
	}
	if h.BodyRoot == nil {
		h.BodyRoot = make([]byte, fieldparams.RootLength)
	}
	if h.StateRoot == nil {
		h.StateRoot = make([]byte, fieldparams.RootLength)
	}
	if h.ParentRoot == nil {
		h.ParentRoot = make([]byte, fieldparams.RootLength)
	}
	return h
}
