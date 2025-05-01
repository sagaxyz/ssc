package types

import (
	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateChainletStack = "update_chainlet_stack"

var _ sdk.Msg = &MsgUpdateChainletStack{}

func NewMsgUpdateChainletStack(creator string, displayName string, image string, version string, checksum string) *MsgUpdateChainletStack {
	return &MsgUpdateChainletStack{
		Creator:     creator,
		DisplayName: displayName,
		Version:     version,
		Image:       image,
		Checksum:    checksum,
	}
}

func (msg *MsgUpdateChainletStack) Route() string {
	return RouterKey
}

func (msg *MsgUpdateChainletStack) Type() string {
	return TypeMsgUpdateChainletStack
}

func (msg *MsgUpdateChainletStack) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateChainletStack) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateChainletStack) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.DisplayName == "" {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "display name cannot be empty")
	}

	if msg.Version == "" {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "version cannot be empty")
	}

	if msg.Image == "" {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "image cannot be empty")
	}

	if msg.Checksum == "" {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "checksum cannot be empty")
	}
	return nil
}
