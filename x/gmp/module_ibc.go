package gmp

import (
	"encoding/json"
	"fmt"
	"strings"

	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/sagaxyz/ssc/x/gmp/keeper"
	"github.com/sagaxyz/ssc/x/gmp/types"
)

type GeneralMessageHandler interface {
	HandleGeneralMessage(ctx sdk.Context, srcChain, srcAddress string, payload []byte) error
	HandleGeneralMessageWithToken(ctx sdk.Context, srcChain, srcAddress string, payload []byte, receiver string, coin sdk.Coin) error
}

// Message is attached in ICS20 packet memo field
type Message struct {
	SourceChain   string `json:"source_chain"`
	SourceAddress string `json:"source_address"`
	Payload       []byte `json:"payload"`
	Type          int64  `json:"type"`
}

type MessageType int

const (
	// TypeUnrecognized means coin type is unrecognized
	TypeUnrecognized = iota
	TypeGeneralMessage
	// TypeGeneralMessageWithToken is a general message with token
	TypeGeneralMessageWithToken
)

type PFMPayload struct {
	Receiver string      `json:"receiver"`
	Channel  string      `json:"channel"`
	Next     *PFMPayload `json:"next"`
}

type ForwardPayload struct {
	Forward *Forward `json:"forward"`
}

type Forward struct {
	Receiver string          `json:"receiver"`
	Port     string          `json:"port"`
	Channel  string          `json:"channel"`
	Next     *ForwardPayload `json:"next,omitempty"`
}

type IBCModule struct {
	keeper keeper.Keeper
	app    porttypes.IBCModule
}

func NewIBCModule(k keeper.Keeper, app porttypes.IBCModule) IBCModule {
	return IBCModule{
		keeper: k,
		app:    app,
	}
}

// OnChanOpenInit implements the IBCModule interface
func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {

	// Require portID is the portID module is bound to
	// boundPort := im.keeper.GetPort(ctx)
	// if boundPort != portID {
	// 	return "", cosmossdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, boundPort)
	// }

	// if version != types.Version {
	// 	return "", cosmossdkerrors.Wrapf(types.ErrInvalidVersion, "got %s, expected %s", version, types.Version)
	// }

	// // Claim channel capability passed back by IBC module
	// if err := im.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
	// 	return "", err
	// }

	return im.app.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, chanCap, counterparty, version)
}

// OnChanOpenTry implements the IBCModule interface
func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {

	// Require portID is the portID module is bound to
	// boundPort := im.keeper.GetPort(ctx)
	// if boundPort != portID {
	// 	return "", cosmossdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, boundPort)
	// }

	// if counterpartyVersion != types.Version {
	// 	return "", cosmossdkerrors.Wrapf(types.ErrInvalidVersion, "invalid counterparty version: got: %s, expected %s", counterpartyVersion, types.Version)
	// }

	// // Module may have already claimed capability in OnChanOpenInit in the case of crossing hellos
	// // (ie chainA and chainB both call ChanOpenInit before one of them calls ChanOpenTry)
	// // If module can already authenticate the capability then module already owns it so we don't need to claim
	// // Otherwise, module does not have channel capability and we must claim it from IBC
	// if !im.keeper.AuthenticateCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)) {
	// 	// Only claim channel capability passed back by IBC module if we do not already own it
	// 	if err := im.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
	// 		return "", err
	// 	}
	// }

	return im.app.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, chanCap, counterparty, counterpartyVersion)
}

// OnChanOpenAck implements the IBCModule interface
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	// if counterpartyVersion != types.Version {
	// 	return cosmossdkerrors.Wrapf(types.ErrInvalidVersion, "invalid counterparty version: %s, expected %s", counterpartyVersion, types.Version)
	// }
	return im.app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

// OnChanOpenConfirm implements the IBCModule interface
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return im.app.OnChanOpenConfirm(ctx, portID, channelID)
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return im.app.OnChanCloseInit(ctx, portID, channelID)
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return im.app.OnChanCloseConfirm(ctx, portID, channelID)
}

