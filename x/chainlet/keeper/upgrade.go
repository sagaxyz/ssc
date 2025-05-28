package keeper

import (
	"errors"
	"fmt"
	"time"

	cosmossdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/gogoproto/proto"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ccvtypes "github.com/cosmos/interchain-security/v5/x/ccv/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

func (k *Keeper) getConsumerConnectionIDs(ctx sdk.Context, chainID string) (controllerConnectionID, hostConnectionID string, err error) {
	// Get controller/local connection ID
	ccvChannelID, found := k.providerKeeper.GetChainToChannel(ctx, chainID)
	if !found {
		err = fmt.Errorf("channel ID for consumer %s not found", chainID)
		return
	}
	ccvChannel, found := k.channelKeeper.GetChannel(ctx, ccvtypes.ProviderPortID, ccvChannelID)
	if !found {
		err = fmt.Errorf("channel %s for consumer %s not found", ccvChannelID, chainID)
		return
	}
	if len(ccvChannel.GetConnectionHops()) == 0 {
		err = fmt.Errorf("no connections for channel %s", ccvChannelID)
		return
	}
	controllerConnectionID = ccvChannel.GetConnectionHops()[0]

	// Get host/counterparty connection ID
	connection, found := k.connectionKeeper.GetConnection(ctx, controllerConnectionID)
	if !found {
		err = fmt.Errorf("connection %s for consumer %s not found", controllerConnectionID, chainID)
		return
	}
	hostConnectionID = connection.Counterparty.ConnectionId
	return
}

func (k *Keeper) InitICA(ctx sdk.Context, chainID string) error {
	controllerConnectionID, hostConnectionID, err := k.getConsumerConnectionIDs(ctx, chainID)
	if err != nil {
		return err
	}

	icaOwner := sdk.AccAddress(address.Module(types.ModuleName)).String()
	portID, err := icatypes.NewControllerPortID(icaOwner)
	if err != nil {
		return err
	}

	// Check if the account exists
	_, found := k.icaKeeper.GetInterchainAccountAddress(ctx, controllerConnectionID, portID)
	if found {
		return nil
	}

	// Register the account
	ctx.Logger().Debug(fmt.Sprintf("registering ICA account for chain %s using connection %s and port %s", chainID, controllerConnectionID, portID))
	metadata := icatypes.NewMetadata(icatypes.Version, controllerConnectionID, hostConnectionID, icaOwner, icatypes.EncodingProtobuf, icatypes.TxTypeSDKMultiMsg)
	ver := string(icatypes.ModuleCdc.MustMarshalJSON(&metadata))
	icaMsg := icacontrollertypes.NewMsgRegisterInterchainAccount(controllerConnectionID, icaOwner, ver)
	handler := k.msgRouter.Handler(icaMsg)
	_, err = handler(ctx, icaMsg)
	if err != nil {
		return err
	}

	return nil
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
func (k Keeper) sendUpgradePlan(ctx sdk.Context, chainlet *types.Chainlet, newVersion string, heightDelta uint64) (height uint64, err error) {
	// Get consumer client id
	clientID, consumerRegistered := k.providerKeeper.GetConsumerClientId(ctx, chainlet.ChainId)
	if !consumerRegistered {
		err = errors.New("consumer not registered yet")
		return
	}
	// Get consumer connection id
	connectionID, _, err := k.getConsumerConnectionIDs(ctx, chainlet.ChainId)
	if err != nil {
		return
	}

	// Check the ICA account
	icaOwner := sdk.AccAddress(address.Module(types.ModuleName)).String()
	portID, err := icatypes.NewControllerPortID(icaOwner)
	if err != nil {
		return
	}

	addr, found := k.icaKeeper.GetInterchainAccountAddress(ctx, connectionID, portID)
	if !found {
		err = fmt.Errorf("missing ICA account for chainlet %s", chainlet.ChainId)
		return
	}

	// Check ICA channel
	_, open := k.icaKeeper.GetOpenActiveChannel(ctx, connectionID, portID)
	if !open {
		err = fmt.Errorf("channel for connection %s and port %s not open", connectionID, portID)
		return
	}

	// Create a MsgSoftwareUpgrade message
	clientState, ex := k.clientKeeper.GetClientState(ctx, clientID)
	if !ex {
		err = fmt.Errorf("client state missing for client ID '%s'", clientID)
		return
	}
	upgradeHeight := clientState.GetLatestHeight().GetRevisionHeight() + heightDelta
	planName, err := upgradePlanName(chainlet.ChainletStackVersion, newVersion)
	if err != nil {
		return
	}
	msg := upgradetypes.MsgSoftwareUpgrade{
		Authority: addr,
		Plan: upgradetypes.Plan{
			Name:   planName,
			Height: int64(upgradeHeight),
			Info:   "Upgrade created by the provider chain",
		},
	}

	// Send the message using ICA
	data, err := icatypes.SerializeCosmosTx(k.cdc, []proto.Message{&msg}, icatypes.EncodingProtobuf)
	if err != nil {
		return
	}
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}
	timeout := ctx.BlockTime().Add(24 * time.Hour).UnixNano() //TODO module param
	_, err = k.icaKeeper.SendTx(ctx, nil, connectionID, portID, packetData, uint64(timeout))
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

func (k *Keeper) HandleUpgradingChainlets(ctx sdk.Context) error {
	iterator := prefix.NewStore(ctx.KVStore(k.storeKey), types.UpgradingChainletsKey).Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		chainID := iterator.Key()
		fmt.Printf("XXX got upgrading chainlet: %s\n", chainID)

		chainlet, err := k.Chainlet(ctx, string(chainID))
		if err != nil {
			return err
		}
		if chainlet.Upgrade == nil {
			panic("no upgrade info")
		}

		clientId, ex := k.providerKeeper.GetConsumerClientId(ctx, chainlet.ChainId)
		if !ex {
			fmt.Printf("XXX upgrading chainlet: %s: not client\n", chainID)
			//TODO log
			continue
		}
		clientState, ex := k.clientKeeper.GetClientState(ctx, clientId)
		if !ex {
			return fmt.Errorf("client state missing for client ID '%s'", clientId)
		}

		//TODO check revision number
		height := clientState.GetLatestHeight().GetRevisionHeight()
		if height > chainlet.Upgrade.Height {
			// Chain failed to stop before the upgrade height => cancel the upgrade
			k.cancelUpgrading(ctx, &chainlet)
			fmt.Printf("XXX upgrading chainlet: %s: cancelled\n", chainID)
			continue
		}
		if height < chainlet.Upgrade.Height-1 {
			continue
		}

		fmt.Printf("XXX upgrading chainlet: %s: DONE\n", chainID)
		k.finishUpgrading(ctx, &chainlet)
	}

	return nil
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
