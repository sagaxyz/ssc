package billing_test

import (
	"testing"

	keepertest "github.com/sagaxyz/ssc/testutil/keeper"
	"github.com/sagaxyz/ssc/testutil/nullify"
	"github.com/sagaxyz/ssc/x/billing"
	"github.com/sagaxyz/ssc/x/billing/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BillingKeeper(t)
	billing.InitGenesis(ctx, *k, genesisState)
	got := billing.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
