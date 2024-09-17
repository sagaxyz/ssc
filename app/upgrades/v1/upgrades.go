package v1

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
)

const Name = "0.x-to-1" //TODO replace x with the actual number

func UpgradeHandler(mm *module.Manager, configurator module.Configurator, paramsKeeper paramskeeper.Keeper, consensusKeeper *consensuskeeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, configurator, vm)
	}
}
