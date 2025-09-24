package types

import (
	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateChainletStack = "create_chainlet_stack"

var _ sdk.Msg = &MsgCreateChainletStack{}

func NewMsgCreateChainletStack(creator string, displayName string, description string, image string, version string, checksum string, fees ChainletStackFees, ccvConsumer bool) *MsgCreateChainletStack {
	return &MsgCreateChainletStack{
		Creator:     creator,
		DisplayName: displayName,
		Description: description,
		Version:     version,
		Image:       image,
		Checksum:    checksum,
		Fees:        fees,
		CcvConsumer: ccvConsumer,
	}
}

func (msg *MsgCreateChainletStack) Route() string {
	return RouterKey
}

func (msg *MsgCreateChainletStack) Type() string {
	return TypeMsgCreateChainletStack
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
	if !coin.Amount.IsPositive() {
		return ErrInvalidCoin
	}
	if msg.Fees.Denom != coin.Denom {
		return ErrInvalidDenom
	}

	coin, err = sdk.ParseCoinNormalized(msg.Fees.SetupFee)
	if err != nil {
		return ErrInvalidCoin
	}
	if !coin.Amount.IsPositive() {
		return ErrInvalidCoin
	}
	if msg.Fees.Denom != coin.Denom {
		return ErrInvalidDenom
	}

	return nil
}
