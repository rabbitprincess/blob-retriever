package retriever

import "time"

const (
	serverTimeout = 60 * time.Second
	minBlobSlot   = 8626176 // mainnet minimum slot which save blobs
)

func NewConfig(beaconUrl string, timeout time.Duration, storageType, storagePath string) *Config {
	if timeout == 0 {
		timeout = serverTimeout
	}
	if storageType == "" {
		storageType = "prysm"
	}
	return &Config{
		BeaconUrl:   beaconUrl,
		Timeout:     serverTimeout,
		StorageType: storageType,
		StoragePath: storagePath,
	}
}

type Config struct {
	Mode        string
	BeaconUrl   string
	Timeout     time.Duration
	StorageType string
	StoragePath string
}
