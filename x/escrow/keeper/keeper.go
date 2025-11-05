package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/sagaxyz/ssc/x/escrow/types"
)

type (
	Keeper struct {
		cdc            codec.BinaryCodec
		storeKey       storetypes.StoreKey
		paramstore     paramtypes.Subspace
		bankKeeper     types.BankKeeper
		billingKeeper  types.BillingKeeper
		chainletKeeper types.ChainletKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	bk types.BankKeeper,
	billingKeeper types.BillingKeeper,
	chainletKeeper types.ChainletKeeper,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		paramstore:     ps,
		bankKeeper:     bk,
		billingKeeper:  billingKeeper,
		chainletKeeper: chainletKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) UpdateKeeper(newKeeper interface{}) {
	switch v := newKeeper.(type) {
	case types.BillingKeeper:
		k.billingKeeper = v
	case types.ChainletKeeper:
		k.chainletKeeper = v
	}
}

func (k Keeper) GetSupportedDenoms(ctx sdk.Context) []string {
	params := k.GetParams(ctx)
	return params.SupportedDenoms
}
