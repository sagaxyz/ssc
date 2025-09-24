package v05

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

const (
	Name            = "0.5"
	tempMinterName  = "developer-credits"
	baseDenom       = "credit"
	recipientBech32 = "saga1a8duyed73q8gmewuakdfgyge52rkdgklysgfam"
)

var (
	// If "credit" has 6 decimals, this is 1,000,000.000000 CREDIT in base units.
	mintAmount = sdkmath.NewInt(1_000_000).MulRaw(1_000_000)
)

// ensureTempMinter creates the temp module account with a proper account number (idempotent).
func ensureTempMinter(ctx sdk.Context, ak authkeeper.AccountKeeper, name string, perms ...string) error {
	addr := authtypes.NewModuleAddress(name)
	if ak.GetAccount(ctx, addr) != nil {
		return nil // already exists
	}
	ma := authtypes.NewEmptyModuleAccount(name, perms...)
	acc := ak.NewAccount(ctx, ma) // assigns account number
	mai, ok := acc.(types.ModuleAccountI)
	if !ok {
		return fmt.Errorf("expected ModuleAccountI, got %T", acc)
	}
	ak.SetModuleAccount(ctx, mai)
	return nil
}

func UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		// Parse recipient AFTER prefixes are configured.
		recipientBz, err := ak.AddressCodec().StringToBytes(recipientBech32)
		if err != nil {
			return nil, fmt.Errorf("invalid recipient addr: %w", err)
		}
		recipient := sdk.AccAddress(recipientBz)

		// Run migrations first (determinism).
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}

		// Create ephemeral minter.
		if err := ensureTempMinter(sdkCtx, ak, tempMinterName, authtypes.Minter); err != nil {
			return nil, err
		}

		coins := sdk.NewCoins(sdk.NewCoin(baseDenom, mintAmount))

		// Mint and send.
		if err := bk.MintCoins(ctx, tempMinterName, coins); err != nil {
			return nil, err
		}
		if err := bk.SendCoinsFromModuleToAccount(ctx, tempMinterName, recipient, coins); err != nil {
			return nil, err
		}

		// Burn any dust left (should be zero).
		if bal := bk.GetAllBalances(ctx, authtypes.NewModuleAddress(tempMinterName)); !bal.IsZero() {
			if err := bk.BurnCoins(ctx, tempMinterName, bal); err != nil {
				return nil, err
			}
		}

		// Remove the temporary minter so no Minter perms remain.
		if acc := ak.GetAccount(sdkCtx, authtypes.NewModuleAddress(tempMinterName)); acc != nil {
			ak.RemoveAccount(sdkCtx, acc)
		}

		return newVM, nil
	}
}
