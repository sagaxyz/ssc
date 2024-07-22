package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sagaxyz/ssc/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateChainletStack_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateChainletStack
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateChainletStack{
				Creator:     "invalid_address",
				DisplayName: "validname",
				Version:     "1234",
				Image:       "123",
				Checksum:    "1234",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateChainletStack{
				Creator:     sample.AccAddress(),
				DisplayName: "validname",
				Version:     "1234",
				Image:       "123",
				Checksum:    "1234",
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
