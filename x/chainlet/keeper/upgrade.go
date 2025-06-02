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

const hostConnectionID = "connection-0"

func (k *Keeper) getConsumerConnectionID(ctx sdk.Context, chainID string) (connectionID string, err error) {
	ccvChannelID, found := k.providerKeeper.GetChainToChannel(ctx, chainID)
	if !found {
		err = errors.New("consumer channel not found")
		return
	}
	ccvChannel, found := k.ibcKeeper.ChannelKeeper.GetChannel(ctx, ccvtypes.ProviderPortID, ccvChannelID)
	if !found {
		err = errors.New("consumer channel not found")
		return
	}
	if len(ccvChannel.GetConnectionHops()) == 0 {
		err = fmt.Errorf("no connections for channel %s", ccvChannel)
		return
	}

	connectionID = ccvChannel.GetConnectionHops()[0]
	return
}

func (k *Keeper) InitICA(ctx sdk.Context, chainID string) error {
	connectionID, err := k.getConsumerConnectionID(ctx, chainID)
	if err != nil {
		return err
	}

	icaOwner := sdk.AccAddress(address.Module(types.ModuleName)).String()
	portID, err := icatypes.NewControllerPortID(icaOwner)
	if err != nil {
		return err
	}

	// Check if the account exists
	_, found := k.icaKeeper.GetInterchainAccountAddress(ctx, connectionID, portID)
	if found {
		return nil
	}

	// Register the account
	ctx.Logger().Debug(fmt.Sprintf("registering ICA account for chain %s using connection %s and port %s", chainID, connectionID, portID))
	metadata := icatypes.NewMetadata(icatypes.Version, connectionID, hostConnectionID, icaOwner, icatypes.EncodingProtobuf, icatypes.TxTypeSDKMultiMsg)
	ver := string(icatypes.ModuleCdc.MustMarshalJSON(&metadata))
	icaMsg := icacontrollertypes.NewMsgRegisterInterchainAccount(connectionID, icaOwner, ver)
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
func (k Keeper) sendUpgradePlan(ctx sdk.Context, chainlet *types.Chainlet, versionFrom, versionTo string) error {
	// Get consumer client id
	clientID, consumerRegistered := k.providerKeeper.GetConsumerClientId(ctx, chainlet.ChainId)
	if !consumerRegistered {
		return errors.New("consumer not registered yet")
	}
	// Get consumer connection id
	connectionID, err := k.getConsumerConnectionID(ctx, chainlet.ChainId)
	if err != nil {
		return err
	}

	// Check the ICA account
	icaOwner := sdk.AccAddress(address.Module(types.ModuleName)).String()
	portID, err := icatypes.NewControllerPortID(icaOwner)
	if err != nil {
		return err
	}

	addr, found := k.icaKeeper.GetInterchainAccountAddress(ctx, connectionID, portID)
	if !found {
		return fmt.Errorf("missing ICA account for chainlet %s", chainlet.ChainId)
	}

	// Check ICA channel
	_, open := k.icaKeeper.GetOpenActiveChannel(ctx, connectionID, portID)
	if !open {
		return fmt.Errorf("channel for connection %s and port %s not open", connectionID, portID)
	}

	// Create a MsgSoftwareUpgrade message
	clientState, ex := k.ibcKeeper.ClientKeeper.GetClientState(ctx, clientID)
	if !ex {
		return fmt.Errorf("client state missing for client ID '%s'", clientID)
	}
	upgradeHeight := clientState.GetLatestHeight().GetRevisionHeight() + 30 //TODO module param minimum/default + msg option
	planName, err := upgradePlanName(versionFrom, versionTo)
	if err != nil {
		return err
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
		return err
	}
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}
	timeout := ctx.BlockTime().Add(24 * time.Hour).UnixNano() //TODO module param
	_, err = k.icaKeeper.SendTx(ctx, nil, connectionID, portID, packetData, uint64(timeout))
	if err != nil {
		return err
	}

	// Mark the chainlet as being upgraded
	err = k.setUpgrading(ctx, chainlet, versionTo, upgradeHeight)
	if err != nil {
		return fmt.Errorf("error while updating chainlet: %w", err)
	}

	return nil
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
		clientState, ex := k.ibcKeeper.ClientKeeper.GetClientState(ctx, clientId)
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
		if height < chainlet.Upgrade.Height - 1 {
			fmt.Printf("XXX upgrading chainlet: %s: not target height yet (%d < %d)\n", chainID, height, chainlet.Upgrade.Height)
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
		return cosmossdkerrors.Wrapf(types.ErrInvalidChainletStack, "stack %s version %s not available", chainlet.ChainletStackName, chainlet.ChainletStackVersion)
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

func (k *Keeper) finishUpgrading(ctx sdk.Context, chainlet *types.Chainlet) {
	chainlet.ChainletStackVersion = chainlet.Upgrade.Version
	chainlet.Upgrade = nil

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletKey)
	store.Set([]byte(chainlet.ChainId), k.cdc.MustMarshal(chainlet))
	store = prefix.NewStore(ctx.KVStore(k.storeKey), types.UpgradingChainletsKey)
	store.Delete([]byte(chainlet.ChainId))
}

func (k *Keeper) cancelUpgrading(ctx sdk.Context, chainlet *types.Chainlet) {
	chainlet.Upgrade = nil
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletKey)
	store.Set([]byte(chainlet.ChainId), k.cdc.MustMarshal(chainlet))

	store = prefix.NewStore(ctx.KVStore(k.storeKey), types.UpgradingChainletsKey)
	store.Delete([]byte(chainlet.ChainId))
}
