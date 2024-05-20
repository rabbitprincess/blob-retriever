package restore

import "github.com/attestantio/go-eth2-client/spec/deneb"

type BlobStore interface {
	Exist(root [32]byte) bool
	Save(root [32]byte, denebSidecar *deneb.BlobSidecar) error
}
