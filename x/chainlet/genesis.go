package chainlet

import (
	"github.com/sagaxyz/ssc/x/chainlet/keeper"
	"github.com/sagaxyz/ssc/x/chainlet/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k *keeper.Keeper, genState types.GenesisState) {
	// Set params
	k.SetParams(ctx, genState.Params)

	// Set the port ID for the chainlet module
	// In IBC v10, port binding is handled automatically when the module is registered in the router
	k.SetPort(ctx, genState.PortId)

	// Set chainlet count
	k.SetChainletCount(ctx, genState.ChainletCount)

	// Import chainlet stacks first (chainlets depend on stacks)
	for _, stack := range genState.ChainletStacks {
		if err := k.ImportChainletStack(ctx, stack); err != nil {
			panic(err)
		}
	}

	// Import chainlets
	for _, chainlet := range genState.Chainlets {
		if err := k.ImportChainlet(ctx, chainlet); err != nil {
			panic(err)
		}
	}

	// this line is used by starport scaffolding # genesis/module/init
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.PortId = k.GetPort(ctx)
	genesis.ChainletCount = k.GetChainletCount(ctx)

	// Export all chainlets
	genesis.Chainlets = k.ExportChainlets(ctx)

	// Export all chainlet stacks
	genesis.ChainletStacks = k.ExportChainletStacks(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
