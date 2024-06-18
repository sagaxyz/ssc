package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/sagaxyz/ssc/x/billing/types"
)

type (
	Keeper struct {
		cdc            codec.BinaryCodec
		storeKey       storetypes.StoreKey
		memKey         storetypes.StoreKey
		paramstore     paramtypes.Subspace
		bankkeeper     types.BankKeeper
		escrowkeeper   types.EscrowKeeper
		accountkeeper  types.AccountKeeper
		stakingkeeper  types.StakingKeeper
		chainletkeeper types.ChainletKeeper
		epochskeeper   types.EpochsKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	bankkeeper types.BankKeeper,
	escrowkeeper types.EscrowKeeper,
	accountkeeper types.AccountKeeper,
	stakingkeeper types.StakingKeeper,
	chainletkeeper types.ChainletKeeper,
	epochskeeper types.EpochsKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		memKey:         memKey,
		paramstore:     ps,
		bankkeeper:     bankkeeper,
		escrowkeeper:   escrowkeeper,
		accountkeeper:  accountkeeper,
		stakingkeeper:  stakingkeeper,
		chainletkeeper: chainletkeeper,
		epochskeeper:   epochskeeper,
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
