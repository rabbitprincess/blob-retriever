package main

import (
	"context"
	"os"

	"github.com/rabbitprincess/blob-retriever/retriever"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "blob_retriever",
		Run: rootRun,
	}
	mode        string
	beaconUrl   string
	storagePath string
	storageType string
	numWorker   uint8
	fromSlot    uint64
	toSlot      uint64
)

func init() {
	fs := rootCmd.PersistentFlags()
	fs.StringVarP(&mode, "mode", "m", "retrieve", "run mode (retrieve / check)")
	fs.StringVarP(&beaconUrl, "beacon", "b", "", "Beacon node URL")
	fs.StringVarP(&storagePath, "storage_path", "p", "", "Path to store blobs")
	// only support prysm for now
	// fs.StringVarP(&storagePath, "storage_type", "s", "prysm", "Type to storage ( prysm or lighthouse )")
	fs.Uint8VarP(&numWorker, "worker", "w", 1, "Number of worker")
	fs.Uint64VarP(&fromSlot, "from", "f", 0, "Start slot")
	fs.Uint64VarP(&toSlot, "to", "t", 0, "End slot")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func rootRun(cmd *cobra.Command, args []string) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx := context.Background()

	cfg := retriever.NewConfig(beaconUrl, 0, storageType, storagePath, numWorker)
	blobRestore := retriever.NewBlobRestore(ctx, logger, cfg)

	logger.Info().Str("mode", mode).Uint64("from slot", fromSlot).Uint64("to slot", toSlot).Msg("Run blob retriever")
	blobRestore.Run(ctx, mode, fromSlot, toSlot)
}
