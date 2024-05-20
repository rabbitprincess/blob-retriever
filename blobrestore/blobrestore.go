package blobrestore

import (
	"context"
	"strconv"
	"time"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/avast/retry-go"
	"github.com/rs/zerolog"
)

const (
	serverTimeout = 60 * time.Second
	minBlobSlot   = 8626176 // mainnet minimum slot which save blobs
)

type BlobRestore struct {
	logger  zerolog.Logger
	client  BeaconClient
	storage BlobStore
}

// NewAPI creates a new Archiver API instance. This API exposes an admin interface to control the archiver.
func NewAPI() *BlobRestore {
	result := &BlobRestore{
		logger: zerolog.New(nil),
	}

	return result
}

func (bs *BlobRestore) RestoreBlob(ctx context.Context, fromSlot, toSlot uint64) error {
	if fromSlot < minBlobSlot {
		fromSlot = minBlobSlot
	}
	if toSlot < fromSlot {
		toSlot = fromSlot
	}

	for slot := fromSlot; slot <= toSlot; slot++ {
		var root phase0.Root
		if err := retry.Do(func() error {
			res, err := bs.client.BeaconBlockHeader(ctx, &api.BeaconBlockHeaderOpts{
				Block: strconv.FormatUint(slot, 10),
			})
			if err != nil {
				bs.logger.Error().Uint64("slot", slot).Err(err).Msg("Failed to get block root | retrying...")
				return err
			}
			root = res.Data.Root
			return nil
		}, retry.Attempts(5), retry.Delay(1*time.Second)); err != nil {
			bs.logger.Error().Uint64("slot", slot).Err(err).Msg("Failed to get block header")
		}

		if bs.storage.Exist(root) {
			bs.logger.Info().Uint64("slot", slot).Str("root", root.String()).Msg("Blob already exists, continue")
			continue
		}

		if err := retry.Do(func() error {
			blobSideCars, err := bs.client.BlobSidecars(ctx, &api.BlobSidecarsOpts{
				Block: root.String(),
			})
			if err != nil {
				return err
			}

			for _, blobSideCar := range blobSideCars.Data {
				if err := bs.storage.Save(root, blobSideCar); err != nil {
					bs.logger.Error().Uint64("slot", slot).Err(err).Msg("Failed to save blob")
					return err
				}
			}
			return nil
		}, retry.Attempts(5), retry.Delay(1*time.Second)); err != nil {
			bs.logger.Error().Uint64("slot", slot).Err(err).Msg("Failed to save blob")
		}

	}
	return nil
}
