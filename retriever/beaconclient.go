package retriever

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
func NewBeaconClient(ctx context.Context, beaconUrl string, beaconType string, timeout time.Duration) (BeaconClient, error) {
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var withEnforceJSON bool
	if beaconType == "prysm" {
		withEnforceJSON = true
	}

	c, err := http.New(cctx,
		http.WithTimeout(timeout),
		http.WithAddress(beaconUrl),
		http.WithLogLevel(zerolog.ErrorLevel),
		http.WithEnforceJSON(withEnforceJSON),
	)
	if err != nil {
		return nil, err
	}

	return c.(*http.Service), nil
}
