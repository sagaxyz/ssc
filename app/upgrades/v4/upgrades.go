package v4

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	liquidmodulekeeper "github.com/cosmos/gaia/v25/x/liquid/keeper"
	liquidtypes "github.com/cosmos/gaia/v25/x/liquid/types"
)

const Name = "3-to-4"

func UpgradeHandler(mm *module.Manager, configurator module.Configurator, stakingKeeper stakingkeeper.Keeper, liquidKeeper liquidmodulekeeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}

		sdkCtx := sdk.UnwrapSDKContext(ctx)

		// set default params to the liquid module
		defaultParams := liquidtypes.DefaultParams()
		err = liquidKeeper.SetParams(sdkCtx, defaultParams)
		if err != nil {
			return vm, fmt.Errorf("error setting params: %w", err)
		}

		// add all validators to the liquid module
		liquidValidators := make(map[string]liquidtypes.LiquidValidator)
		allVals, err := stakingKeeper.GetAllValidators(ctx)
		if err != nil {
			return vm, fmt.Errorf("unable to get all validators: %w", err)
		}
		for _, val := range allVals {
			if _, ok := liquidValidators[val.OperatorAddress]; !ok {
				liquidValidators[val.OperatorAddress] = liquidtypes.NewLiquidValidator(val.OperatorAddress)
			}
		}

		for _, liquidVal := range liquidValidators {
			if err := liquidKeeper.SetLiquidValidator(ctx, liquidVal); err != nil {
				return vm, fmt.Errorf("error migrating liquid validator: %w", err)
			}
		}

		// because this is a new module, we don't need to migrate anything else

		return newVM, nil
	}
}
