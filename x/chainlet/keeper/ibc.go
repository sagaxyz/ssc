package keeper

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	chainlettypes "github.com/sagaxyz/saga-sdk/x/chainlet/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
)

// TransmitCreateUpgradePacket transmits the packet over IBC with the specified source port and source channel
func (k Keeper) TransmitCreateUpgradePacket(
	ctx sdk.Context,
	packetData chainlettypes.CreateUpgradePacketData,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) (uint64, error) {
	packetBytes, err := packetData.GetBytes()
	if err != nil {
		return 0, errorsmod.Wrapf(sdkerrors.ErrJSONMarshal, "cannot marshal the packet: %s", err)
	}

	return k.ibcKeeperFn().ChannelKeeper.SendPacket(ctx, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, packetBytes)
}

// OnAcknowledgementCreateUpgradePacket responds to the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementCreateUpgradePacket(ctx sdk.Context, packet channeltypes.Packet, data chainlettypes.CreateUpgradePacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		chainlet, err := k.Chainlet(ctx, data.ChainId)
		if err != nil {
			return err
		}
		if chainlet.Upgrade == nil {
			return nil
		}

		// Verify channel maches chain ID
		clientID, consumerRegistered := k.providerKeeper.GetConsumerClientId(ctx, chainlet.ConsumerId)
		if !consumerRegistered {
			return errors.New("consumer not registered yet")
		}
		err = k.verifyChannel(ctx, clientID, packet.SourceChannel)
		if err != nil {
			return err
		}

		// Cancel if the upgrade plan matches the current upgrade
		planName, err := upgradePlanName(chainlet.ChainletStackVersion, chainlet.Upgrade.Version)
		if err != nil {
			return err
		}
		if data.Name == planName {
			k.cancelUpgrading(ctx, &chainlet)
			ctx.Logger().Info(fmt.Sprintf("cancelled upgrade %s for chainlet %s: error ack: %s\n", planName, chainlet.ChainId, ack))
		}

		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck chainlettypes.CreateUpgradePacketAck

		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}

// OnTimeoutCreateUpgradePacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutCreateUpgradePacket(ctx sdk.Context, packet channeltypes.Packet, data chainlettypes.CreateUpgradePacketData) error {
	chainlet, err := k.Chainlet(ctx, data.ChainId)
	if err != nil {
		return err
	}
	if chainlet.Upgrade == nil {
		return nil
	}

	// Verify channel maches chain ID
	clientID, consumerRegistered := k.providerKeeper.GetConsumerClientId(ctx, chainlet.ConsumerId)
	if !consumerRegistered {
		return errors.New("consumer not registered yet")
	}
	err = k.verifyChannel(ctx, clientID, packet.SourceChannel)
	if err != nil {
		return err
	}

	// Cancel if the upgrade plan matches the current upgrade
	planName, err := upgradePlanName(chainlet.ChainletStackVersion, chainlet.Upgrade.Version)
	if err != nil {
		return err
	}
	if data.Name == planName {
		k.cancelUpgrading(ctx, &chainlet)
		ctx.Logger().Info(fmt.Sprintf("cancelled upgrade %s for chainlet %s: timed out\n", planName, chainlet.ChainId))
	}
	return nil
}

// TransmitCancelUpgradePacket transmits the packet over IBC with the specified source port and source channel
func (k Keeper) TransmitCancelUpgradePacket(
	ctx sdk.Context,
	packetData chainlettypes.CancelUpgradePacketData,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) (uint64, error) {
	packetBytes, err := packetData.GetBytes()
	if err != nil {
		return 0, errorsmod.Wrapf(sdkerrors.ErrJSONMarshal, "cannot marshal the packet: %s", err)
	}

	return k.ibcKeeperFn().ChannelKeeper.SendPacket(ctx, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, packetBytes)
}

// OnAcknowledgementCancelUpgradePacket responds to the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementCancelUpgradePacket(ctx sdk.Context, packet channeltypes.Packet, data chainlettypes.CancelUpgradePacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		ctx.Logger().Error(fmt.Sprintf("failed to cancel upgrade for chainlet %s: error ack: %s", data.ChainId, ack))
		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck chainlettypes.CancelUpgradePacketAck

		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		chainlet, err := k.Chainlet(ctx, data.ChainId)
		if err != nil {
			return err
		}
		if chainlet.Upgrade == nil {
			ctx.Logger().Error(fmt.Sprintf("failed to cancel upgrade for chainlet %s: no upgrade", chainlet.ChainId))
			return nil
		}

		// Verify channel maches chain ID
		clientID, consumerRegistered := k.providerKeeper.GetConsumerClientId(ctx, chainlet.ConsumerId)
		if !consumerRegistered {
			return errors.New("consumer not registered yet")
		}
		err = k.verifyChannel(ctx, clientID, packet.SourceChannel)
		if err != nil {
			return err
		}

		// Cancel if the upgrade plan matches the current upgrade
		planName, err := upgradePlanName(chainlet.ChainletStackVersion, chainlet.Upgrade.Version)
		if err != nil {
			return err
		}
		if data.Plan == planName {
			k.cancelUpgrading(ctx, &chainlet)
			ctx.Logger().Info(fmt.Sprintf("cancelled upgrade %s for chainlet %s\n", planName, chainlet.ChainId))
		} else {
			ctx.Logger().Error(fmt.Sprintf("failed to cancel upgrade for chainlet %s: plan does not match (%s != %s)\n", chainlet.ChainId, data.Plan, planName))
		}

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}

// OnTimeoutCancelUpgradePacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutCancelUpgradePacket(ctx sdk.Context, packet channeltypes.Packet, data chainlettypes.CancelUpgradePacketData) error {
	ctx.Logger().Error(fmt.Sprintf("failed to cancel upgrade for chainlet %s: timed out", data.ChainId))
	return nil
}

// OnRecvConfirmUpgradePacket processes packet reception
func (k Keeper) OnRecvConfirmUpgradePacket(ctx sdk.Context, packet channeltypes.Packet, data chainlettypes.ConfirmUpgradePacketData) (packetAck chainlettypes.ConfirmUpgradePacketAck, err error) {
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	chainlet, err := k.Chainlet(ctx, data.ChainId)
	if err != nil {
		return
	}
	if chainlet.Upgrade == nil {
		err = fmt.Errorf("chainlet %s is not being upgraded", data.ChainId)
		return
	}

	// Verify channel maches chain ID
	clientID, consumerRegistered := k.providerKeeper.GetConsumerClientId(ctx, chainlet.ConsumerId)
	if !consumerRegistered {
		err = errors.New("consumer not registered yet")
		return
	}
	err = k.verifyChannel(ctx, clientID, packet.SourceChannel)
	if err != nil {
		return
	}

	ctx.Logger().Info(fmt.Sprintf("finished upgrading chainlet %s to version %s\n", chainlet.ChainId, chainlet.Upgrade.Version))
	err = k.finishUpgrading(ctx, &chainlet)
	if err != nil {
		return
	}

	return packetAck, nil
}
