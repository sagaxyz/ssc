package keeper

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v10/modules/core/24-host"
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

// OnAcknowledgementCreateUpgradePacket responds to the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementCreateUpgradePacket(ctx sdk.Context, packet channeltypes.Packet, data chainlettypes.CreateUpgradePacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		fmt.Printf("XXX packet ack error: data %+v ack %+v: %s\n", data, ack, dispatchedAck.Error)
		//TODO cancel upgrade
		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck chainlettypes.CreateUpgradePacketAck

		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		fmt.Printf("XXX packet ack OK: data %+v ack %+v: %s\n", data, ack, dispatchedAck)
		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}

// OnTimeoutCreateUpgradePacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutCreateUpgradePacket(ctx sdk.Context, packet channeltypes.Packet, data chainlettypes.CreateUpgradePacketData) error {
	fmt.Printf("XXX packet timeout: data %+v\n", data)
	return nil
}

// OnRecvConfirmUpgradePacket processes packet reception
func (k Keeper) OnRecvConfirmUpgradePacket(ctx sdk.Context, packet channeltypes.Packet, data chainlettypes.ConfirmUpgradePacketData) (packetAck chainlettypes.ConfirmUpgradePacketAck, err error) {
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
	//if data.Height != chainlet.Upgrade.Height - 1 {
	//	err = fmt.Errorf("unexpected upgrade height: %d != %d", data.Height, chainlet.Upgrade.Height -1)
	//	fmt.Printf("XXX upgrading %s: %s\n", data.ChainId, err)
	//	return
	//}
	err = k.finishUpgrading(ctx, &chainlet)
	if err != nil {
		fmt.Printf("XXX upgrading %s: %s\n", data.ChainId, err)
		return
	}
	fmt.Printf("XXX upgrading %s: DONE\n")

	return packetAck, nil
}

