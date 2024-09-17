package types

import (
	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetPeers = "set_peers"

var _ sdk.Msg = &MsgSetPeers{}

func NewMsgSetPeers(validator string, chainId string, peers ...string) *MsgSetPeers {
	return &MsgSetPeers{
		Validator: validator,
		ChainId:   chainId,
		Peers:     peers,
	}
}

func (msg *MsgSetPeers) Route() string {
	return RouterKey
}

func (msg *MsgSetPeers) Type() string {
	return TypeMsgSetPeers
}

func (msg *MsgSetPeers) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}

	if msg.ChainId == "" {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "value cannot be empty")
	}

	//TODO check peers

	return nil
}
