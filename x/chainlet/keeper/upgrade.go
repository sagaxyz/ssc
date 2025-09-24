package keeper

import (
	"errors"
	"fmt"
	"time"

	cosmossdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	ccvtypes "github.com/cosmos/interchain-security/v7/x/ccv/types"
	sdkchainlettypes "github.com/sagaxyz/saga-sdk/x/chainlet/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

func (k *Keeper) getConsumerConnectionIDs(ctx sdk.Context, chainID string) (controllerConnectionID, hostConnectionID string, err error) {
	// Get controller/local connection ID
	ccvChannelID, found := k.providerKeeper.GetConsumerIdToChannelId(ctx, chainID)
	if !found {
		err = fmt.Errorf("channel ID for consumer %s not found", chainID)
		return
	}
	ccvChannel, found := k.channelKeeper.GetChannel(ctx, ccvtypes.ProviderPortID, ccvChannelID)
	if !found {
		err = fmt.Errorf("channel %s for consumer %s not found", ccvChannelID, chainID)
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
		err = fmt.Errorf("connection %s for consumer %s not found", controllerConnectionID, chainID)
		return
	}
	hostConnectionID = connection.Counterparty.ConnectionId
	return
}

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
func (k Keeper) sendUpgradePlan(ctx sdk.Context, chainlet *types.Chainlet, newVersion string, heightDelta uint64, channelID string) (height uint64, err error) {
	// Get consumer client id
	clientID, consumerRegistered := k.providerKeeper.GetConsumerClientId(ctx, chainlet.ChainId)
	if !consumerRegistered {
		err = errors.New("consumer not registered yet")
		return
	}

	// Verify channel corresponds to the correct client
	channel, found := k.channelKeeper.GetChannel(ctx, sdkchainlettypes.PortID, channelID)
	if !found {
		err = fmt.Errorf("channel %s not found (port: %s)", channelID, sdkchainlettypes.PortID)
		return
	}
	if len(channel.ConnectionHops) == 0 {
		err = fmt.Errorf("no connection hops for channel %s", channelID)
		return
	}
	connectionID := channel.ConnectionHops[0]
	connection, found := k.connectionKeeper.GetConnection(ctx, connectionID)
	if !found {
		err = fmt.Errorf("connection not found: %s", connectionID)
		return
	}
	if connection.ClientId != clientID {
		err = fmt.Errorf("client ID of the provided channel does not match consumer client id (%s != %s)", connection.ClientId, clientID)
		return
	}

	clientRevisionHeight := k.clientKeeper.GetClientLatestHeight(ctx, clientID).GetRevisionHeight()
	clientRevisionNumber := k.clientKeeper.GetClientLatestHeight(ctx, clientID).GetRevisionNumber()

	// Create the IBC packet
	upgradeHeight := clientRevisionHeight + heightDelta
	planName, err := upgradePlanName(chainlet.ChainletStackVersion, newVersion)
	if err != nil {
		return
	}
	packetData := sdkchainlettypes.CreateUpgradePacketData{
		Name:   planName,
		Height: upgradeHeight,
		Info:   "Upgrade created by the provider chain",
	}
	err = packetData.ValidateBasic()
	if err != nil {
		return
	}

	// Timeout
	//p := k.GetParams(ctx)
	TimeoutHeight := uint64(300) //TODO p
	TimeoutTime := 2 * time.Hour //TODO p
	timeoutHeight := clienttypes.Height{
		RevisionNumber: clientRevisionNumber, 
		RevisionHeight: clientRevisionHeight + TimeoutHeight,
	}
	timeoutTimestamp := uint64(ctx.BlockTime().Add(TimeoutTime).UnixNano())

	//TODO remove
	timeoutHeight = clienttypes.Height{
		RevisionNumber: 1,
		RevisionHeight: 10000000000,
	}

	fmt.Printf("XXX sending packet %+v to %s with port %s, timeout %s and %s\n", packetData, channelID, types.PortID, timeoutHeight, timeoutTimestamp)
	_, err = k.TransmitCreateUpgradePacket(ctx, packetData, types.PortID, channelID, timeoutHeight, timeoutTimestamp)
	if err != nil {
		return
	}

	// Mark the chainlet as being upgraded
	err = k.setUpgrading(ctx, chainlet, newVersion, upgradeHeight)
	if err != nil {
		err = fmt.Errorf("error while updating chainlet: %w", err)
		return
	}

	height = upgradeHeight
	return
}

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

func (k *Keeper) cancelUpgrading(ctx sdk.Context, chainlet *types.Chainlet) {
	chainlet.Upgrade = nil

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletKey)
	store.Set([]byte(chainlet.ChainId), k.cdc.MustMarshal(chainlet))
	store = prefix.NewStore(ctx.KVStore(k.storeKey), types.UpgradingChainletsKey)
	store.Delete([]byte(chainlet.ChainId))
}
