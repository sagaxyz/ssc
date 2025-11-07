package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sagaxyz/ssc/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {
	defaultParams := DefaultParams()
	tests := []struct {
		name string
		msg  MsgUpdateParams
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateParams{
				Authority: "invalid_address",
				Params:    &defaultParams,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params:    &defaultParams,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
