package types

import (
	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgWithdraw = "withdraw"

var _ sdk.Msg = &MsgWithdraw{}

func NewMsgWithdraw(creator string, chainId string) *MsgWithdraw {
	return &MsgWithdraw{
		Creator: creator,
		ChainId: chainId,
	}
}

func (msg *MsgWithdraw) Route() string {
	return RouterKey
}

func (msg *MsgWithdraw) Type() string {
	return TypeMsgWithdraw
}

func (msg *MsgWithdraw) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
