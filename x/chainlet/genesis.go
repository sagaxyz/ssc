package chainlet

import (
	"github.com/sagaxyz/ssc/x/chainlet/keeper"
	"github.com/sagaxyz/ssc/x/chainlet/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k *keeper.Keeper, genState types.GenesisState) {
	k.InitializeChainletCount(ctx)
	// this line is used by starport scaffolding # genesis/module/init

	k.SetParams(ctx, genState.Params)

	// Set the port ID for the chainlet module
	// In IBC v10, port binding is handled automatically when the module is registered in the router
	k.SetPort(ctx, genState.PortId)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.PortId = k.GetPort(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
