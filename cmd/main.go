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

	beaconUrl   string
	storagePath string
	storageType string
	fromSlot    uint64
	toSlot      uint64
)

func init() {
	fs := rootCmd.PersistentFlags()
	fs.StringVarP(&beaconUrl, "beacon", "b", "", "Beacon node URL")
	fs.StringVarP(&storagePath, "storage_path", "p", "", "Path to store blobs")
	// only support prysm for now
	// fs.StringVarP(&storagePath, "storage_type", "s", "prysm", "Type to storage ( prysm or lighthouse )")
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

	blobRestore := blob.NewBlobRestore(ctx, logger, &blob.Config{
		BeaconUrl:   beaconUrl,
		StoragePath: storagePath,
		StorageType: storageType,
	})

	blobRestore.RestoreBlob(ctx, fromSlot, toSlot)
}
