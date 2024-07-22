package chainlet_test

import (
	"testing"

	keepertest "github.com/sagaxyz/ssc/testutil/keeper"
	"github.com/sagaxyz/ssc/testutil/nullify"

	"github.com/sagaxyz/ssc/x/chainlet"
	"github.com/sagaxyz/ssc/x/chainlet/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ChainletKeeper(t)
	chainlet.InitGenesis(ctx, k, genesisState)
	got := chainlet.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.NumChainlets, got.NumChainlets)
	// this line is used by starport scaffolding # genesis/test/assert
}
