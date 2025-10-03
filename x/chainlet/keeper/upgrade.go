package keeper

import (
	"errors"
	"fmt"

	cosmossdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ccvtypes "github.com/cosmos/interchain-security/v7/x/ccv/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

//nolint:unused
func (k *Keeper) getConsumerConnectionIDs(ctx sdk.Context, consumerID string) (controllerConnectionID, hostConnectionID string, err error) {
	// Get controller/local connection ID
	ccvChannelID, found := k.providerKeeper.GetConsumerIdToChannelId(ctx, consumerID)
	if !found {
		err = fmt.Errorf("channel ID for consumer ID %s not found", consumerID)
		return
	}
	ccvChannel, found := k.channelKeeper.GetChannel(ctx, ccvtypes.ProviderPortID, ccvChannelID)
	if !found {
		err = fmt.Errorf("channel %s for consumer %s not found", ccvChannelID, consumerID)
		return
	}
	if len(ccvChannel.ConnectionHops) == 0 {
		err = fmt.Errorf("no connections for channel %s", ccvChannelID)
		return
	}
	controllerConnectionID = ccvChannel.ConnectionHops[0]

	// Get host/counterparty connection ID
	connection, found := k.connectionKeeper.GetConnection(ctx, controllerConnectionID)
	if !found {
		err = fmt.Errorf("connection %s for consumer ID %s not found", controllerConnectionID, consumerID)
		return
	}
	hostConnectionID = connection.Counterparty.ConnectionId
	return
}

//nolint:unused
func upgradePlanName(from, to string) (plan string, err error) {
	major, minor, _, _, err := versions.Parse(from)
	if err != nil {
		return
	}
	if major == 0 {
		plan = fmt.Sprintf("0.%d-to", minor)
	} else {
		plan = fmt.Sprintf("%d-to", major)
	}

	major, minor, _, _, err = versions.Parse(to)
	if err != nil {
		return
	}
	if major == 0 {
		plan = fmt.Sprintf("%s-0.%d", plan, minor)
	} else {
		plan = fmt.Sprintf("%s-%d", plan, major)
	}

	return
}
func (k Keeper) sendUpgradePlan(ctx sdk.Context, chainlet *types.Chainlet, newVersion string, heightDelta uint64) (height uint64, err error) {
	err = errors.New("not implemented")
	return
}

//nolint:unused
func (k *Keeper) setUpgrading(ctx sdk.Context, chainlet *types.Chainlet, version string, height uint64) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletKey)

	key := []byte(chainlet.ChainId)
	if !store.Has(key) {
		return cosmossdkerrors.Wrapf(types.ErrInvalidChainId, "chainlet with chainId %s not found", chainlet.ChainId)
	}

	avail, err := k.chainletStackVersionAvailable(ctx, chainlet.ChainletStackName, version)
	if err != nil {
		return cosmossdkerrors.Wrapf(types.ErrInvalidChainletStack, "cannot upgrade to stack %s version %s: %s", chainlet.ChainletStackName, version, err)
	}
	if !avail {
		return cosmossdkerrors.Wrapf(types.ErrInvalidChainletStack, "stack %s version %s not available", chainlet.ChainletStackName, version)
	}

	chainlet.Upgrade = &types.Upgrade{
		Height:  height,
		Version: version,
	}

	updatedValue := k.cdc.MustMarshal(chainlet)
	store.Set(key, updatedValue)

	store = prefix.NewStore(ctx.KVStore(k.storeKey), types.UpgradingChainletsKey)
	store.Set([]byte(chainlet.ChainId), k.cdc.MustMarshal(&types.UpgradingChainlet{}))
	return nil
}

//nolint:unused
func (k *Keeper) finishUpgrading(ctx sdk.Context, chainlet *types.Chainlet) error {
	if chainlet.Upgrade == nil {
		return fmt.Errorf("chainlet %s is not being upgraded", chainlet.ChainId)
	}
	chainlet.ChainletStackVersion = chainlet.Upgrade.Version
	chainlet.Upgrade = nil

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletKey)
	store.Set([]byte(chainlet.ChainId), k.cdc.MustMarshal(chainlet))
	store = prefix.NewStore(ctx.KVStore(k.storeKey), types.UpgradingChainletsKey)
	store.Delete([]byte(chainlet.ChainId))
	return nil
}

//nolint:unused
func (k *Keeper) cancelUpgrading(ctx sdk.Context, chainlet *types.Chainlet) {
	chainlet.Upgrade = nil

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletKey)
	store.Set([]byte(chainlet.ChainId), k.cdc.MustMarshal(chainlet))
	store = prefix.NewStore(ctx.KVStore(k.storeKey), types.UpgradingChainletsKey)
	store.Delete([]byte(chainlet.ChainId))
}
