package main

import (
	"context"
	"os"

	"github.com/rabbitprincess/blob-retriever/blob"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "blob_restorer",
		Run: rootRun,
	}
	mode        string
	beaconUrl   string
	storagePath string
	storageType string
	fromSlot    uint64
	toSlot      uint64
)

func init() {
	fs := rootCmd.PersistentFlags()
	fs.StringVarP(&mode, "mode", "m", "check", "run mode (retrieve / check)")
	fs.StringVarP(&beaconUrl, "beacon", "b", "https://ethereum-beacon-api.publicnode.com", "Beacon node URL")
	fs.StringVarP(&storagePath, "storage_path", "p", "./data", "Path to store blobs")
	// only support prysm for now
	// fs.StringVarP(&storagePath, "storage_type", "s", "prysm", "Type to storage ( prysm or lighthouse )")
	fs.Uint64VarP(&fromSlot, "from", "f", 9084721, "Start slot")
	fs.Uint64VarP(&toSlot, "to", "t", 9084733, "End slot")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func rootRun(cmd *cobra.Command, args []string) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx := context.Background()

	cfg := blob.NewConfig(beaconUrl, 0, storageType, storagePath)
	blobRestore := blob.NewBlobRestore(ctx, logger, cfg)

	logger.Info().Str("mode", mode).Uint64("from slot", fromSlot).Uint64("to slot", toSlot).Msg("Run blob retriever")
	blobRestore.Run(ctx, mode, fromSlot, toSlot)
}