// OnRecvPacket implements the IBCModule interface
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	// this line is used by starport scaffolding # oracle/packet/module/recv

	var data transfertypes.FungibleTokenPacketData
	var err error
	if err := types.ModuleCdc.UnmarshalJSON(modulePacket.GetData(), &data); err != nil {
		ctx.Logger().Debug(fmt.Sprintf("cannot unmarshal ICS-20 transfer packet data: %s", err.Error()))
		return im.app.OnRecvPacket(ctx, modulePacket, relayer)
	}

	var msg Message
	if err = json.Unmarshal([]byte(data.GetMemo()), &msg); err != nil {
		ctx.Logger().Debug(fmt.Sprintf("cannot unmarshal memo: %s", err.Error()))
		return im.app.OnRecvPacket(ctx, modulePacket, relayer)
	}

	if msg.Payload == nil {
		return im.app.OnRecvPacket(ctx, modulePacket, relayer)
	}

	// // authenticate the message with packet sender + channel-id
	// // TODO: authenticate the message with channel-id
	// if data.Sender != AxelarGMPAcc {
	// 	return ack
	// }

	switch msg.Type {
	case TypeGeneralMessage:
		ctx.Logger().Info(fmt.Sprintf("Got pure general message: %v", msg))
		return nil //?
	case TypeGeneralMessageWithToken:
		ctx.Logger().Info(fmt.Sprintf("Got general message with token: %v", msg))
		payloadType, err := abi.NewType("string", "", nil)
		if err != nil {
			ctx.Logger().Info(fmt.Sprintf("failed to create reflection: %s", err.Error()))
			return channeltypes.NewErrorAcknowledgement(cosmossdkerrors.Wrapf(transfertypes.ErrInvalidMemo, "unable to define new abi type (%s)", err.Error()))
		}

		// Add 4 bytes to the payload to match the length of the payload
		payloadData := msg.Payload
		payloadData = append(make([]byte, 4), payloadData...)

		args, err := abi.Arguments{{Type: payloadType}}.Unpack(payloadData)
		if err != nil {
			ctx.Logger().Info(fmt.Sprintf("failed to unpack: %s", err.Error()))
			return channeltypes.NewErrorAcknowledgement(cosmossdkerrors.Wrapf(transfertypes.ErrInvalidMemo, "unable to unpack payload (%s)", err.Error()))
		}
		pfmPayload := args[0].(string)
		ctx.Logger().Info(fmt.Sprintf("Got pfmPayload: %v", pfmPayload))
		// pfmPayload is like saga1azf8fv5x9l5n9lh8h5s4l9m9ju76xhd9fhjjqk,channel-1
		forwardAddress, channelID := strings.Split(pfmPayload, ",")[0], strings.Split(pfmPayload, ",")[1]
		updatedPfmPayload := &PFMPayload{forwardAddress, channelID, nil}
		ctx.Logger().Info(fmt.Sprintf("Updated pfmPayload: %+v", updatedPfmPayload))
		// Now update modulePacket with new memo
		// Convert payload to the new structure
		forwardPayload := convertToForwardPayload(updatedPfmPayload)
		updatedMemo, err := json.Marshal(forwardPayload)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(cosmossdkerrors.Wrapf(transfertypes.ErrInvalidMemo, "memo convertion error: %s", err.Error()))
		}
		data.Memo = string(updatedMemo)
		modulePacket.Data, err = types.ModuleCdc.MarshalJSON(&data)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(cosmossdkerrors.Wrapf(transfertypes.ErrInvalidMemo, "cannot marshal updated data: %s", err.Error()))
		}
		return im.app.OnRecvPacket(ctx, modulePacket, relayer)

	default:
		return channeltypes.NewErrorAcknowledgement(cosmossdkerrors.Wrapf(transfertypes.ErrInvalidMemo, "unrecognized message type (%d)", msg.Type))
	}
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	return im.app.OnAcknowledgementPacket(ctx, modulePacket, acknowledgement, relayer)
}

// OnTimeoutPacket implements the IBCModule interface
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	return im.app.OnTimeoutPacket(ctx, modulePacket, relayer)
}

// Recursive function to convert PFMPayload to ForwardPayload
func convertToForwardPayload(pfm *PFMPayload) *ForwardPayload {
	if pfm == nil {
		return nil
	}
	return &ForwardPayload{
		Forward: &Forward{
			Receiver: pfm.Receiver,
			Port:     "transfer",
			Channel:  pfm.Channel,
			Next:     convertToForwardPayload(pfm.Next),
		},
	}
}

func parseDenom(packet channeltypes.Packet, denom string) string {
	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), denom) {
		// remove prefix added by sender chain
		voucherPrefix := transfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())
		unprefixedDenom := denom[len(voucherPrefix):]

		// coin denomination used in sending from the escrow address
		denom = unprefixedDenom

		// The denomination used to send the coins is either the native denom or the hash of the path
		// if the denomination is not native.
		denomTrace := transfertypes.ParseDenomTrace(unprefixedDenom)
		if denomTrace.Path != "" {
			denom = denomTrace.IBCDenom()
		}

		return denom
	}

	prefixedDenom := transfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel()) + denom
	denom = transfertypes.ParseDenomTrace(prefixedDenom).IBCDenom()

	return denom
}
