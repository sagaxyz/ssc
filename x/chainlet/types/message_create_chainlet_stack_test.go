package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sagaxyz/ssc/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateChainletStack_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateChainletStack
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateChainletStack{
				Creator:     "invalid_address",
				DisplayName: "validname",
				Version:     "1234",
				Image:       "123",
				Checksum:    "1234",
				Fees: ChainletStackFees{
					"1000upsaga",
					"day",
					"1000upsaga",
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateChainletStack{
				Creator:     sample.AccAddress(),
				DisplayName: "validname",
				Version:     "1234",
				Image:       "123",
				Checksum:    "1234",
				Fees: ChainletStackFees{
					"1000upsaga",
					"day",
					"1000upsaga",
				},
			},
		}, {
			name: "invalid epoch",
			msg: MsgCreateChainletStack{
				Creator:     sample.AccAddress(),
				DisplayName: "validname",
				Version:     "1234",
				Image:       "123",
				Checksum:    "1234",
				Fees: ChainletStackFees{
					"1000upsaga",
					"2",
					"1000upsaga",
				},
			},
			err: ErrInvalidEpoch,
		}, {
			name: "invalid setup fee",
			msg: MsgCreateChainletStack{
				Creator:     sample.AccAddress(),
				DisplayName: "validname",
				Version:     "1234",
				Image:       "123",
				Checksum:    "1234",
				Fees: ChainletStackFees{
					"1000upsaga",
					"day",
					"-100upsaga",
				},
			},
			err: ErrInvalidCoin,
		}, {
			name: "invalid epoch fee",
			msg: MsgCreateChainletStack{
				Creator:     sample.AccAddress(),
				DisplayName: "validname",
				Version:     "1234",
				Image:       "123",
				Checksum:    "1234",
				Fees: ChainletStackFees{
					"-1000upsaga",
					"day",
					"1000upsaga",
				},
			},
			err: ErrInvalidCoin,
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
