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
	billingkeeper "github.com/sagaxyz/ssc/x/billing/keeper"

	aclkeeper "github.com/sagaxyz/saga-sdk/x/acl/keeper"
	chainletkeeper "github.com/sagaxyz/ssc/x/chainlet/keeper"
)

const (
	Name            = "1.0"
	tempMinterName  = "developer-credits"
	baseDenom       = "credit"
	recipientBech32 = "saga1a8duyed73q8gmewuakdfgyge52rkdgklysgfam"
)

var (
	// 1,000,000 CREDIT with 6 decimals
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
	aclKeeper aclkeeper.Keeper,
	chainletKeeper chainletkeeper.Keeper,
	billingKeeper billingkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		// ------------------------------------------------------------------
		// 1. Run module migrations, treating provider as existing
		// ------------------------------------------------------------------

		if _, exists := vm[ccvprovidertypes.ModuleName]; !exists {
			vm[ccvprovidertypes.ModuleName] = ccvprovider.AppModule{}.ConsensusVersion()
		}

		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}

		// ------------------------------------------------------------------
		// 2. Initialize provider store WITHOUT validator updates
		// ------------------------------------------------------------------

		genState := ccvprovidertypes.DefaultGenesisState()
		providerKeeper.InitGenesis(sdkCtx, genState)
		newVM[ccvprovidertypes.ModuleName] = ccvprovider.AppModule{}.ConsensusVersion()

		// ------------------------------------------------------------------
		// 3. Fix chainlet params (match SPC behavior)
		// ------------------------------------------------------------------

		chainletParams := chainletKeeper.GetParams(sdkCtx) // or ctx if your keeper uses context.Context
		chainletParams.ChainletStackProtections = true
		chainletParams.EnableCCV = false
		chainletKeeper.SetParams(sdkCtx, chainletParams)

		// ------------------------------------------------------------------
		// 4. Patch ACL genesis:
		//    - enable = true
		//    - Admins = SPC allowed list
		//    - Allowed = SPC allowed list
		// ------------------------------------------------------------------

		aclParams := aclKeeper.GetParams(sdkCtx)
		aclParams.Enable = true
		aclKeeper.SetParams(sdkCtx, aclParams)

		aclGen := aclKeeper.ExportGenesis(sdkCtx) // returns acltypes.GenesisState

		addresses := []string{
			"saga1rdssl22ysxyendrkh2exw9zm7hvj8d2ju346g3",
			"saga1rcs5sw5yy9r04xsultcqv6tj73408qnawmlxqw",
			"saga1yuvju0cztlahsf6f37z9j83vwyzgj6pzhx090f",
			"saga1gme3rzzddpf4hkdngpruz5e4739lqsyyakgu0j",
			"saga1sz83y27774xwrahwmv5afutv86grc286hcf7w5",
			"saga16p4cejpaqpuha65hqyj85k5lx4umw7qzku37eg",
			"saga1u2a8ktctqhpx655ysw7ru27t6hqt9wlq4fn5ca",
			"saga1uccxg0ud23424ssuddqnkgjlsz2f6rvqlgjf9t",
			"saga17x049ugfafggn823dsnf32fhj5qlhlxrrzdz22",
			"saga17gk4chqd0lrkyamrxdmu62czmu0dpnemmxlymn",
		}

		aclGen.Allowed = append([]string(nil), addresses...)
		aclGen.Admins = append([]string(nil), addresses...)

		aclKeeper.InitGenesis(sdkCtx, aclGen)

		// After this:
		//   spcd/sscd q acl list-allowed -> those 10 addresses
		//   spcd/sscd q acl params       -> enable: true

		// ------------------------------------------------------------------
		// 5. One-shot dev credit mint + cleanup
		// ------------------------------------------------------------------

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

		// Set platform validators
		billingparams := billingKeeper.GetParams(sdkCtx)
		billingparams.PlatformValidators = []string{}
		billingKeeper.SetParams(sdkCtx, billingparams)

		// Done

		return newVM, nil
	}
}
