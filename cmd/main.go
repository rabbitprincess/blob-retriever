package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rabbitprincess/blob-retriever/retriever"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var (
	mode        string
	beaconUrl   string
	beaconType  string
	dataPath    string
	storageType string
	numWorker   uint64
	fromSlot    uint64
	toSlot      uint64
)

func main() {
	app := &cli.App{
		Name:  "blob_retriever",
		Usage: "Retrieve and check blobs",
		Flags: flags(),
		Action: func(c *cli.Context) error {
			return rootRun()
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func rootRun() error {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := retriever.NewConfig(beaconUrl, beaconType, 0, storageType, dataPath, numWorker)
	blobRestore := retriever.NewBlobRetriever(ctx, logger, cfg)

	logger.Info().Str("mode", mode).Uint64("from slot", fromSlot).Uint64("to slot", toSlot).Msg("Run blob retriever")

	go func() {
		blobRestore.Run(ctx, mode, fromSlot, toSlot)
		cancel()
	}()

	// Wait main routine to stop
	handleInterrupt(logger, cancel)
	return nil
}

func handleInterrupt(logger zerolog.Logger, cancelFunc context.CancelFunc) {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	sig := <-sigChannel
	logger.Info().Msgf("Received signal %s, shutting down...", sig)
	cancelFunc()
}
