package v1

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	ccvprovider "github.com/cosmos/interchain-security/v7/x/ccv/provider"
	ccvproviderkeeper "github.com/cosmos/interchain-security/v7/x/ccv/provider/keeper"
	ccvprovidertypes "github.com/cosmos/interchain-security/v7/x/ccv/provider/types"
)

const (
	Name            = "0.5-to-1.0"
	tempMinterName  = "developer-credits"
	baseDenom       = "credit"
	recipientBech32 = "saga1a8duyed73q8gmewuakdfgyge52rkdgklysgfam"
)

var (
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
	mai, ok := acc.(sdk.ModuleAccountI)
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
	providerKeeper ccvproviderkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		// Prevent RunMigrations from thinking provider is new
		if _, exists := vm[ccvprovidertypes.ModuleName]; !exists {
			vm[ccvprovidertypes.ModuleName] = ccvprovider.AppModule{}.ConsensusVersion()
		}

		// Run migrations (provider will now be skipped)
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}

		// Manually initialize provider store (no validator updates)
		genState := ccvprovidertypes.DefaultGenesisState()
		providerKeeper.InitGenesis(sdkCtx, genState)

		newVM[ccvprovidertypes.ModuleName] = ccvprovider.AppModule{}.ConsensusVersion()

		if err := ensureTempMinter(sdkCtx, ak, tempMinterName, authtypes.Minter); err != nil {
			return nil, err
		}

		coins := sdk.NewCoins(sdk.NewCoin(baseDenom, mintAmount))

		if err := bk.MintCoins(ctx, tempMinterName, coins); err != nil {
			return nil, err
		}
		recipientBz, err := ak.AddressCodec().StringToBytes(recipientBech32)
		if err != nil {
			return nil, fmt.Errorf("invalid recipient addr: %w", err)
		}
		recipient := sdk.AccAddress(recipientBz)

		if err := bk.SendCoinsFromModuleToAccount(ctx, tempMinterName, recipient, coins); err != nil {
			return nil, err
		}

		if bal := bk.GetAllBalances(ctx, authtypes.NewModuleAddress(tempMinterName)); !bal.IsZero() {
			if err := bk.BurnCoins(ctx, tempMinterName, bal); err != nil {
				return nil, err
			}
		}

		if acc := ak.GetAccount(sdkCtx, authtypes.NewModuleAddress(tempMinterName)); acc != nil {
			ak.RemoveAccount(sdkCtx, acc)
		}

		return newVM, nil
	}
}
