package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

type Keeper struct {
	cdc            codec.BinaryCodec
	storeKey       storetypes.StoreKey
	memKey         storetypes.StoreKey
	paramstore     paramtypes.Subspace
	billingKeeper  types.BillingKeeper
	providerKeeper types.ProviderKeeper
	escrowKeeper   types.EscrowKeeper
	dacKeeper      types.DacKeeper

	stackVersions map[string]*versions.Versions // display name => version tree
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	providerKeeper types.ProviderKeeper,
	billingKeeper types.BillingKeeper,
	escrowKeeper types.EscrowKeeper,
	dacKeeper types.DacKeeper,
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
		billingKeeper:  billingKeeper,
		providerKeeper: providerKeeper,
		escrowKeeper:   escrowKeeper,
		dacKeeper:      dacKeeper,
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) StackVersions(stackName string) *versions.Versions {
	return k.stackVersions[stackName]
}
