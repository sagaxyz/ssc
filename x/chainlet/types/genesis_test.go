package types_test

import (
	"testing"

	"github.com/sagaxyz/ssc/x/chainlet/types"

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
				Params: types.Params{
					ChainletStackProtections:         false,
					NEpochDeposit:                    "30",
					AutomaticChainletUpgrades:        true,
					AutomaticChainletUpgradeInterval: 100,
				},
				PortId:         "chainlet",
				Chainlets:      []types.Chainlet{},
				ChainletStacks: []types.ChainletStack{},
				ChainletCount:  0,
			},
			valid: true,
		},
		{
			desc: "valid genesis state with chainlets",
			genState: &types.GenesisState{
				Params: types.Params{
					ChainletStackProtections:         false,
					NEpochDeposit:                    "30",
					AutomaticChainletUpgrades:        true,
					AutomaticChainletUpgradeInterval: 100,
				},
				PortId: "chainlet",
				Chainlets: []types.Chainlet{
					{ChainId: "chain-1"},
					{ChainId: "chain-2"},
				},
				ChainletStacks: []types.ChainletStack{
					{DisplayName: "stack-1"},
				},
				ChainletCount: 2,
			},
			valid: true,
		},
		{
			desc: "invalid genesis state - duplicate chainlet IDs",
			genState: &types.GenesisState{
				Params: types.Params{
					ChainletStackProtections:         false,
					NEpochDeposit:                    "30",
					AutomaticChainletUpgrades:        true,
					AutomaticChainletUpgradeInterval: 100,
				},
				PortId: "chainlet",
				Chainlets: []types.Chainlet{
					{ChainId: "chain-1"},
					{ChainId: "chain-1"}, // duplicate
				},
				ChainletStacks: []types.ChainletStack{},
				ChainletCount:  2,
			},
			valid: false,
		},
		{
			desc: "invalid genesis state - duplicate stack names",
			genState: &types.GenesisState{
				Params: types.Params{
					ChainletStackProtections:         false,
					NEpochDeposit:                    "30",
					AutomaticChainletUpgrades:        true,
					AutomaticChainletUpgradeInterval: 100,
				},
				PortId:    "chainlet",
				Chainlets: []types.Chainlet{},
				ChainletStacks: []types.ChainletStack{
					{DisplayName: "stack-1"},
					{DisplayName: "stack-1"}, // duplicate
				},
				ChainletCount: 0,
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
