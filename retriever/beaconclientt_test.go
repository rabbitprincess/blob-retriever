package retriever

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/stretchr/testify/require"
)

const (
	beaconUrl = "https://ethereum-beacon-api.publicnode.com"
)

func TestGetBlockRootBySlot(t *testing.T) {
	if beaconUrl == "" {
		t.Skip("beaconUrl is not set")
	}

	ctx := context.Background()
	bc, err := NewBeaconClient(ctx, beaconUrl, time.Second*5)
	require.NoError(t, err)

	for _, test := range []struct {
		slot uint64
		root string
	}{
		{8639468, "0x9e6234130da3a3c4c1388811ae8cebcca966fa6951573f0ab418b4cdb433fadc"},
		{8636176, "0xc5a02ee7b62f49e62489c7e2efe1b6470308836f44d33a0181c1b98c5e5ab368"},
		{8759552, "0x4cdf12b7bd58317d9e09deb3e37ef8ba33d179a41d5f9b07d753c17a5a61a154"},
		{8751935, "0x781f81676622ca98edd3d3013420c44b5679c56ddd452b1e1da5334db4cbd3fe"},
	} {
		header, err := bc.BeaconBlockHeader(ctx, &api.BeaconBlockHeaderOpts{
			Block: strconv.FormatUint(test.slot, 10),
		})
		if err != nil {
			require.Equal(t, test.root, "", test.slot)
		} else {
			require.Equal(t, test.root, header.Data.Root.String(), test.slot)
		}
		if test.root != "" {
			header, err = bc.BeaconBlockHeader(ctx, &api.BeaconBlockHeaderOpts{
				Block: test.root,
			})
			require.NoError(t, err)
			require.Equal(t, test.slot, uint64(header.Data.Header.Message.Slot), test.root)
		}
	}
}
