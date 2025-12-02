package billing

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sagaxyz/ssc/x/billing/keeper"
	"github.com/sagaxyz/ssc/x/billing/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set params
	k.SetParams(ctx, genState.Params)

	// Import billing history
	for _, bh := range genState.BillingHistory {
		k.ImportBillingHistory(ctx, bh)
	}

	// Import validator payout history
	for _, vph := range genState.ValidatorPayoutHistory {
		k.ImportValidatorPayoutHistory(ctx, vph)
	}

	// this line is used by starport scaffolding # genesis/module/init
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Export billing history
	genesis.BillingHistory = k.ExportBillingHistory(ctx)

	// Export validator payout history
	genesis.ValidatorPayoutHistory = k.ExportValidatorPayoutHistory(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
