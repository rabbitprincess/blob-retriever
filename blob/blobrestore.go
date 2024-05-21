package blob

import (
	"context"
	"strconv"
	"time"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/avast/retry-go"
	"github.com/rs/zerolog"
)

type BlobRestore struct {
	cfg     *Config
	logger  zerolog.Logger
	client  BeaconClient
	storage BlobStore
}

// NewBlobRestore
func NewBlobRestore(ctx context.Context, logger zerolog.Logger, cfg *Config) *BlobRestore {
	client, err := NewBeaconClient(ctx, cfg.BeaconUrl, uint64(cfg.Timeout.Seconds()))
	if err != nil {
		logger.Panic().Err(err).Msg("Failed to create beacon client")
	}
	storage, err := NewPrysmBlobStorage(cfg.StoragePath)
	if err != nil {
		logger.Panic().Err(err).Msg("Failed to create blob storage")
	}
	return &BlobRestore{
		cfg:     cfg,
		logger:  logger,
		client:  client,
		storage: storage,
	}
}

func (bs *BlobRestore) RestoreBlob(ctx context.Context, fromSlot, toSlot uint64) error {
	if fromSlot < minBlobSlot {
		bs.logger.Warn().Uint64("fromSlot", fromSlot).Uint64("minBlobSlot", minBlobSlot).Msg("fromSlot is less than minBlobSlot, set fromSlot to minBlobSlot")
		fromSlot = minBlobSlot
	}
	if toSlot < fromSlot {
		bs.logger.Warn().Uint64("toSlot", toSlot).Uint64("fromSlot", fromSlot).Msg("toSlot is less than fromSlot, set toSlot to fromSlot")
		toSlot = fromSlot
	}

	for slot := fromSlot; slot <= toSlot; slot++ {
		var root phase0.Root
		data, sidecars, err := bs.GetV1BlobFromApi(ctx, slot)
		if err != nil {
			bs.logger.Panic().Uint64("slot", slot).Err(err).Msg("Failed to get block header")
		}
		if data.Root.IsZero() {
			bs.logger.Error().Uint64("slot", slot).Msg("Block not exist in api")
			continue
		}
		if bs.storage.Exist(root) {
			bs.logger.Info().Uint64("slot", slot).Str("root", root.String()).Msg("Blob already exists in storage, continue")
			continue
		}

		// Save blob sidecar
		for _, sidecar := range sidecars {
			if err := bs.storage.Save(root, sidecar); err != nil {
				bs.logger.Error().Uint64("slot", slot).Err(err).Msg("Failed to save blob sidecar")
				continue
			}
		}
	}
	return nil
}

func (bs *BlobRestore) CheckBlobSidecar(ctx context.Context, slot, index uint64) (bool, error) {
	header, sidecar, err := bs.GetV1BlobFromApi(ctx, slot)
	if err != nil {
		return false, err
	}

	if len(sidecar) <= int(index) || sidecar[index] == nil {
		return false, nil
	}
	return bs.storage.Valid(header.Root, sidecar[index])
}

func (bs *BlobRestore) GetV1BlobFromApi(ctx context.Context, slot uint64) (*apiv1.BeaconBlockHeader, []*deneb.BlobSidecar, error) {
	var header *apiv1.BeaconBlockHeader
	var sidecars []*deneb.BlobSidecar
	if err := retry.Do(func() error {
		res, err := bs.client.BeaconBlockHeader(ctx, &api.BeaconBlockHeaderOpts{
			Block: strconv.FormatUint(slot, 10),
		})
		if err != nil {
			bs.logger.Error().Uint64("slot", slot).Err(err).Msg("Failed to get block root | retrying...")
			return err
		}
		header = res.Data

		if !res.Data.Root.IsZero() {
			blobSideCars, err := bs.client.BlobSidecars(ctx, &api.BlobSidecarsOpts{
				Block: header.Root.String(),
			})
			if err != nil {
				return err
			}
			sidecars = blobSideCars.Data
		}

		return nil
	}, retry.Attempts(5), retry.Delay(1*time.Second)); err != nil {
		bs.logger.Error().Uint64("slot", slot).Err(err).Msg("Failed to get block header")
		return nil, nil, err
	}

	return header, sidecars, nil
}
