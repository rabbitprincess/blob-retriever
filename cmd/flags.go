package main

import (
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
)

var (
	mode      string
	apiUrl    string
	apiType   string
	dataPath  string
	dataType  string
	numWorker uint64
	fromSlot  uint64
	toSlot    uint64
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
			Name:        "api_url",
			Aliases:     []string{"u"},
			Value:       getEnv("API_URL", ""),
			Usage:       "Beacon node URL",
			Destination: &apiUrl,
		},
		&cli.StringFlag{
			Name:        "api_type",
			Aliases:     []string{"a"},
			Value:       getEnv("API_TYPE", "any"),
			Usage:       "Beacon node network type (any or prysm)",
			Destination: &apiType,
		},
		&cli.StringFlag{
			Name:        "data_path",
			Aliases:     []string{"d"},
			Value:       getEnv("DATA_PATH", "./blobs"),
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
			Usage:       "from slot. minimum is 8626176",
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
