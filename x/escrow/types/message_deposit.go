package types

import (
	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDeposit = "deposit"

var _ sdk.Msg = &MsgDeposit{}

func NewMsgDeposit(creator string, amount string, chainId string) *MsgDeposit {
	return &MsgDeposit{
		Creator: creator,
		Amount:  amount,
		ChainId: chainId,
	}
}

func (msg *MsgDeposit) Route() string {
	return RouterKey
}

func (msg *MsgDeposit) Type() string {
	return TypeMsgDeposit
}

func (msg *MsgDeposit) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	coin, err := sdk.ParseCoinNormalized(msg.Amount)
	if err != nil {
		return cosmossdkerrors.Wrapf(ErrInvalidCoin, "invalid coin (%s)", msg.Amount)
	}

	if !coin.Amount.GT(sdk.ZeroInt()) {
		return cosmossdkerrors.Wrapf(ErrInvalidCoin, "must send more than 0 coins (%s)", coin.Amount.String())
	}
	return nil
}
