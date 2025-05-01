package types

import (
	cosmossdkerrors "cosmossdk.io/errors"
	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sagaxyz/ssc/x/chainlet/tags"
)

const TypeMsgLaunchChainlet = "launch_chainlet"

const ChainletGasLimit uint64 = 30000000

var _ sdk.Msg = &MsgLaunchChainlet{}

func NewMsgLaunchChainlet(
	creator string,
	maintainers []string,
	chainletStackName string,
	chainletStackVersion string,
	chainletName string,
	chainId string,
	denom string,
	params ChainletParams,
	tags []string,
	serviceChainlet bool) *MsgLaunchChainlet {
	return &MsgLaunchChainlet{
		Creator:              creator,
		Maintainers:          maintainers,
		ChainletStackName:    chainletStackName,
		ChainletStackVersion: chainletStackVersion,
		ChainletName:         chainletName,
		ChainId:              chainId,
		Denom:                denom,
		Params:               params,
		Tags:                 tags,
		IsServiceChainlet:    serviceChainlet,
	}
}

func (msg *MsgLaunchChainlet) Route() string {
	return RouterKey
}

func (msg *MsgLaunchChainlet) Type() string {
	return TypeMsgLaunchChainlet
}

func (msg *MsgLaunchChainlet) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgLaunchChainlet) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgLaunchChainlet) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.ChainletName == "" {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "chainlet name cannot be empty")
	}

	for _, maintainer := range msg.Maintainers {
		_, err := sdk.AccAddressFromBech32(maintainer)
		if err != nil {
			return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid maintainer address (%s)", err)
		}
	}

	if msg.ChainletStackName == "" {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "chainlet stack name cannot be empty")
	}

	if msg.ChainletStackVersion == "" {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "chainlet stack version cannot be empty")
	}

	valid := validateChainId(msg.ChainId)
	if !valid {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "chain id %s is invalid", msg.ChainId)
	}

	if msg.Denom == "" {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "denom cannot be empty")
	}

	valid = validateDenom(msg.Denom)
	if !valid {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "denom %s is invalid", msg.Denom)
	}

	for _, genacct := range msg.Params.GenAcctBalances.List {
		_, err := sdk.AccAddressFromBech32(genacct.Address)
		if err != nil {
			return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "could not parse address from supplied string %s. Error: %s", genacct.Address, err)
		}

		_, err = math.ParseUint(genacct.Balance)
		if err != nil {
			return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "could not parse uint from supplied string %s. Error: %s", genacct.Balance, err)
		}
	}
	if msg.Params.FixedBaseFee != "" {
		i, ok := math.NewIntFromString(msg.Params.FixedBaseFee)
		if !ok {
			return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "could not parse fixed base fee '%s'", msg.Params.FixedBaseFee)
		}
		if i.IsNegative() {
			return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "fixed base fee cannot be negative")
		}
	}
	if msg.Params.FeeAccount != "" {
		_, err := sdk.AccAddressFromBech32(msg.Params.FeeAccount)
		if err != nil {
			return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid fee account address: %s", err)
		}
	}
	if msg.Params.GasLimit > ChainletGasLimit {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "gas limit too high, must be less than %v", ChainletGasLimit)
	}
	if _, ok := tags.VerifyAndTruncate(msg.Tags); !ok {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid tags")
	}
	return nil
}
