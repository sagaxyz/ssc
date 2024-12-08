package v02

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	ibcclientkeeper "github.com/cosmos/ibc-go/v8/modules/core/02-client/keeper"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

const Name = "0.1-to-0.2"

func UpgradeHandler(mm *module.Manager, configurator module.Configurator, paramsKeeper paramskeeper.Keeper, consensuskeeper *consensuskeeper.Keeper, clientKeeper ibcclientkeeper.Keeper) upgradetypes.UpgradeHandler {

	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		params := clientKeeper.GetParams(sdkCtx)
		params.AllowedClients = append(params.AllowedClients, exported.Localhost)
		clientKeeper.SetParams(sdkCtx, params)
		return newVM, nil
	}
}
