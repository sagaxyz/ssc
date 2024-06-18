package types

import (
	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"k8s.io/utils/strings/slices"
)

const TypeMsgCreateChainletStack = "create_chainlet_stack"

var _ sdk.Msg = &MsgCreateChainletStack{}

func NewMsgCreateChainletStack(creator string, displayName string, description string, image string, version string, checksum string, fees ChainletStackFees) *MsgCreateChainletStack {
	return &MsgCreateChainletStack{
		Creator:     creator,
		DisplayName: displayName,
		Description: description,
		Version:     version,
		Image:       image,
		Checksum:    checksum,
		Fees:        fees,
	}
}

func (msg *MsgCreateChainletStack) Route() string {
	return RouterKey
}

func (msg *MsgCreateChainletStack) Type() string {
	return TypeMsgCreateChainletStack
}

func (msg *MsgCreateChainletStack) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateChainletStack) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateChainletStack) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.DisplayName == "" {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "display name cannot be empty")
	}

	if len(msg.Description) > 120 {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "description cannot be longer than 120 characters")
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

	coin, err := sdk.ParseCoinNormalized(msg.Fees.EpochFee)
	if err != nil {
		return ErrInvalidCoin
	}

	if !coin.Amount.GT(sdk.ZeroInt()) {
		return ErrInvalidCoin
	}

	coin, err = sdk.ParseCoinNormalized(msg.Fees.SetupFee)
	if err != nil {
		return ErrInvalidCoin
	}

	if !coin.Amount.GT(sdk.ZeroInt()) {
		return ErrInvalidCoin
	}

	if !slices.Contains([]string{"day", "week", "hour", "minute"}, msg.Fees.EpochLength) {
		return ErrInvalidEpoch
	}
	return nil
}
