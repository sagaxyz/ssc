package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/sagaxyz/ssc/x/billing/types"
)

type (
	Keeper struct {
		cdc            codec.BinaryCodec
		storeKey       storetypes.StoreKey
		paramstore     paramtypes.Subspace
		bankkeeper     types.BankKeeper
		escrowkeeper   types.EscrowKeeper
		accountkeeper  types.AccountKeeper
		stakingkeeper  types.StakingKeeper
		chainletkeeper types.ChainletKeeper
		epochskeeper   types.EpochsKeeper
		authority      string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	bankkeeper types.BankKeeper,
	escrowkeeper types.EscrowKeeper,
	accountkeeper types.AccountKeeper,
	stakingkeeper types.StakingKeeper,
	chainletkeeper types.ChainletKeeper,
	epochskeeper types.EpochsKeeper,
	authority string,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		paramstore:     ps,
		bankkeeper:     bankkeeper,
		escrowkeeper:   escrowkeeper,
		accountkeeper:  accountkeeper,
		stakingkeeper:  stakingkeeper,
		chainletkeeper: chainletkeeper,
		epochskeeper:   epochskeeper,
		authority:      authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) UpdateKeeper(newKeeper interface{}) {
	if newk, ok := newKeeper.(types.ChainletKeeper); ok {
		k.chainletkeeper = newk
	} else if newk, ok := newKeeper.(types.EscrowKeeper); ok {
		k.escrowkeeper = newk
	} else if newk, ok := newKeeper.(types.EpochsKeeper); ok {
		k.epochskeeper = newk
	}
}

func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) SetPlatformValidators(ctx sdk.Context, vals []string) error {
	params := k.GetParams(ctx)
	params.PlatformValidators = vals
	k.SetParams(ctx, params)
	return nil
}

func (k Keeper) GetPlatformValidators(ctx sdk.Context) []string {
	params := k.GetParams(ctx)
	return params.PlatformValidators
}
