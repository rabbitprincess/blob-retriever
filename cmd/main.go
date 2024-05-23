package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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
	beaconType  string
	dataPath    string
	storageType string
	numWorker   uint8
	fromSlot    uint64
	toSlot      uint64
)

func init() {
	fs := rootCmd.PersistentFlags()
	fs.StringVarP(&mode, "mode", "m", "retrieve", "run mode (retrieve / check)")
	fs.StringVarP(&beaconUrl, "beacon_url", "b", "", "Beacon node URL")
	fs.StringVarP(&beaconType, "beacon_type", "n", "any", "Beacon node network type. ( any or prysm )")
	fs.StringVarP(&dataPath, "data", "d", "", "data path to store blobs")
	// only support prysm for now
	// fs.StringVarP(&storageType, "storage_type", "s", "prysm", "Type to storage ( prysm or lighthouse )")
	fs.Uint8VarP(&numWorker, "worker", "w", 1, "number of worker")
	fs.Uint64VarP(&fromSlot, "from", "f", 0, "start slot")
	fs.Uint64VarP(&toSlot, "to", "t", 0, "end slot")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func rootRun(cmd *cobra.Command, args []string) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx := context.Background()

	cfg := retriever.NewConfig(beaconUrl, beaconType, 0, storageType, dataPath, numWorker)
	blobRestore := retriever.NewBlobRestore(ctx, logger, cfg)

	logger.Info().Str("mode", mode).Uint64("from slot", fromSlot).Uint64("to slot", toSlot).Msg("Run blob retriever")

	interrupt := handleKillSig(func() {
	}, logger)

	go func() {
		blobRestore.Run(ctx, mode, fromSlot, toSlot)
		interrupt.C <- struct{}{}
	}()

	// Wait main routine to stop
	<-interrupt.C
}

type interrupt struct {
	C chan struct{}
}

func handleKillSig(handler func(), logger zerolog.Logger) interrupt {
	i := interrupt{
		C: make(chan struct{}),
	}

	sigChannel := make(chan os.Signal, 1)

	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for signal := range sigChannel {
			logger.Info().Msgf("Receive signal %s, Shutting down...", signal)
			handler()
			close(i.C)
		}
	}()
	return i
}
