package escrow_test

import (
	"testing"

	keepertest "github.com/sagaxyz/ssc/testutil/keeper"
	"github.com/sagaxyz/ssc/testutil/nullify"
	"github.com/sagaxyz/ssc/x/escrow"
	"github.com/sagaxyz/ssc/x/escrow/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.EscrowKeeper(t)
	escrow.InitGenesis(ctx, *k, genesisState)
	got := escrow.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
