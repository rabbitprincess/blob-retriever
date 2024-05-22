package retriever

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/avast/retry-go"
	"github.com/gammazero/workerpool"
	"github.com/rabbitprincess/blob-retriever/storage"
	"github.com/rs/zerolog"
)

type BlobRestore struct {
	cfg     *Config
	logger  zerolog.Logger
	wp      *workerpool.WorkerPool
	client  BeaconClient
	storage storage.BlobStore
}

// NewBlobRestore
func NewBlobRestore(ctx context.Context, log zerolog.Logger, cfg *Config) *BlobRestore {
	wp := workerpool.New(int(cfg.NumWorker))
	client, err := NewBeaconClient(ctx, cfg.BeaconUrl, cfg.Timeout)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create beacon client")
	}
	storage, err := storage.NewPrysmBlobStorage(log, cfg.StoragePath)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create blob storage")
	}
	return &BlobRestore{
		cfg:     cfg,
		logger:  log,
		wp:      wp,
		client:  client,
		storage: storage,
	}
}

func (bs *BlobRestore) Run(ctx context.Context, mode string, fromSlot, toSlot uint64) error {
	if fromSlot < minBlobSlot {
		bs.logger.Warn().Uint64("fromSlot", fromSlot).Uint64("minBlobSlot", minBlobSlot).Msg("fromSlot is less than minBlobSlot, set fromSlot to minBlobSlot")
		fromSlot = minBlobSlot
	}
	if toSlot < fromSlot {
		bs.logger.Warn().Uint64("toSlot", toSlot).Uint64("fromSlot", fromSlot).Msg("toSlot is less than fromSlot, set toSlot to fromSlot")
		toSlot = fromSlot
	}

	for slot := fromSlot; slot <= toSlot; slot++ {
		bs.wp.Submit(func() {
			header, sidecars, err := bs.GetV1BlobFromApi(ctx, slot)
			if err != nil || header.Root.IsZero() {
				bs.logger.Info().Uint64("slot", slot).Err(err).Msg("block is not existing, continue...")
				return
			}
			if len(sidecars) == 0 {
				bs.logger.Info().Uint64("slot", slot).Msg("blob sidecars not exist, continue...")
				return
			}

			switch mode {
			case "retrieve":
				if err := bs.RestoreBlob(ctx, slot, header, sidecars); err != nil {
					bs.logger.Panic().Uint64("slot", slot).Err(err).Msg("Failed to restore blob")
				}
			case "check":
				if err := bs.CheckBlobSidecar(ctx, slot, header, sidecars); err != nil {
					bs.logger.Error().Uint64("slot", slot).Err(err).Msg("Failed to check blob sidecar")
				}
			}
		})
	}
	return nil
}

func (bs *BlobRestore) RestoreBlob(ctx context.Context, slot uint64, header *apiv1.BeaconBlockHeader, sidecars []*deneb.BlobSidecar) error {
	if bs.storage.Exist(header.Root) {
		bs.logger.Info().Str("root", header.Root.String()).Msg("Blob already exists in storage, continue")
		return nil
	}

	for _, sidecar := range sidecars {
		if err := bs.storage.Save(header.Root, sidecar); err != nil {
			bs.logger.Error().Uint64("slot", slot).Str("root", header.Root.String()).Err(err).Msg("Failed to save blob sidecar")
			return err
		}
		bs.logger.Info().Uint64("slot", slot).Str("root", header.Root.String()).Uint64("index", uint64(sidecar.Index)).Msg("Blob sidecar saved")
	}
	return nil
}

func (bs *BlobRestore) CheckBlobSidecar(ctx context.Context, slot uint64, header *apiv1.BeaconBlockHeader, sidecars []*deneb.BlobSidecar) error {
	for _, sidecar := range sidecars {
		valid, err := bs.storage.Valid(header.Root, sidecar)
		if err != nil {
			return err
		}
		if !valid {
			err = fmt.Errorf("Blob sidecar is not valid")
			bs.logger.Error().Err(err).Uint64("slot", slot).Str("root", header.Root.String()).Uint64("index", uint64(sidecar.Index)).Msg("Blob sidecar is not valid")
			return err
		}
		bs.logger.Info().Uint64("slot", slot).Str("root", header.Root.String()).Uint64("index", uint64(sidecar.Index)).Msg("Blob sidecar is valid")
	}
	return nil
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
	}, retry.Attempts(1), retry.Delay(100*time.Millisecond)); err != nil {
		bs.logger.Error().Uint64("slot", slot).Err(err).Msg("Failed to get block header")
		return nil, nil, err
	}

	return header, sidecars, nil
}
