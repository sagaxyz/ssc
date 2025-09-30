package v02

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const Name = "0.1-to-0.2"

//nolint:staticcheck
func UpgradeHandler(mm *module.Manager, configurator module.Configurator, paramsKeeper paramskeeper.Keeper, consensusKeeper *consensuskeeper.Keeper, baseAppLegacySS paramstypes.Subspace) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		cp := baseapp.GetConsensusParams(sdkCtx, baseAppLegacySS)
		if cp == nil {
			return nil, fmt.Errorf("consensus parameters are undefined")
		}
		err := consensusKeeper.ParamsStore.Set(ctx, *cp)
		if err != nil {
			return nil, fmt.Errorf("failed to set consensus params: %w", err)
		}
		return mm.RunMigrations(ctx, configurator, vm)
	}
}
