package restore

import "time"

const (
	serverTimeout = 60 * time.Second
	minBlobSlot   = 8626176 // mainnet minimum slot which save blobs
)

func NewConfig(beaconUrl string, timeout time.Duration, storagePath string) *Config {
	if timeout == 0 {
		timeout = serverTimeout
	}
	return &Config{
		beaconUrl:   beaconUrl,
		timeOut:     serverTimeout,
		storagePath: storagePath,
	}
}

type Config struct {
	beaconUrl string
	timeOut   time.Duration

	storagePath string
}
