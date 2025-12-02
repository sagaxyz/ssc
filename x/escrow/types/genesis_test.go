package types_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sagaxyz/ssc/x/escrow/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Params:           types.DefaultParams(),
				ChainletAccounts: []types.ChainletAccount{},
				Pools:            []types.DenomPool{},
				Funders:          []types.GenesisFunder{},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "valid genesis state with data",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ChainletAccounts: []types.ChainletAccount{
					{ChainId: "chain-1"},
					{ChainId: "chain-2"},
				},
				Pools: []types.DenomPool{
					{ChainId: "chain-1", Denom: "usaga", Balance: sdk.NewCoin("usaga", math.NewInt(1000)), Shares: math.LegacyNewDec(1000)},
				},
				Funders: []types.GenesisFunder{
					{ChainId: "chain-1", Denom: "usaga", Address: "saga1abc", Funder: types.Funder{Shares: math.LegacyNewDec(500)}},
				},
			},
			valid: true,
		},
		{
			desc: "invalid - duplicate chainlet accounts",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ChainletAccounts: []types.ChainletAccount{
					{ChainId: "chain-1"},
					{ChainId: "chain-1"},
				},
				Pools:   []types.DenomPool{},
				Funders: []types.GenesisFunder{},
			},
			valid: false,
		},
		{
			desc: "invalid - duplicate pools",
			genState: &types.GenesisState{
				Params:           types.DefaultParams(),
				ChainletAccounts: []types.ChainletAccount{},
				Pools: []types.DenomPool{
					{ChainId: "chain-1", Denom: "usaga", Balance: sdk.NewCoin("usaga", math.NewInt(1000)), Shares: math.LegacyNewDec(1000)},
					{ChainId: "chain-1", Denom: "usaga", Balance: sdk.NewCoin("usaga", math.NewInt(500)), Shares: math.LegacyNewDec(500)},
				},
				Funders: []types.GenesisFunder{},
			},
			valid: false,
		},
		{
			desc: "invalid - duplicate funders",
			genState: &types.GenesisState{
				Params:           types.DefaultParams(),
				ChainletAccounts: []types.ChainletAccount{},
				Pools:            []types.DenomPool{},
				Funders: []types.GenesisFunder{
					{ChainId: "chain-1", Denom: "usaga", Address: "saga1abc", Funder: types.Funder{Shares: math.LegacyNewDec(500)}},
					{ChainId: "chain-1", Denom: "usaga", Address: "saga1abc", Funder: types.Funder{Shares: math.LegacyNewDec(300)}},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
