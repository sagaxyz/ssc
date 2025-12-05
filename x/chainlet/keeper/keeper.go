package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

type Keeper struct {
	cdc               codec.Codec
	storeKey          storetypes.StoreKey
	paramstore        paramtypes.Subspace
	providerMsgServer types.ProviderMsgServer
	stakingKeeper     types.StakingKeeper
	clientKeeper      types.ClientKeeper
	channelKeeper     types.ChannelKeeper
	connectionKeeper  types.ConnectionKeeper
	billingKeeper     types.BillingKeeper
	providerKeeper    types.ProviderKeeper
	escrowKeeper      types.EscrowKeeper
	aclKeeper         types.AclKeeper

	stackVersions      map[string]*versions.Versions // display name => version tree
	stackVersionParams map[string]map[string]types.ChainletStackParams
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	providerMsgServer types.ProviderMsgServer,
	stakingKeeper types.StakingKeeper,
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
		cdc:               cdc,
		storeKey:          storeKey,
		paramstore:        ps,
		providerMsgServer: providerMsgServer,
		stakingKeeper:     stakingKeeper,
		clientKeeper:      clientKeeper,
		channelKeeper:     channelKeeper,
		connectionKeeper:  connectionKeeper,
		billingKeeper:     billingKeeper,
		providerKeeper:    providerKeeper,
		escrowKeeper:      escrowKeeper,
		aclKeeper:         aclKeeper,
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) StackVersions(stackName string) *versions.Versions {
	return k.stackVersions[stackName]
}

// GetPort returns the portID for the IBC app module. Used in ExportGenesis
func (k *Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	portBytes := store.Get(types.PortKey)
	if portBytes == nil {
		return ""
	}
	return string(portBytes)
}

// SetPort sets the portID for the IBC app module. Used in InitGenesis
func (k *Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PortKey, []byte(portID))
}
