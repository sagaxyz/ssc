package gmp_test

import (
	"encoding/json"
	"testing"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/sagaxyz/ssc/x/gmp"
	"github.com/sagaxyz/ssc/x/gmp/types"
	"github.com/stretchr/testify/require"

	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

type mockIBCModule struct {
	lastPacket  channeltypes.Packet
	lastRelayer sdk.AccAddress
	lastCalled  string
	returnAck   ibcexported.Acknowledgement
}

func (m *mockIBCModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	m.lastPacket = packet
	m.lastRelayer = relayer
	m.lastCalled = "OnRecvPacket"
	if m.returnAck != nil {
		return m.returnAck
	}
	return channeltypes.NewResultAcknowledgement([]byte("mock"))
}

func (m *mockIBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {
	return "", nil
}

func (m *mockIBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	return "", nil
}

func (m *mockIBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID,
	counterpartyChannelID,
	counterpartyVersion string,
) error {
	return nil
}

func (m *mockIBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

func (m *mockIBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

func (m *mockIBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

func (m *mockIBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	return nil
}

func (m *mockIBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	return nil
}

func makePacketWithMemo(t *testing.T, memo string) channeltypes.Packet {
	data := transfertypes.FungibleTokenPacketData{
		Denom:    "foo",
		Amount:   "1",
		Sender:   "sender",
		Receiver: "receiver",
		Memo:     memo,
	}
	bz, err := types.ModuleCdc.MarshalJSON(&data)
	require.NoError(t, err)
	return channeltypes.Packet{
		Data: bz,
	}
}

func TestOnRecvPacket(t *testing.T) {
	relayer := sdk.AccAddress(address.Module("relayer"))
	ctx := sdk.Context{}.WithLogger(log.NewNopLogger())

	t.Run("invalid packet data", func(t *testing.T) {
		mock := &mockIBCModule{}
		mod := gmp.NewIBCModule(mock)
		// invalid JSON
		packet := channeltypes.Packet{Data: []byte("{invalid-json")}
		ack := mod.OnRecvPacket(ctx, packet, relayer)
		require.Equal(t, "OnRecvPacket", mock.lastCalled)
		require.NotNil(t, ack)
	})

	t.Run("invalid memo", func(t *testing.T) {
		mock := &mockIBCModule{}
		mod := gmp.NewIBCModule(mock)
		// valid packet, invalid memo
		data := transfertypes.FungibleTokenPacketData{
			Denom:    "foo",
			Amount:   "1",
			Sender:   "sender",
			Receiver: "receiver",
			Memo:     "{invalid-json",
		}
		bz, err := types.ModuleCdc.MarshalJSON(&data)
		require.NoError(t, err)
		packet := channeltypes.Packet{Data: bz}
		ack := mod.OnRecvPacket(ctx, packet, relayer)
		require.Equal(t, "OnRecvPacket", mock.lastCalled)
		require.NotNil(t, ack)
	})

	t.Run("empty payload", func(t *testing.T) {
		mock := &mockIBCModule{}
		mod := gmp.NewIBCModule(mock)
		msg := gmp.Message{
			SourceChain:   "chainA",
			SourceAddress: "addr",
			Payload:       []byte{},
			Type:          gmp.TypeGeneralMessage,
		}
		memo, err := json.Marshal(msg)
		require.NoError(t, err)
		packet := makePacketWithMemo(t, string(memo))
		ack := mod.OnRecvPacket(ctx, packet, relayer)
		require.Equal(t, "OnRecvPacket", mock.lastCalled)
		require.NotNil(t, ack)
	})

	t.Run("TypeGeneralMessage", func(t *testing.T) {
		mock := &mockIBCModule{}
		mod := gmp.NewIBCModule(mock)
		msg := gmp.Message{
			SourceChain:   "chainA",
			SourceAddress: "addr",
			Payload:       []byte("payload"),
			Type:          gmp.TypeGeneralMessage,
		}
		memo, err := json.Marshal(msg)
		require.NoError(t, err)
		packet := makePacketWithMemo(t, string(memo))
		ack := mod.OnRecvPacket(ctx, packet, relayer)
		require.Equal(t, "OnRecvPacket", mock.lastCalled)
		require.NotNil(t, ack)
	})

	t.Run("TypeGeneralMessageWithToken valid ABI", func(t *testing.T) {
		mock := &mockIBCModule{}
		mod := gmp.NewIBCModule(mock)
		payloadType, err := abi.NewType("string", "", nil)
		require.NoError(t, err)
		args := abi.Arguments{{Type: payloadType}}
		encoded, err := args.Pack("pfm-memo")
		require.NoError(t, err)
		msg := gmp.Message{
			SourceChain:   "chainA",
			SourceAddress: "addr",
			Payload:       encoded,
			Type:          gmp.TypeGeneralMessageWithToken,
		}
		memo, err := json.Marshal(msg)
		require.NoError(t, err)
		packet := makePacketWithMemo(t, string(memo))
		ack := mod.OnRecvPacket(ctx, packet, relayer)
		require.Equal(t, "OnRecvPacket", mock.lastCalled)
		require.NotNil(t, ack)
	})

	t.Run("TypeGeneralMessageWithToken invalid ABI", func(t *testing.T) {
		mock := &mockIBCModule{}
		mod := gmp.NewIBCModule(mock)
		msg := gmp.Message{
			SourceChain:   "chainA",
			SourceAddress: "addr",
			Payload:       []byte("{}"), // not valid ABI encoding
			Type:          gmp.TypeGeneralMessageWithToken,
		}
		memo, err := json.Marshal(msg)
		require.NoError(t, err)
		packet := makePacketWithMemo(t, string(memo))
		ack := mod.OnRecvPacket(ctx, packet, relayer)
		require.Equal(t, "OnRecvPacket", mock.lastCalled)
		require.NotNil(t, ack)
	})

	t.Run("unrecognized type", func(t *testing.T) {
		mock := &mockIBCModule{}
		mod := gmp.NewIBCModule(mock)
		msg := gmp.Message{
			SourceChain:   "chainA",
			SourceAddress: "addr",
			Payload:       []byte("payload"),
			Type:          int64(999),
		}
		memo, err := json.Marshal(msg)
		require.NoError(t, err)
		packet := makePacketWithMemo(t, string(memo))
		ack := mod.OnRecvPacket(ctx, packet, relayer)
		require.NotNil(t, ack)
		// Should not call underlying OnRecvPacket for unrecognized type
		require.NotEqual(t, "OnRecvPacket", mock.lastCalled)
	})
}
