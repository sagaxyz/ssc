package types_test

import (
	"testing"
	"time"

	"github.com/sagaxyz/ssc/x/peers/types"

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
				Params:        types.DefaultParams(),
				PeerData:      []types.GenesisPeerData{},
				ChainCounters: []types.GenesisChainCounter{},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "valid genesis state with data",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				PeerData: []types.GenesisPeerData{
					{ChainId: "chain-1", ValidatorAddress: "val1", Data: types.Data{Updated: time.Now(), Addresses: []string{"peer1"}}},
					{ChainId: "chain-1", ValidatorAddress: "val2", Data: types.Data{Updated: time.Now(), Addresses: []string{"peer2"}}},
				},
				ChainCounters: []types.GenesisChainCounter{
					{ChainId: "chain-1", Counter: types.Counter{Number: 2}},
				},
			},
			valid: true,
		},
		{
			desc: "invalid - duplicate peer data",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				PeerData: []types.GenesisPeerData{
					{ChainId: "chain-1", ValidatorAddress: "val1", Data: types.Data{Updated: time.Now(), Addresses: []string{"peer1"}}},
					{ChainId: "chain-1", ValidatorAddress: "val1", Data: types.Data{Updated: time.Now(), Addresses: []string{"peer2"}}},
				},
				ChainCounters: []types.GenesisChainCounter{},
			},
			valid: false,
		},
		{
			desc: "invalid - duplicate chain counters",
			genState: &types.GenesisState{
				Params:   types.DefaultParams(),
				PeerData: []types.GenesisPeerData{},
				ChainCounters: []types.GenesisChainCounter{
					{ChainId: "chain-1", Counter: types.Counter{Number: 2}},
					{ChainId: "chain-1", Counter: types.Counter{Number: 3}},
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
