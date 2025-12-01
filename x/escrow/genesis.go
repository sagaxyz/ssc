package escrow

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sagaxyz/ssc/x/escrow/keeper"
	"github.com/sagaxyz/ssc/x/escrow/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set params
	k.SetParams(ctx, genState.Params)

	// Import chainlet accounts
	for _, acc := range genState.ChainletAccounts {
		k.ImportChainletAccount(ctx, acc)
	}

	// Import pools
	for _, pool := range genState.Pools {
		k.ImportPool(ctx, pool)
	}

	// Import funders
	for _, gf := range genState.Funders {
		k.ImportFunder(ctx, gf.ChainId, gf.Denom, gf.Address, gf.Funder)
	}

	// this line is used by starport scaffolding # genesis/module/init
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Export chainlet accounts
	genesis.ChainletAccounts = k.ExportChainletAccounts(ctx)

	// Export pools
	genesis.Pools = k.ExportPools(ctx)

	// Export funders
	genesis.Funders = k.ExportFunders(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
