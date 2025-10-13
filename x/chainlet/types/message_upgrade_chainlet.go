package types

import (
	"time"

	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpgradeChainlet = "upgrade_chainlet"

var _ sdk.Msg = &MsgUpgradeChainlet{}

func NewMsgUpgradeChainlet(creator string, chainId string, stackVersion string, heightDelta uint64, channelID string, unbondingPeriod *time.Duration) *MsgUpgradeChainlet {
	return &MsgUpgradeChainlet{
		Creator:         creator,
		ChainId:         chainId,
		StackVersion:    stackVersion,
		HeightDelta:     heightDelta,
		ChannelId:       channelID,
		UnbondingPeriod: unbondingPeriod,
	}
}

func (msg *MsgUpgradeChainlet) Route() string {
	return RouterKey
}

func (msg *MsgUpgradeChainlet) Type() string {
	return TypeMsgUpgradeChainlet
}

func (msg *MsgUpgradeChainlet) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
