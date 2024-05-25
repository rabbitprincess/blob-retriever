package main

import (
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
)

func flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "mode",
			Aliases:     []string{"m"},
			Value:       getEnv("MODE", "retrieve"),
			Usage:       "run mode (retrieve / check)",
			Destination: &mode,
		},
		&cli.StringFlag{
			Name:        "beacon_url",
			Aliases:     []string{"b"},
			Value:       getEnv("BEACON_URL", ""),
			Usage:       "Beacon node URL",
			Destination: &beaconUrl,
		},
		&cli.StringFlag{
			Name:        "beacon_type",
			Aliases:     []string{"n"},
			Value:       getEnv("BEACON_TYPE", "any"),
			Usage:       "Beacon node network type (any or prysm)",
			Destination: &beaconType,
		},
		&cli.StringFlag{
			Name:        "data",
			Aliases:     []string{"d"},
			Value:       getEnv("DATA_PATH", ""),
			Usage:       "data path to store blobs",
			Destination: &dataPath,
		},
		&cli.Uint64Flag{
			Name:        "worker",
			Aliases:     []string{"w"},
			Value:       getEnvAsUint64("NUM_WORKER", 1),
			Usage:       "number of workers",
			Destination: &numWorker,
		},
		&cli.Uint64Flag{
			Name:        "from",
			Aliases:     []string{"f"},
			Value:       getEnvAsUint64("FROM_SLOT", 0),
			Usage:       "from slot",
			Destination: &fromSlot,
		},
		&cli.Uint64Flag{
			Name:        "to",
			Aliases:     []string{"t"},
			Value:       getEnvAsUint64("TO_SLOT", 0),
			Usage:       "to slot",
			Destination: &toSlot,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsUint64(name string, defaultValue uint64) uint64 {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseUint(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}
