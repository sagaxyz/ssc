package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	icatypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/types"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"

	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

type Keeper struct {
	cdc              codec.Codec
	storeKey         storetypes.StoreKey
	paramstore       paramtypes.Subspace
	stakingKeeper    types.StakingKeeper
	icaKeeper        types.ICAKeeper
	msgRouter        icatypes.MessageRouter
	clientKeeper     types.ClientKeeper
	channelKeeper    types.ChannelKeeper
	connectionKeeper types.ConnectionKeeper
	billingKeeper    types.BillingKeeper
	providerKeeper   types.ProviderKeeper
	escrowKeeper     types.EscrowKeeper
	aclKeeper        types.AclKeeper

	ibcKeeperFn func() *ibckeeper.Keeper

	stackVersions map[string]*versions.Versions // display name => version tree
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	ibcKeeperFn func() *ibckeeper.Keeper,
	stakingKeeper types.StakingKeeper,
	icaKeeper types.ICAKeeper,
	msgRouter icatypes.MessageRouter,
	clientKeeper types.ClientKeeper,
	channelKeeper types.ChannelKeeper,
	connectionKeeper types.ConnectionKeeper,
	providerKeeper types.ProviderKeeper,
	billingKeeper types.BillingKeeper,
	escrowKeeper types.EscrowKeeper,
	aclKeeper types.AclKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		paramstore:       ps,
		ibcKeeperFn:      ibcKeeperFn,
		stakingKeeper:    stakingKeeper,
		icaKeeper:        icaKeeper,
		msgRouter:        msgRouter,
		clientKeeper:     clientKeeper,
		channelKeeper:    channelKeeper,
		connectionKeeper: connectionKeeper,
		billingKeeper:    billingKeeper,
		providerKeeper:   providerKeeper,
		escrowKeeper:     escrowKeeper,
		aclKeeper:        aclKeeper,
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) StackVersions(stackName string) *versions.Versions {
	return k.stackVersions[stackName]
}

// ----------------------------------------------------------------------------
// IBC Keeper Logic
// ----------------------------------------------------------------------------

// GetPort returns the portID for the IBC app module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.PortKey))
}

// SetPort sets the portID for the IBC app module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PortKey, []byte(portID))
}
