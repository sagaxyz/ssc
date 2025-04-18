package utils

import (
	"testing"

	"github.com/test-go/testify/require"
)

func TestValidateConfig(t *testing.T) {
	testcases := []struct {
		name        string
		config      *networkConfig
		errContains string
	}{
		{
			name:   "pass",
			config: defaultConfig(),
		},
		{
			name: "fail - invalid chains",
			config: &networkConfig{
				nChains: 0,
			},
			errContains: "invalid number of chains",
		},
		{
			name: "fail - invalid validators",
			config: &networkConfig{
				nChains:       1,
				nValsPerChain: 0,
			},
			errContains: "invalid number of validators",
		},
		{
			name: "fail - invalid relayer path index",
			config: &networkConfig{
				nChains:       3,
				nValsPerChain: 1,
				relayerPaths:  []RelayerPath{{0, 1}, {3, 2}},
			},
			errContains: "relayer path contains invalid index: 3; max 2 allowed",
		},
		{
			name: "fail - invalid relayer path length",
			config: &networkConfig{
				nChains:       2,
				nValsPerChain: 1,
				relayerPaths:  []RelayerPath{{0, 1}, {3, 2}},
			},
			errContains: "incorrect number of relayer paths; expected max. of 1, got 2",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.validate()
			if tc.errContains != "" {
				require.Contains(t, err.Error(), tc.errContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
