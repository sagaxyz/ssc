package v4

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
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

		// set params to the liquid module
		params := liquidtypes.Params{
			GlobalLiquidStakingCap:    math.LegacyNewDecFromIntWithPrec(math.NewInt(80), 2), // 0.80
			ValidatorLiquidStakingCap: math.LegacyOneDec(),                                  // 1
		}
		err = liquidKeeper.SetParams(ctx, params)
		if err != nil {
			return vm, fmt.Errorf("error setting params: %w", err)
		}

		// add all validators to the liquid module
		allVals, err := stakingKeeper.GetAllValidators(ctx)
		if err != nil {
			return vm, fmt.Errorf("unable to get all validators: %w", err)
		}
		for _, val := range allVals {
			lVal := liquidtypes.NewLiquidValidator(val.OperatorAddress)
			if err := liquidKeeper.SetLiquidValidator(ctx, lVal); err != nil {
				return vm, fmt.Errorf("error migrating liquid validator: %w", err)
			}
		}

		// because this is a new module, we don't need to migrate anything else

		return newVM, nil
	}
}
