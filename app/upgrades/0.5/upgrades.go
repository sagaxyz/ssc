package v05

import (
	"context"

	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

const (
	Name           = "0.5"
	tempMinterName = "developer-credits"
	baseDenom      = "ucday"
)

var addr = "saga1a8duyed73q8gmewuakdfgyge52rkdgklysgfam"

var mintAmount = sdkmath.NewInt(1_000_000).MulRaw(1_000_000)

func UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		// 1) Run module migrations first for determinism
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}

		sdkCtx := sdk.UnwrapSDKContext(ctx)

		// decode here (after bech32 prefixes are configured)
		recipientAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return nil, err
		}
		// 2) Create a temporary minter module account in state (NOT in maccPerms)
		//    (If it already exists for some reason, we reuse it.)
		if ak.GetModuleAddress(tempMinterName) == nil {
			ma := authtypes.NewEmptyModuleAccount(tempMinterName, authtypes.Minter)
			ak.SetModuleAccount(sdkCtx, ma)
		}

		// 3) Mint the one-off amount
		coins := sdk.NewCoins(sdk.NewCoin(baseDenom, mintAmount))
		if err := bk.MintCoins(ctx, tempMinterName, coins); err != nil {
			return nil, err
		}

		// 4) Send all freshly minted coins to the single recipient
		if err := bk.SendCoinsFromModuleToAccount(ctx, tempMinterName, recipientAddr, coins); err != nil {
			return nil, err
		}

		// 5) Safety: if anything remains on the temp account, burn it (should be zero)
		if bal := bk.GetAllBalances(ctx, ak.GetModuleAddress(tempMinterName)); !bal.IsZero() {
			if err := bk.BurnCoins(ctx, tempMinterName, bal); err != nil {
				return nil, err
			}
		}

		// 6) Remove the temporary module account from state so no minter remains
		if acc := ak.GetAccount(sdkCtx, ak.GetModuleAddress(tempMinterName)); acc != nil {
			ak.RemoveAccount(sdkCtx, acc)
		}

		return newVM, nil
	}
}
