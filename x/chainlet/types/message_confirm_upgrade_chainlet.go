package types

import (
	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgConfirmUpgradeChainlet = "upgrade_chainlet"

var _ sdk.Msg = &MsgConfirmUpgradeChainlet{}

func NewMsgConfirmUpgradeChainlet(creator string, chainId string, stackVersion string, heightDelta uint64) *MsgConfirmUpgradeChainlet {
	return &MsgConfirmUpgradeChainlet{
		Creator:      creator,
		ChainId:      chainId,
		StackVersion: stackVersion,
		HeightDelta:  heightDelta,
	}
}

func (msg *MsgConfirmUpgradeChainlet) Route() string {
	return RouterKey
}

func (msg *MsgConfirmUpgradeChainlet) Type() string {
	return TypeMsgConfirmUpgradeChainlet
}

func (msg *MsgConfirmUpgradeChainlet) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
