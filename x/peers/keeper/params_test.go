package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	testkeeper "github.com/sagaxyz/ssc/testutil/keeper"
	"github.com/sagaxyz/ssc/x/peers/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.PeersKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
