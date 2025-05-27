package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

type Keeper struct {
	cdc            codec.Codec
	storeKey       storetypes.StoreKey
	paramstore     paramtypes.Subspace
	stakingKeeper types.StakingKeeper
	billingKeeper  types.BillingKeeper
	providerKeeper types.ProviderKeeper
	escrowKeeper   types.EscrowKeeper
	aclKeeper      types.AclKeeper
	icaKeeper      icacontrollerkeeper.Keeper //TODO expected interface
	msgRouter      icatypes.MessageRouter
	ibcKeeper      *ibckeeper.Keeper //TODO expected interface

	stackVersions map[string]*versions.Versions // display name => version tree
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	stakingKeeper types.StakingKeeper,
	ica icacontrollerkeeper.Keeper,
	msgRouter icatypes.MessageRouter,
	ibcKeeper *ibckeeper.Keeper,
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
		cdc:            cdc,
		storeKey:       storeKey,
		paramstore:     ps,
		stakingKeeper: stakingKeeper,
		icaKeeper:      ica,
		msgRouter:      msgRouter,
		ibcKeeper:      ibcKeeper,
		billingKeeper:  billingKeeper,
		providerKeeper: providerKeeper,
		escrowKeeper:   escrowKeeper,
		aclKeeper:      aclKeeper,
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) StackVersions(stackName string) *versions.Versions {
	return k.stackVersions[stackName]
}
