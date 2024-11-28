package keeper

import (
	"context"
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sagaxyz/ssc/x/gmp/types"
)

// AxelarGMPAcc is the address that receives the message from a cosmos chain
const AxelarGMPAcc = "axelar1dv4u5k73pzqrxlzujxg3qp8kvc3pje7jtdvu72npnt5zhq05ejcsn5qme5"

type MessageType int

const (
	// TypeUnrecognized means coin type is unrecognized
	TypeUnrecognized = iota
	// TypeGeneralMessage is a pure message
	TypeGeneralMessage
	// TypeGeneralMessageWithToken is a general message with token
	TypeGeneralMessageWithToken
	// TypeSendToken is a direct token transfer
	TypeSendToken
)

// Message is attached in ICS20 packet memo field
type Message struct {
	DestinationChain   string `json:"destination_chain"`
	DestinationAddress string `json:"destination_address"`
	Payload            []byte `json:"payload"`
	Type               int64  `json:"type"`
}

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

func (k msgServer) Transfer(goCtx context.Context, msg *types.MsgMultiSend) (*types.MsgMultiSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// sender, err := sdk.AccAddressFromBech32(msg.Sender)
	// if err != nil {
	// 	return nil, err
	// }

	// build payload that can be decoded by solidity
	addressesType, err := abi.NewType("address[]", "address[]", nil)
	if err != nil {
		return nil, err
	}

	var addresses []common.Address
	addresses = append(addresses, common.HexToAddress(msg.ToAddress))
	// for _, receiver := range msg.ToAddress {
	// 	addresses = append(addresses, common.HexToAddress(receiver))
	// }

	payload, err := abi.Arguments{{Type: addressesType}}.Pack(addresses)
	if err != nil {
		return nil, err
	}

	message := Message{
		DestinationChain:   msg.DestinationChain,
		DestinationAddress: msg.ToAddress,
		Payload:            payload,
		Type:               TypeGeneralMessageWithToken,
	}

	bz, err := json.Marshal(&message)
	if err != nil {
		return nil, err
	}

	ibcMessage := ibctransfertypes.NewMsgTransfer(
		ibctransfertypes.PortID,
		"channel-1", // hard-coded channel id for demo
		*msg.Amount,
		msg.FromAddress,
		AxelarGMPAcc, // TODO: use config
		clienttypes.ZeroHeight(),
		uint64(ctx.BlockTime().Add(6*time.Hour).UnixNano()),
		string(bz),
	)

	res, err := k.keeper.Transfer(goCtx, ibcMessage)
	if err != nil {
		return nil, err
	}

	return res, nil
}

var _ types.MsgServer = msgServer{}
