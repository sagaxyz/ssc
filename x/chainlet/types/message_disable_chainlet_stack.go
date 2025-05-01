package types

import (
	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDisableChainletStackVersion = "disable_chainlet_stack_version"

var _ sdk.Msg = &MsgDisableChainletStackVersion{}

func NewMsgDisableChainletStackVersion(creator string, displayName string, version string) *MsgDisableChainletStackVersion {
	return &MsgDisableChainletStackVersion{
		Creator:     creator,
		DisplayName: displayName,
		Version:     version,
	}
}

func (msg *MsgDisableChainletStackVersion) Route() string {
	return RouterKey
}

func (msg *MsgDisableChainletStackVersion) Type() string {
	return TypeMsgDisableChainletStackVersion
}

func (msg *MsgDisableChainletStackVersion) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDisableChainletStackVersion) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDisableChainletStackVersion) ValidateBasic() error {
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

	return nil
}
