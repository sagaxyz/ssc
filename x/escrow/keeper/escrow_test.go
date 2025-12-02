package keeper

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/sagaxyz/ssc/x/escrow/types"
)

func TestScalingFactor(t *testing.T) {
	tests := []struct {
		name     string
		pool     types.DenomPool
		expected math.LegacyDec
	}{
		{
			name: "1:1 ratio",
			pool: types.DenomPool{
				ChainId: "test-chain",
				Denom:   "utoken",
				Balance: sdk.NewCoin("utoken", math.NewInt(1000)),
				Shares:  math.LegacyNewDec(1000),
			},
			expected: math.LegacyOneDec(),
		},
		{
			name: "2:1 shares to tokens",
			pool: types.DenomPool{
				ChainId: "test-chain",
				Denom:   "utoken",
				Balance: sdk.NewCoin("utoken", math.NewInt(500)),
				Shares:  math.LegacyNewDec(1000),
			},
			expected: math.LegacyNewDec(2),
		},
		{
			name: "zero balance returns 1",
			pool: types.DenomPool{
				ChainId: "test-chain",
				Denom:   "utoken",
				Balance: sdk.NewCoin("utoken", math.NewInt(0)),
				Shares:  math.LegacyNewDec(0),
			},
			expected: math.LegacyOneDec(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sf := ScalingFactor(tc.pool)
			require.True(t, tc.expected.Equal(sf), "expected %s, got %s", tc.expected, sf)
		})
	}
}

func TestWithdrawTokenCalculation_UsesFloor(t *testing.T) {
	// This test verifies that token calculation uses floor (truncation) behavior,
	// not rounding. This is critical for security: users should never receive
	// more tokens than their proportional share.

	tests := []struct {
		name           string
		poolShares     math.LegacyDec
		poolBalance    math.Int
		funderShares   math.LegacyDec
		expectedTokens math.Int
	}{
		{
			name:           "exact division - no rounding needed",
			poolShares:     math.LegacyNewDec(100),
			poolBalance:    math.NewInt(100),
			funderShares:   math.LegacyNewDec(50),
			expectedTokens: math.NewInt(50),
		},
		{
			name:           "should floor, not round up",
			poolShares:     math.LegacyNewDec(100),
			poolBalance:    math.NewInt(100),
			funderShares:   math.LegacyMustNewDecFromStr("33.6"), // Would round to 34 with banker's rounding
			expectedTokens: math.NewInt(33),                      // Should be 33 (floor)
		},
		{
			name:           "edge case - 0.5 remainder should floor to 0",
			poolShares:     math.LegacyNewDec(100),
			poolBalance:    math.NewInt(100),
			funderShares:   math.LegacyMustNewDecFromStr("0.5"),
			expectedTokens: math.NewInt(0), // Should be 0 (floor), not 1 (round)
		},
		{
			name:           "complex ratio with truncation",
			poolShares:     math.LegacyNewDec(300),
			poolBalance:    math.NewInt(100),
			funderShares:   math.LegacyNewDec(100), // 100 shares = 100/300 * 100 = 33.33... tokens
			expectedTokens: math.NewInt(33),        // Should floor to 33
		},
		{
			name:           "high precision shares",
			poolShares:     math.LegacyMustNewDecFromStr("1000000000000000000"),
			poolBalance:    math.NewInt(1000000000),
			funderShares:   math.LegacyMustNewDecFromStr("999999999999999999"),
			expectedTokens: math.NewInt(999999999), // Should floor
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pool := types.DenomPool{
				ChainId: "test-chain",
				Denom:   "utoken",
				Balance: sdk.NewCoin("utoken", tc.poolBalance),
				Shares:  tc.poolShares,
			}

			sf := ScalingFactor(pool)
			require.False(t, sf.IsZero(), "scaling factor should not be zero")

			// This mirrors the actual calculation in withdrawOne():
			// tokensDec := f.Shares.QuoTruncate(sf)
			// amt := tokensDec.TruncateInt()
			tokensDec := tc.funderShares.QuoTruncate(sf)
			amt := tokensDec.TruncateInt()

			require.True(t, tc.expectedTokens.Equal(amt),
				"expected %s tokens, got %s (tokensDec: %s)",
				tc.expectedTokens, amt, tokensDec)

			// Verify that using Quo (which rounds) could give different results
			// This demonstrates why QuoTruncate is necessary
			tokenWithRounding := tc.funderShares.Quo(sf).TruncateInt()
			if !tc.expectedTokens.Equal(tokenWithRounding) {
				t.Logf("Note: Quo would have given %s instead of %s (QuoTruncate)",
					tokenWithRounding, amt)
			}
		})
	}
}

func TestWithdrawTokenCalculation_NeverExceedsEntitlement(t *testing.T) {
	// Property test: withdrawn tokens should never exceed the proportional share
	// tokens_out <= funder_shares / total_shares * total_balance

	testCases := []struct {
		poolShares   string
		poolBalance  int64
		funderShares string
	}{
		{"100", 100, "33"},
		{"100", 100, "33.5"},
		{"100", 100, "33.9"},
		{"100", 100, "34"},
		{"1000", 333, "500"},
		{"777", 999, "123.456789"},
	}

	for _, tc := range testCases {
		pool := types.DenomPool{
			ChainId: "test-chain",
			Denom:   "utoken",
			Balance: sdk.NewCoin("utoken", math.NewInt(tc.poolBalance)),
			Shares:  math.LegacyMustNewDecFromStr(tc.poolShares),
		}
		funderShares := math.LegacyMustNewDecFromStr(tc.funderShares)

		sf := ScalingFactor(pool)
		tokensDec := funderShares.QuoTruncate(sf)
		amt := tokensDec.TruncateInt()

		// Calculate exact entitlement: funderShares / poolShares * poolBalance
		exactEntitlement := funderShares.MulInt(pool.Balance.Amount).Quo(pool.Shares)

		// amt should be <= floor(exactEntitlement)
		require.True(t, math.LegacyNewDecFromInt(amt).LTE(exactEntitlement),
			"withdrawn amount %s exceeds entitlement %s for shares %s/%s on balance %d",
			amt, exactEntitlement, funderShares, tc.poolShares, tc.poolBalance)
	}
}
