package retriever

import "time"

const (
	serverTimeout = 60 * time.Second
	minBlobSlot   = 8626176 // mainnet minimum slot which save blobs
)

func NewConfig(beaconUrl, beaconType string, timeout time.Duration, storageType, storagePath string, numWorker uint64) *Config {
	if timeout == 0 {
		timeout = serverTimeout
	}
	if storageType == "" {
		storageType = "prysm"
	}
	return &Config{
		BeaconApiUrl:  beaconUrl,
		BeaconApiType: beaconType,
		Timeout:       serverTimeout,
		StorageType:   storageType,
		StoragePath:   storagePath,
		NumWorker:     numWorker,
	}
}

type Config struct {
	Mode          string
	BeaconApiUrl  string
	BeaconApiType string
	Timeout       time.Duration
	StorageType   string
	StoragePath   string
	NumWorker     uint64
}
