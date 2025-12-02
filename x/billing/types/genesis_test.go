package types_test

import (
	"testing"

	"github.com/sagaxyz/ssc/x/billing/types"
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
				Params:                 types.DefaultParams(),
				BillingHistory:         []types.SaveBillingHistory{},
				ValidatorPayoutHistory: []types.ValidatorPayoutHistory{},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "valid genesis state with data",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				BillingHistory: []types.SaveBillingHistory{
					{ChainletId: "chain-1", EpochIdentifier: "day", EpochNumber: 1, BilledAmount: "100usaga"},
					{ChainletId: "chain-1", EpochIdentifier: "day", EpochNumber: 2, BilledAmount: "100usaga"},
				},
				ValidatorPayoutHistory: []types.ValidatorPayoutHistory{
					{ValidatorAddress: "val1", EpochIdentifier: "day", EpochNumber: 1, RewardAmount: "50usaga"},
				},
			},
			valid: true,
		},
		{
			desc: "invalid - duplicate billing history",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				BillingHistory: []types.SaveBillingHistory{
					{ChainletId: "chain-1", EpochIdentifier: "day", EpochNumber: 1, BilledAmount: "100usaga"},
					{ChainletId: "chain-1", EpochIdentifier: "day", EpochNumber: 1, BilledAmount: "200usaga"},
				},
				ValidatorPayoutHistory: []types.ValidatorPayoutHistory{},
			},
			valid: false,
		},
		{
			desc: "invalid - duplicate validator payout history",
			genState: &types.GenesisState{
				Params:         types.DefaultParams(),
				BillingHistory: []types.SaveBillingHistory{},
				ValidatorPayoutHistory: []types.ValidatorPayoutHistory{
					{ValidatorAddress: "val1", EpochIdentifier: "day", EpochNumber: 1, RewardAmount: "50usaga"},
					{ValidatorAddress: "val1", EpochIdentifier: "day", EpochNumber: 1, RewardAmount: "100usaga"},
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
