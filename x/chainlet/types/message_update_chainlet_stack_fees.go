package types

import (
	"strings"

	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (m *MsgUpdateChainletStackFees) ValidateBasic() error {
	if strings.TrimSpace(m.Creator) == "" {
		return cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", m.Creator)
	}
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return err
	}
	if strings.TrimSpace(m.ChainletStackName) == "" {
		return ErrInvalidChainletStack
	}
	if len(m.Fees) == 0 {
		return ErrInvalidFees.Wrap("fees cannot be empty")
	}
	seen := make(map[string]struct{})
	for i, f := range m.Fees {
		coin, err := sdk.ParseCoinNormalized(strings.TrimSpace(f.EpochFee))
		if err != nil {
			return ErrInvalidFees.Wrapf("fee[%d]=%q: %v", i, f.EpochFee, err)
		}
		if !coin.IsPositive() {
			return ErrInvalidFees.Wrapf("fee[%d]=%q must be positive", i, f.EpochFee)
		}
		if _, ok := seen[coin.Denom]; ok {
			return ErrDuplicateDenom.Wrapf("denom=%s", coin.Denom)
		}
		seen[coin.Denom] = struct{}{}
	}
	return nil
}
