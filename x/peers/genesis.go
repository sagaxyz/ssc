package peers

import (
	"github.com/sagaxyz/ssc/x/peers/keeper"
	"github.com/sagaxyz/ssc/x/peers/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set params
	k.SetParams(ctx, genState.Params)

	// Import peer data
	for _, pd := range genState.PeerData {
		k.ImportPeerData(ctx, pd.ChainId, pd.ValidatorAddress, pd.Data)
	}

	// Import chain counters
	for _, cc := range genState.ChainCounters {
		k.ImportChainCounter(ctx, cc.ChainId, cc.Counter)
	}

	// this line is used by starport scaffolding # genesis/module/init
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Export peer data
	genesis.PeerData = k.ExportPeerData(ctx)

	// Export chain counters
	genesis.ChainCounters = k.ExportChainCounters(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
