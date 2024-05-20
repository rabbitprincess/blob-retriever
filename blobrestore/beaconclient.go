package blobrestore

import (
	"context"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/rs/zerolog"
)

type BeaconClient interface {
	client.BlobSidecarsProvider
	client.BeaconBlockHeadersProvider
}

// NewBeaconClient returns a new HTTP beacon client.
func NewBeaconClient(ctx context.Context, beaconUrl string, timeout uint64) (BeaconClient, error) {
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c, err := http.New(cctx,
		http.WithAddress(beaconUrl),
		http.WithTimeout(time.Duration(timeout)),
		http.WithLogLevel(zerolog.ErrorLevel),
		// http.WithEnforceJSON(cfg.EnforceJSON),
	)
	if err != nil {
		return nil, err
	}

	return c.(*http.Service), nil
}
