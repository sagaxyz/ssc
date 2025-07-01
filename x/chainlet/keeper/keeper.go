package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

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

	ibcKeeperFn  func() *ibckeeper.Keeper
	scopedKeeper exported.ScopedKeeper

	stackVersions map[string]*versions.Versions // display name => version tree
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	ibcKeeperFn func() *ibckeeper.Keeper,
	scopedKeeper exported.ScopedKeeper,
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
		scopedKeeper:     scopedKeeper,
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

// ChanCloseInit defines a wrapper function for the channel Keeper's function.
func (k Keeper) ChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	capName := host.ChannelCapabilityPath(portID, channelID)
	chanCap, ok := k.scopedKeeper.GetCapability(ctx, capName)
	if !ok {
		return errorsmod.Wrapf(channeltypes.ErrChannelCapabilityNotFound, "could not retrieve channel capability at: %s", capName)
	}
	return k.ibcKeeperFn().ChannelKeeper.ChanCloseInit(ctx, portID, channelID, chanCap)
}

// IsBound checks if the IBC app module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the port Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	cap := k.ibcKeeperFn().PortKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}

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

// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.scopedKeeper.AuthenticateCapability(ctx, cap, name)
}

// ClaimCapability allows the IBC app module to claim a capability that core IBC
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

// ScopedKeeper returns the ScopedKeeper
func (k Keeper) ScopedKeeper() exported.ScopedKeeper {
	return k.scopedKeeper
}
