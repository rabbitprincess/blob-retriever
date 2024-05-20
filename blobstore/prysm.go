package blobstore

import (
	"math"

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/db/filesystem"
	fieldparams "github.com/prysmaticlabs/prysm/v5/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
)

func NewPrysmBlobStorage(path string, start, end uint64) *PrysmBlobStorage {
	blobStorage, err := filesystem.NewBlobStorage(
		filesystem.WithBasePath(path),
		filesystem.WithBlobRetentionEpochs(math.MaxUint64),
		filesystem.WithSaveFsync(true),
	)
	if err != nil {
		return nil
	}
	return &PrysmBlobStorage{blobStorage: blobStorage}
}

type PrysmBlobStorage struct {
	blobStorage *filesystem.BlobStorage
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
	ROBlob, err := blocks.NewROBlobWithRoot(sidecar, root)
	if err != nil {
		return err
	}
	return p.blobStorage.Save(blocks.NewVerifiedROBlob(ROBlob))
}

func ConvSideCar(denebSidecar *deneb.BlobSidecar) *ethpb.BlobSidecar {
	var sidecar *ethpb.BlobSidecar = HydrateBlobSidecar(nil)
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
