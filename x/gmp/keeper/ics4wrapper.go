package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v10/modules/core/exported"
	"github.com/sagaxyz/ssc/x/gmp/types"
)

var _ porttypes.ICS4Wrapper = &Keeper{}

// GetAppVersion implements types.ICS4Wrapper.
func (k Keeper) GetAppVersion(ctx sdk.Context, portID string, channelID string) (string, bool) {
	return k.ics4Wrapper.GetAppVersion(ctx, portID, channelID)
}

// SendPacket implements types.ICS4Wrapper.
func (k Keeper) SendPacket(ctx sdk.Context, sourcePort string, sourceChannel string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, data []byte) (sequence uint64, err error) {
	return k.ics4Wrapper.SendPacket(ctx, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}

// WriteAcknowledgement implements types.ICS4Wrapper.
func (k Keeper) WriteAcknowledgement(ctx sdk.Context, packet exported.PacketI, ack exported.Acknowledgement) error {
	// check if we have a packet for this sequence
	storedPacket, found := k.GetPacket(ctx, packet.GetSourceChannel(), packet.GetSourcePort(), packet.GetSequence())
	if found {
		packet = storedPacket
		// if not found, we skip this and use the packet coming in the call

		// delete the packet from the store, as we won't use it anymore
		k.DeletePacket(ctx, packet.GetSourceChannel(), packet.GetSourcePort(), packet.GetSequence())

		ctx.Logger().Debug("gmp: found packet in store, using it", "channel", packet.GetSourceChannel(), "port", packet.GetSourcePort(), "sequence", packet.GetSequence())
	}
	return k.ics4Wrapper.WriteAcknowledgement(ctx, packet, ack)
}

func (k Keeper) SetPacket(ctx sdk.Context, packet exported.PacketI) {
	store := ctx.KVStore(k.storeKey)

	newPacket := channeltypes.Packet{
		Data:               packet.GetData(),
		Sequence:           packet.GetSequence(),
		SourcePort:         packet.GetSourcePort(),
		SourceChannel:      packet.GetSourceChannel(),
		DestinationPort:    packet.GetDestPort(),
		DestinationChannel: packet.GetDestChannel(),
		TimeoutHeight:      clienttypes.NewHeight(packet.GetTimeoutHeight().GetRevisionNumber(), packet.GetTimeoutHeight().GetRevisionHeight()),
		TimeoutTimestamp:   packet.GetTimeoutTimestamp(),
	}
	bz, err := newPacket.Marshal()
	if err != nil {
		panic(err)
	}

	store.Set(types.PacketKey(packet.GetSourceChannel(), packet.GetSourcePort(), packet.GetSequence()), bz)
}

func (k Keeper) GetPacket(ctx sdk.Context, channelID, port string, sequence uint64) (channeltypes.Packet, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PacketKey(channelID, port, sequence))
	if len(bz) == 0 {
		return channeltypes.Packet{}, false
	}
	var packet channeltypes.Packet
	err := packet.Unmarshal(bz)
	if err != nil {
		return channeltypes.Packet{}, false
	}
	return packet, true
}

func (k Keeper) DeletePacket(ctx sdk.Context, channelID, port string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PacketKey(channelID, port, sequence))
}
