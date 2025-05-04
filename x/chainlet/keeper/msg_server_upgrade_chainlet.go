package keeper

import (
	"context"
	"errors"
	"fmt"
	"time"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/gogoproto/proto"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	"golang.org/x/exp/slices"

	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

var icaOwner = authtypes.NewModuleAddress(types.ModuleName).String()

func (k msgServer) UpgradeChainlet(goCtx context.Context, msg *types.MsgUpgradeChainlet) (*types.MsgUpgradeChainletResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := msg.ValidateBasic()
	if err != nil {
		return &types.MsgUpgradeChainletResponse{}, err
	}

	ogChainlet, err := k.Chainlet(ctx, msg.ChainId)
	if err != nil {
		return &types.MsgUpgradeChainletResponse{}, err
	}

	if !slices.Contains(ogChainlet.Maintainers, msg.Creator) {
		return nil, fmt.Errorf("address %s not whitelisted for creating or updating chainlet stacks", msg.Creator)
	}
	majorUpgrade, err := versions.CheckUpgrade(ogChainlet.ChainletStackVersion, msg.StackVersion)
	if err != nil {
		return nil, err
	}
	if majorUpgrade {
		err = k.sendUpgradePlan(ctx, &ogChainlet, ogChainlet.ChainletStackVersion, msg.StackVersion)
		if err != nil {
			return nil, fmt.Errorf("error sending upgrade: %s", err)
		}

		return &types.MsgUpgradeChainletResponse{}, nil
	}

	err = k.UpgradeChainletStackVersion(ctx, msg.ChainId, msg.StackVersion)
	if err != nil {
		return nil, fmt.Errorf("error while updating chainlet: %s", err)
	}

	return &types.MsgUpgradeChainletResponse{}, ctx.EventManager().EmitTypedEvent(&types.EventUpdateChainlet{
		ChainId:      msg.ChainId,
		StackVersion: msg.StackVersion,
	})
}

func (k Keeper) sendUpgradePlan(ctx context.Context, chainlet *types.Chainlet, versionFrom, versionTo string) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	portID, err := icatypes.NewControllerPortID(icaOwner)
	if err != nil {
		return err
	}

	// Get consumer client id
	clientID, consumerRegistered := k.providerKeeper.GetConsumerClientId(sdkCtx, chainlet.ChainId)
	if !consumerRegistered {
		return errors.New("consumer not registered yet")
	}
	// Get consumer connection id
	ccvChannelID, found := k.providerKeeper.GetChainToChannel(sdkCtx, chainlet.ChainId)
	if !found{
		return errors.New("consumer channel not found")
	}
    ccvChannel, found := k.ibcKeeper.ChannelKeeper.GetChannel(sdkCtx, portID, ccvChannelID)
	if !found {
		return errors.New("consumer channel not found")
	}
	connectionID := ccvChannel.GetConnectionHops()[0] //TODO check len

	// Check ICA channel
	_, open := k.icaKeeper.GetOpenActiveChannel(sdkCtx, connectionID, portID)
	if !open {
		return fmt.Errorf("channel for connection %s and port %s not open", connectionID, portID)
	}

	// Create a MsgSoftwareUpgrade message
	clientState, ex := k.ibcKeeper.ClientKeeper.GetClientState(sdkCtx, clientID)
	if !ex {
		return fmt.Errorf("client state missing for client ID '%s'", clientID)
	}
	upgradeHeight := int64(clientState.GetLatestHeight().GetRevisionHeight()) + 1337 //TODO module param instead of a constant
	msg := upgradetypes.MsgSoftwareUpgrade{
		Authority: icaOwner,
		Plan: upgradetypes.Plan{
			Name:   "0.10-to-1", //TODO generate automatically
			Height: upgradeHeight,
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
	timeout := sdkCtx.BlockTime().Add(24 * time.Hour).UnixNano()
	//TODO do not use SendTx
	_, err = k.icaKeeper.SendTx(sdkCtx, nil, connectionID, portID, packetData, uint64(timeout))
	if err != nil {
		return err
	}

	// Mark the chainlet as being upgraded
	err = k.SetUpgrading(sdkCtx, chainlet, versionTo, upgradeHeight)
	if err != nil {
		return fmt.Errorf("error while updating chainlet: %w", err)
	}

	return nil
}
