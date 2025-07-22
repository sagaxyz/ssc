package keeper

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"

	"github.com/sagaxyz/saga-sdk/x/chainlet/types"
)

// TransmitConfirmUpgradePacket transmits the packet over IBC with the specified source port and source channel
func (k Keeper) TransmitConfirmUpgradePacket(
	ctx sdk.Context,
	packetData types.ConfirmUpgradePacketData,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) (uint64, error) {
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return 0, errorsmod.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	packetBytes, err := packetData.GetBytes()
	if err != nil {
		return 0, errorsmod.Wrapf(sdkerrors.ErrJSONMarshal, "cannot marshal the packet: %s", err)
	}

	return k.ibcKeeperFn().ChannelKeeper.SendPacket(ctx, channelCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, packetBytes)
}

// OnRecvConfirmUpgradePacket processes packet reception
func (k Keeper) OnRecvConfirmUpgradePacket(ctx sdk.Context, packet channeltypes.Packet, data types.ConfirmUpgradePacketData) (packetAck types.ConfirmUpgradePacketAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	fmt.Printf("XXX received packet %+v\n", data)

	chainlet, err := k.Chainlet(ctx, data.ChainId)
	if err != nil {
		fmt.Printf("XXX upgrading %s: %s\n", data.ChainId, err)
		return
	}
	if chainlet.Upgrade == nil {
		err = fmt.Errorf("chain %s is not being upgraded", data.ChainId)
		fmt.Printf("XXX upgrading %s: %s\n", data.ChainId, err)
		return
	}
	if data.Height != chainlet.Upgrade.Height - 1 {
		err = fmt.Errorf("unexpected upgrade height: %d != %d", data.Height, chainlet.Upgrade.Height)
		fmt.Printf("XXX upgrading %s: %s\n", data.ChainId, err)
		return
	}
	err = k.finishUpgrading(ctx, &chainlet)
	if err != nil {
		fmt.Printf("XXX upgrading %s: %s\n", data.ChainId, err)
		return
	}
	fmt.Printf("XXX upgrading %s: DONE\n")

	return packetAck, nil
}

// OnAcknowledgementConfirmUpgradePacket responds to the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementConfirmUpgradePacket(ctx sdk.Context, packet channeltypes.Packet, data types.ConfirmUpgradePacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:

		// TODO: failed acknowledgement logic
		_ = dispatchedAck.Error

		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.ConfirmUpgradePacketAck

		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		// TODO: successful acknowledgement logic

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}

// OnTimeoutConfirmUpgradePacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutConfirmUpgradePacket(ctx sdk.Context, packet channeltypes.Packet, data types.ConfirmUpgradePacketData) error {

	// TODO: packet timeout logic

	return nil
}
