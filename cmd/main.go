package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rabbitprincess/blob-retriever/retriever"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
	blobRetriever := retriever.NewBlobRetriever(ctx, logger, cfg)
	if blobRetriever == nil {
		logger.Error().Msg("Failed to create blob retriever")
		return nil
	}

	logger.Info().Str("mode", mode).Uint64("from slot", fromSlot).Uint64("to slot", toSlot).Msg("Run blob retriever")

	go func() {
		blobRetriever.Run(ctx, mode, fromSlot, toSlot)
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
