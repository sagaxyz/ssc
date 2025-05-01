package types

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChainletParamsMarshalling(t *testing.T) {
	tests := []struct {
		name   string
		params ChainletParams
		err    error
	}{
		{
			name: "valid",
			params: ChainletParams{
				GasLimit:          1337,
				CreateEmptyBlocks: true,
				DacEnable:         true,
				GenAcctBalances: GenesisAccountBalances{
					List: []*AccountBalance{
						{
							Address: "cosmos1wze8mn5nsgl9qrgazq6a92fvh7m5e6psjcx2du",
							Balance: "123usaga",
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "valid - multiple genesis balances",
			params: ChainletParams{
				GasLimit:          1337,
				CreateEmptyBlocks: true,
				DacEnable:         true,
				GenAcctBalances: GenesisAccountBalances{
					List: []*AccountBalance{
						{
							Address: "cosmos1wze8mn5nsgl9qrgazq6a92fvh7m5e6psjcx2du",
							Balance: "123usaga",
						},
						{
							Address: "cosmos1wze8mn5nsgl9qrgazq6a92fvh7m5e6psjcx2du",
							Balance: "321usaga",
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "valid - no genesis balances",
			params: ChainletParams{
				GasLimit:          1337,
				CreateEmptyBlocks: true,
				DacEnable:         true,
				GenAcctBalances: GenesisAccountBalances{
					List: []*AccountBalance{},
				},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.params)
			require.NoError(t, err)

			var params ChainletParams
			err = json.Unmarshal(data, &params)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}

			t.Logf("expected params: %+v\n", tt.params)
			t.Logf("data: %s\n", data)
			t.Logf("actual params: %+v\n", params)

			require.True(t, reflect.DeepEqual(tt.params, params))
		})
	}
}
