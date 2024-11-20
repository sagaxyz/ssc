package keeper_test

import (
	"testing"

	testkeeper "github.com/sagaxyz/ssc/testutil/keeper"

	"github.com/sagaxyz/ssc/x/chainlet/types"

	"github.com/stretchr/testify/require"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.ChainletKeeper(t)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
