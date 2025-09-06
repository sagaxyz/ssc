package gmp

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/sagaxyz/ssc/x/gmp/types"
)

// Message is attached in ICS20 packet memo field
type Message struct {
	SourceChain   string `json:"source_chain"`
	SourceAddress string `json:"source_address"`
	Payload       []byte `json:"payload"`
	Type          int64  `json:"type"`
}

const (
	// TypeUnrecognized means coin type is unrecognized
	TypeUnrecognized = iota
	TypeGeneralMessage
	// TypeGeneralMessageWithToken is a general message with token
	TypeGeneralMessageWithToken
)

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
	app porttypes.IBCModule
}

func NewIBCModule(app porttypes.IBCModule) IBCModule {
	return IBCModule{
		app: app,
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
	if err := types.ModuleCdc.UnmarshalJSON(modulePacket.GetData(), &data); err != nil {
		ctx.Logger().Debug(fmt.Sprintf("cannot unmarshal ICS-20 transfer packet data: %s", err.Error()))
		return im.app.OnRecvPacket(ctx, modulePacket, relayer)
	}

	var msg Message
	if err := json.Unmarshal([]byte(data.GetMemo()), &msg); err != nil {
		ctx.Logger().Debug(fmt.Sprintf("cannot unmarshal memo: %s", err.Error()))
		return im.app.OnRecvPacket(ctx, modulePacket, relayer)
	}

	if msg.Payload == nil {
		return im.app.OnRecvPacket(ctx, modulePacket, relayer)
	}

	switch msg.Type {
	case TypeGeneralMessage:
		ctx.Logger().Debug(fmt.Sprintf("Got pure general message: %v", msg))
		return im.app.OnRecvPacket(ctx, modulePacket, relayer)
	case TypeGeneralMessageWithToken:
		ctx.Logger().Debug(fmt.Sprintf("Got general message with token: %v", msg))
		payloadType, err := abi.NewType("string", "", nil)
		if err != nil {
			ctx.Logger().Info(fmt.Sprintf("failed to create reflection: %s", err.Error()))
			return channeltypes.NewErrorAcknowledgement(cosmossdkerrors.Wrapf(transfertypes.ErrInvalidMemo, "unable to define new abi type (%s)", err.Error()))
		}
		payloadData := msg.Payload
		ctx.Logger().Debug(fmt.Sprintf("ABI-encoded data: %v", hex.EncodeToString(payloadData)))

		args, err := abi.Arguments{{Type: payloadType}}.Unpack(payloadData)
		if err != nil {
			ctx.Logger().Debug(fmt.Sprintf("failed to unpack: %s", err.Error()))
			return channeltypes.NewErrorAcknowledgement(cosmossdkerrors.Wrapf(transfertypes.ErrInvalidMemo, "unable to unpack payload (%s)", err.Error()))
		}
		pfmPayload := args[0].(string)
		data.Memo = string(pfmPayload)
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
