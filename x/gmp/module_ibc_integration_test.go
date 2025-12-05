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

	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
)

// TestGMPInTransferStack verifies that GMP middleware properly wraps the transfer module
// and processes packets before forwarding them to the underlying module
func TestGMPInTransferStack(t *testing.T) {
	relayer := sdk.AccAddress(address.Module("relayer"))
	ctx := sdk.Context{}.WithLogger(log.NewNopLogger())

	// Create a mock transfer module that tracks if it was called
	mockTransfer := &mockIBCModule{
		lastCalled: "",
	}

	// Build the stack: transfer -> GMP
	transferStack := mockTransfer
	gmpStack := gmp.NewIBCModule(transferStack)

	t.Run("GMP processes TypeGeneralMessageWithToken and forwards to transfer", func(t *testing.T) {
		// Reset mock state
		mockTransfer.lastCalled = ""

		// Create a GMP message with TypeGeneralMessageWithToken
		payloadType, err := abi.NewType("string", "", nil)
		require.NoError(t, err)
		args := abi.Arguments{{Type: payloadType}}
		encoded, err := args.Pack("forward-memo")
		require.NoError(t, err)

		msg := gmp.Message{
			SourceChain:   "chainA",
			SourceAddress: "addr",
			Payload:       encoded,
			Type:          gmp.TypeGeneralMessageWithToken,
		}
		memo, err := json.Marshal(msg)
		require.NoError(t, err)

		// Create packet with GMP memo
		data := transfertypes.FungibleTokenPacketData{
			Denom:    "foo",
			Amount:   "1",
			Sender:   "sender",
			Receiver: "receiver",
			Memo:     string(memo),
		}
		bz, err := types.ModuleCdc.MarshalJSON(&data)
		require.NoError(t, err)
		packet := channeltypes.Packet{Data: bz}

		// Process packet through GMP stack
		ack := gmpStack.OnRecvPacket(ctx, types.Version, packet, relayer)

		// Verify GMP processed the packet and forwarded to transfer
		require.Equal(t, "OnRecvPacket", mockTransfer.lastCalled)
		require.NotNil(t, ack)

		// Verify the packet data was modified (memo should be replaced with unpacked payload)
		require.NotNil(t, mockTransfer.lastPacket)
		var modifiedData transfertypes.FungibleTokenPacketData
		err = types.ModuleCdc.UnmarshalJSON(mockTransfer.lastPacket.Data, &modifiedData)
		require.NoError(t, err)
		// The memo should be replaced with the unpacked ABI payload
		require.Equal(t, "forward-memo", modifiedData.Memo)
	})

	t.Run("GMP forwards non-GMP packets directly to transfer", func(t *testing.T) {
		// Reset mock state
		mockTransfer.lastCalled = ""

		// Create a regular transfer packet without GMP memo
		data := transfertypes.FungibleTokenPacketData{
			Denom:    "foo",
			Amount:   "1",
			Sender:   "sender",
			Receiver: "receiver",
			Memo:     "regular-memo",
		}
		bz, err := types.ModuleCdc.MarshalJSON(&data)
		require.NoError(t, err)
		packet := channeltypes.Packet{Data: bz}

		// Process packet through GMP stack
		ack := gmpStack.OnRecvPacket(ctx, types.Version, packet, relayer)

		// Verify packet was forwarded to transfer
		require.Equal(t, "OnRecvPacket", mockTransfer.lastCalled)
		require.NotNil(t, ack)

		// Verify packet data was not modified
		var forwardedData transfertypes.FungibleTokenPacketData
		err = types.ModuleCdc.UnmarshalJSON(mockTransfer.lastPacket.Data, &forwardedData)
		require.NoError(t, err)
		require.Equal(t, "regular-memo", forwardedData.Memo)
	})

	t.Run("GMP processes TypeGeneralMessage and forwards to transfer", func(t *testing.T) {
		// Reset mock state
		mockTransfer.lastCalled = ""

		msg := gmp.Message{
			SourceChain:   "chainA",
			SourceAddress: "addr",
			Payload:       []byte("payload"),
			Type:          gmp.TypeGeneralMessage,
		}
		memo, err := json.Marshal(msg)
		require.NoError(t, err)

		data := transfertypes.FungibleTokenPacketData{
			Denom:    "foo",
			Amount:   "1",
			Sender:   "sender",
			Receiver: "receiver",
			Memo:     string(memo),
		}
		bz, err := types.ModuleCdc.MarshalJSON(&data)
		require.NoError(t, err)
		packet := channeltypes.Packet{Data: bz}

		ack := gmpStack.OnRecvPacket(ctx, types.Version, packet, relayer)

		// Verify GMP processed and forwarded
		require.Equal(t, "OnRecvPacket", mockTransfer.lastCalled)
		require.NotNil(t, ack)
	})
}

