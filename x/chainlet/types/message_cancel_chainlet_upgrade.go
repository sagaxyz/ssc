package types

import (
	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCancelChainletUpgrade = "cancel_chainlet_upgrade"

var _ sdk.Msg = &MsgCancelChainletUpgrade{}

func NewMsgCancelChainletUpgrade(creator string, chainId string, stackVersion string, channelID string) *MsgCancelChainletUpgrade {
	return &MsgCancelChainletUpgrade{
		Creator:      creator,
		ChainId:      chainId,
		Version: stackVersion,
		ChannelId:    channelID,
	}
}

func (msg *MsgCancelChainletUpgrade) Route() string {
	return RouterKey
}

func (msg *MsgCancelChainletUpgrade) Type() string {
	return TypeMsgCancelChainletUpgrade
}

func (msg *MsgCancelChainletUpgrade) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
