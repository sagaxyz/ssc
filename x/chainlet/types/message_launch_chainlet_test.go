package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sagaxyz/ssc/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgLaunchChainlet_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgLaunchChainlet
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgLaunchChainlet{
				Creator:              "invalid_address",
				Maintainers:          []string{},
				ChainletStackName:    "sagaevm",
				ChainletStackVersion: "4.2",
				ChainletName:         "cassio",
				ChainId:              "cassio_123-4",
				Denom:                "asaga",
				Params:               ChainletParams{},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid",
			msg: MsgLaunchChainlet{
				Creator:              sample.AccAddress(),
				Maintainers:          []string{},
				ChainletStackName:    "sagaevm",
				ChainletStackVersion: "4.2",
				ChainletName:         "cassio",
				ChainId:              "cassio_123-4",
				Denom:                "asaga",
				Params: ChainletParams{
					GenAcctBalances: GenesisAccountBalances{
						List: []*AccountBalance{
							{
								Address: "cosmos1wze8mn5nsgl9qrgazq6a92fvh7m5e6psjcx2du",
								Balance: "123",
							},
						},
					},
				},
			},
		}, {
			name: "invalid chain id",
			msg: MsgLaunchChainlet{
				Creator:              sample.AccAddress(),
				Maintainers:          []string{},
				ChainletStackName:    "sagaevm",
				ChainletStackVersion: "4.2",
				ChainletName:         "cassio",
				ChainId:              "cassio_123",
				Denom:                "asaga",
				Params:               ChainletParams{},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "invalid chain id - no numbers in chain names",
			msg: MsgLaunchChainlet{
				Creator:              sample.AccAddress(),
				Maintainers:          []string{},
				ChainletStackName:    "sagaevm",
				ChainletStackVersion: "4.2",
				ChainletName:         "cassio",
				ChainId:              "cassio1_123-1",
				Denom:                "asaga",
				Params:               ChainletParams{},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "invalid chainlet stack",
			msg: MsgLaunchChainlet{
				Creator:              sample.AccAddress(),
				Maintainers:          []string{},
				ChainletStackName:    "",
				ChainletStackVersion: "4.2",
				ChainletName:         "cassio",
				ChainId:              "cassio_123-4",
				Denom:                "asaga",
				Params:               ChainletParams{},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "invalid chainlet stack version",
			msg: MsgLaunchChainlet{
				Creator:              sample.AccAddress(),
				Maintainers:          []string{},
				ChainletStackName:    "sagaevm",
				ChainletStackVersion: "",
				ChainletName:         "cassio",
				ChainId:              "cassio_123-4",
				Denom:                "asaga",
				Params:               ChainletParams{},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "invalid chainlet name",
			msg: MsgLaunchChainlet{
				Creator:              sample.AccAddress(),
				Maintainers:          []string{},
				ChainletStackName:    "sagaevm",
				ChainletStackVersion: "4.2",
				ChainletName:         "",
				ChainId:              "cassio_123-4",
				Denom:                "asaga",
				Params:               ChainletParams{},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "chainlet denom too long",
			msg: MsgLaunchChainlet{
				Creator:              sample.AccAddress(),
				Maintainers:          []string{},
				ChainletStackName:    "sagaevm",
				ChainletStackVersion: "4.2",
				ChainletName:         "cassio",
				ChainId:              "cassio_123-4",
				Denom:                "asaga",
				Params: ChainletParams{
					GenAcctBalances: GenesisAccountBalances{
						List: []*AccountBalance{
							{
								Address: "cosmos1wze8mn5nsgl9qrgazq6a92fvh7m5e6psjcx2du",
								Balance: "123",
							},
						},
					},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "chainlet denom too short",
			msg: MsgLaunchChainlet{
				Creator:              sample.AccAddress(),
				Maintainers:          []string{},
				ChainletStackName:    "sagaevm",
				ChainletStackVersion: "4.2",
				ChainletName:         "cassio",
				ChainId:              "cassio_123-4",
				Denom:                "asaga",
				Params: ChainletParams{
					GenAcctBalances: GenesisAccountBalances{
						List: []*AccountBalance{
							{
								Address: "cosmos1wze8mn5nsgl9qrgazq6a92fvh7m5e6psjcx2du",
								Balance: "123",
							},
						},
					},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "gas limit too high",
			msg: MsgLaunchChainlet{
				Creator:              sample.AccAddress(),
				Maintainers:          []string{},
				ChainletStackName:    "sagaevm",
				ChainletStackVersion: "4.2",
				ChainletName:         "cassio",
				ChainId:              "cassio_123-4",
				Denom:                "asaga",
				Params: ChainletParams{
					GenAcctBalances: GenesisAccountBalances{
						List: []*AccountBalance{
							{
								Address: "cosmos1wze8mn5nsgl9qrgazq6a92fvh7m5e6psjcx2du",
								Balance: "123",
							},
						},
					},
					GasLimit: 100000001,
				},
			},
			err: sdkerrors.ErrInvalidRequest,
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
