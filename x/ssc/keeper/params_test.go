package keeper_test

import (
	"testing"

	testkeeper "github.com/sagaxyz/ssc/testutil/keeper"
	"github.com/sagaxyz/ssc/x/ssc/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.SscKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
