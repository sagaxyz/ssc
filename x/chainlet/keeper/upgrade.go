package keeper

import (
	"errors"
	"fmt"

	cosmossdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	sdkchainlettypes "github.com/sagaxyz/saga-sdk/x/chainlet/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

func UpgradePlanName(from, to string) (plan string, err error) {
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
	clientID, consumerRegistered := k.providerKeeper.GetConsumerClientId(ctx, chainlet.ConsumerId)
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

	lh := k.clientKeeper.GetClientLatestHeight(ctx, clientID)
	clientRevisionHeight := lh.GetRevisionHeight()
	clientRevisionNumber := lh.GetRevisionNumber()

	// Create the IBC packet
	upgradeHeight := clientRevisionHeight + heightDelta
	planName, err := UpgradePlanName(chainlet.ChainletStackVersion, newVersion)
	if err != nil {
		return
	}
	packetData := sdkchainlettypes.CreateUpgradePacketData{
		ChainId: chainlet.ChainId,
		Name:    planName,
		Height:  upgradeHeight,
		Info:    "Upgrade created by the provider chain",
	}
	err = packetData.ValidateBasic()
	if err != nil {
		return
	}

	// Timeout
	p := k.GetParams(ctx)
	var timeoutTimestamp uint64
	if p.UpgradeTimeoutTime > 0 {
		un := ctx.BlockTime().Add(p.UpgradeTimeoutTime).UnixNano()
		if un < 0 {
			err = errors.New("timeout negative")
			return
		}
		timeoutTimestamp = uint64(un)
	}
	var timeoutHeight clienttypes.Height
	if p.UpgradeTimeoutHeight > 0 {
		timeoutHeight = clienttypes.Height{
			RevisionNumber: clientRevisionNumber,
			RevisionHeight: clientRevisionHeight + p.UpgradeTimeoutHeight,
		}
	}

	k.Logger(ctx).Debug(fmt.Sprintf("sending packet to chainlet %s to create an upgrade to version %s\n", chainlet.ChainId, chainlet.ChainletStackVersion))
	_, err = k.TransmitCreateUpgradePacket(ctx, packetData, sdkchainlettypes.PortID, channelID, timeoutHeight, timeoutTimestamp)
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

// Verifies channel matches the client ID of a consumer
func (k Keeper) verifyChannel(ctx sdk.Context, clientID string, channelID string) error {
	channel, found := k.channelKeeper.GetChannel(ctx, sdkchainlettypes.PortID, channelID)
	if !found {
		return fmt.Errorf("channel %s not found (port: %s)", channelID, sdkchainlettypes.PortID)
	}
	if len(channel.ConnectionHops) == 0 {
		return fmt.Errorf("no connection hops for channel %s", channelID)
	}
	connectionID := channel.ConnectionHops[0]
	connection, found := k.connectionKeeper.GetConnection(ctx, connectionID)
	if !found {
		return fmt.Errorf("connection not found: %s", connectionID)
	}
	if connection.ClientId != clientID {
		return fmt.Errorf("client ID of the provided channel does not match consumer client id (%s != %s)", connection.ClientId, clientID)
	}

	return nil
}

func (k Keeper) sendCancelUpgradePlan(ctx sdk.Context, chainlet *types.Chainlet, channelID string) (err error) {
	if chainlet.Upgrade == nil {
		panic("no upgrade")
	}
	// Get consumer client id
	clientID, consumerRegistered := k.providerKeeper.GetConsumerClientId(ctx, chainlet.ConsumerId)
	if !consumerRegistered {
		err = errors.New("consumer not registered yet")
		return
	}
	err = k.verifyChannel(ctx, clientID, channelID)
	if err != nil {
		return
	}

	// Create the IBC packet
	planName, err := UpgradePlanName(chainlet.ChainletStackVersion, chainlet.Upgrade.Version)
	if err != nil {
		return
	}
	packetData := sdkchainlettypes.CancelUpgradePacketData{
		ChainId: chainlet.ChainId,
		Plan:    planName,
	}
	err = packetData.ValidateBasic()
	if err != nil {
		return
	}

	// Timeout
	p := k.GetParams(ctx)
	var timeoutTimestamp uint64
	if p.UpgradeTimeoutTime > 0 {
		un := ctx.BlockTime().Add(p.UpgradeTimeoutTime).UnixNano()
		if un < 0 {
			err = errors.New("timeout negative")
			return
		}
		timeoutTimestamp = uint64(un)
	}
	var timeoutHeight clienttypes.Height
	lh := k.clientKeeper.GetClientLatestHeight(ctx, clientID)
	if p.UpgradeTimeoutHeight > 0 {
		timeoutHeight = clienttypes.Height{
			RevisionNumber: lh.GetRevisionNumber(),
			RevisionHeight: lh.GetRevisionHeight() + p.UpgradeTimeoutHeight,
		}
	}

	k.Logger(ctx).Debug(fmt.Sprintf("sending packet to chainlet %s to cancel upgrade to version %s\n", chainlet.ChainId, chainlet.ChainletStackVersion))
	_, err = k.TransmitCancelUpgradePacket(ctx, packetData, sdkchainlettypes.PortID, channelID, timeoutHeight, timeoutTimestamp)
	if err != nil {
		return
	}

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
