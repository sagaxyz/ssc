package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmdb "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"

	"github.com/sagaxyz/ssc/x/escrow/types"
)

// testKeeper creates a minimal keeper for unit tests without external dependencies
func testKeeper(t *testing.T) (*Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"EscrowParams",
	)

	k := NewKeeper(cdc, storeKey, paramsSubspace, nil, nil, nil, nil)
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}

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

func TestClearPoolFunders(t *testing.T) {
	k, ctx := testKeeper(t)

	chainID := "test-chain"
	denom := "utoken"

	// Set up multiple funders for the pool
	funders := []string{
		"saga1abc123",
		"saga1def456",
		"saga1ghi789",
	}

	for _, addr := range funders {
		k.setFunder(ctx, chainID, denom, addr, types.Funder{
			Shares: math.LegacyNewDec(100),
		})
	}

	// Verify funders exist
	for _, addr := range funders {
		_, exists := k.getFunder(ctx, chainID, denom, addr)
		require.True(t, exists, "funder %s should exist before clearing", addr)
	}

	// Clear all funders
	k.clearPoolFunders(ctx, chainID, denom)

	// Verify all funders are removed
	for _, addr := range funders {
		_, exists := k.getFunder(ctx, chainID, denom, addr)
		require.False(t, exists, "funder %s should not exist after clearing", addr)
	}

	// Verify reverse index is also cleared
	store := ctx.KVStore(k.storeKey)
	for _, addr := range funders {
		key := types.ByFunderKey(addr, chainID, denom)
		require.False(t, store.Has(key), "reverse index for %s should be cleared", addr)
	}
}

func TestClearPoolFunders_OnlyAffectsSpecificPool(t *testing.T) {
	k, ctx := testKeeper(t)

	chainID1 := "chain-1"
	chainID2 := "chain-2"
	denom1 := "utoken"
	denom2 := "usaga"

	// Set up funders across different pools
	k.setFunder(ctx, chainID1, denom1, "addr1", types.Funder{Shares: math.LegacyNewDec(100)})
	k.setFunder(ctx, chainID1, denom2, "addr2", types.Funder{Shares: math.LegacyNewDec(100)})
	k.setFunder(ctx, chainID2, denom1, "addr3", types.Funder{Shares: math.LegacyNewDec(100)})

	// Clear only chainID1/denom1 pool
	k.clearPoolFunders(ctx, chainID1, denom1)

	// Verify only the targeted pool's funders are removed
	_, exists := k.getFunder(ctx, chainID1, denom1, "addr1")
	require.False(t, exists, "addr1 in chain-1/utoken should be cleared")

	_, exists = k.getFunder(ctx, chainID1, denom2, "addr2")
	require.True(t, exists, "addr2 in chain-1/usaga should NOT be cleared")

	_, exists = k.getFunder(ctx, chainID2, denom1, "addr3")
	require.True(t, exists, "addr3 in chain-2/utoken should NOT be cleared")
}

func TestClearPoolFunders_EmptyPool(t *testing.T) {
	k, ctx := testKeeper(t)

	// Clearing an empty pool should not panic
	k.clearPoolFunders(ctx, "nonexistent-chain", "utoken")

	// Verify no funders exist (nothing to verify, just ensure no panic)
	store := ctx.KVStore(k.storeKey)
	pfx := prefix.NewStore(store, types.FunderPrefix("nonexistent-chain", "utoken"))
	it := pfx.Iterator(nil, nil)
	defer it.Close()
	require.False(t, it.Valid(), "should have no funders")
}

func TestDepositIntoDrainedPool_SharePricingVulnerability(t *testing.T) {
	// This test documents the vulnerability that was fixed:
	// When a pool's balance is drained to zero but shares remain,
	// new deposits would get 1:1 shares, allowing existing shareholders
	// to claim a portion of new deposits.
	//
	// The fix ensures that when balance hits zero, all funders are
	// cleared and shares are reset to zero.

	k, ctx := testKeeper(t)

	chainID := "test-chain"
	denom := "utoken"

	// Test case 1: Balance exactly zero - should clear
	t.Run("balance_zero_clears_funders", func(t *testing.T) {
		drainedPool := types.DenomPool{
			ChainId: chainID,
			Denom:   denom,
			Balance: sdk.NewCoin(denom, math.ZeroInt()),
			Shares:  math.LegacyNewDec(1000),
		}
		k.setFunder(ctx, chainID, denom, "old-funder", types.Funder{
			Shares: math.LegacyNewDec(1000),
		})
		k.setPool(ctx, drainedPool)

		// Simulate the check in BillAccount
		pool, _ := k.getPool(ctx, chainID, denom)
		if pool.Balance.IsZero() {
			k.clearPoolFunders(ctx, chainID, denom)
			pool.Shares = math.LegacyZeroDec()
			k.setPool(ctx, pool)
		}

		pool, _ = k.getPool(ctx, chainID, denom)
		require.True(t, pool.Shares.IsZero(), "pool shares should be zero")
		_, exists := k.getFunder(ctx, chainID, denom, "old-funder")
		require.False(t, exists, "old funder should be cleared")
	})

	// Test case 2: Balance > 0 - should NOT clear (proportional math works)
	t.Run("positive_balance_preserves_funders", func(t *testing.T) {
		chainID2 := "test-chain-2"
		lowBalancePool := types.DenomPool{
			ChainId: chainID2,
			Denom:   denom,
			Balance: sdk.NewCoin(denom, math.NewInt(1)),
			Shares:  math.LegacyNewDec(132323124),
		}
		k.setFunder(ctx, chainID2, denom, "old-funder-2", types.Funder{
			Shares: math.LegacyNewDec(132323124),
		})
		k.setPool(ctx, lowBalancePool)

		pool, _ := k.getPool(ctx, chainID2, denom)
		if pool.Balance.IsZero() {
			k.clearPoolFunders(ctx, chainID2, denom)
			pool.Shares = math.LegacyZeroDec()
			k.setPool(ctx, pool)
		}

		pool, _ = k.getPool(ctx, chainID2, denom)
		require.False(t, pool.Shares.IsZero(), "pool shares should NOT be zero when balance > 0")
		require.Equal(t, int64(1), pool.Balance.Amount.Int64(), "pool balance should remain")
		_, exists := k.getFunder(ctx, chainID2, denom, "old-funder-2")
		require.True(t, exists, "funder should NOT be cleared when balance > 0")
	})
}

func TestVulnerabilityScenario_EndToEnd(t *testing.T) {
	// End-to-end test of the vulnerability fix:
	// 1. Create initial funders with deposits
	// 2. Simulate BillAccount draining pool to zero (verify funders cleared)
	// 3. New depositor adds funds
	// 4. Verify new depositor receives 100% of shares (not diluted)

	k, ctx := testKeeper(t)

	chainID := "vuln-test-chain"
	denom := "utoken"

	// Step 1: Create initial pool with multiple funders
	initialPool := types.DenomPool{
		ChainId: chainID,
		Denom:   denom,
		Balance: sdk.NewCoin(denom, math.NewInt(10000)),
		Shares:  math.LegacyNewDec(10000),
	}
	k.setPool(ctx, initialPool)
	k.setChainlet(ctx, types.ChainletAccount{ChainId: chainID})

	// Add multiple funders
	k.setFunder(ctx, chainID, denom, "funder-alice", types.Funder{
		Shares: math.LegacyNewDec(6000),
	})
	k.setFunder(ctx, chainID, denom, "funder-bob", types.Funder{
		Shares: math.LegacyNewDec(4000),
	})

	// Verify initial state
	pool, _ := k.getPool(ctx, chainID, denom)
	require.Equal(t, int64(10000), pool.Balance.Amount.Int64())
	require.Equal(t, "10000.000000000000000000", pool.Shares.String())

	_, aliceExists := k.getFunder(ctx, chainID, denom, "funder-alice")
	_, bobExists := k.getFunder(ctx, chainID, denom, "funder-bob")
	require.True(t, aliceExists, "alice should exist initially")
	require.True(t, bobExists, "bob should exist initially")

	// Step 2: Simulate billing draining the pool to zero
	// This mimics what BillAccount does when it drains the pool
	pool.Balance = sdk.NewCoin(denom, math.ZeroInt())

	// Apply the fix: when balance is zero, clear funders and reset shares
	if pool.Balance.IsZero() {
		k.clearPoolFunders(ctx, chainID, denom)
		pool.Shares = math.LegacyZeroDec()
	}
	k.setPool(ctx, pool)

	// Verify funders were cleared
	pool, _ = k.getPool(ctx, chainID, denom)
	require.True(t, pool.Balance.IsZero(), "pool balance should be zero after drain")
	require.True(t, pool.Shares.IsZero(), "pool shares should be zero after drain")

	_, aliceExists = k.getFunder(ctx, chainID, denom, "funder-alice")
	_, bobExists = k.getFunder(ctx, chainID, denom, "funder-bob")
	require.False(t, aliceExists, "alice should be cleared after drain")
	require.False(t, bobExists, "bob should be cleared after drain")

	// Step 3: New depositor adds funds
	// Simulate deposit logic (without bank transfer since we don't have mock)
	newDepositAmount := math.NewInt(50000)
	var newShares math.LegacyDec

	pool, _ = k.getPool(ctx, chainID, denom)
	if pool.Balance.IsPositive() && pool.Shares.IsPositive() {
		// Proportional shares (won't execute since pool is empty)
		newShares = pool.Shares.MulInt(newDepositAmount).QuoInt(pool.Balance.Amount)
	} else {
		// Bootstrap 1:1 shares
		newShares = math.LegacyNewDecFromInt(newDepositAmount)
	}

	pool.Shares = pool.Shares.Add(newShares)
	pool.Balance = pool.Balance.Add(sdk.NewCoin(denom, newDepositAmount))
	k.setFunder(ctx, chainID, denom, "funder-new", types.Funder{Shares: newShares})
	k.setPool(ctx, pool)

	// Step 4: Verify new depositor receives 100% of shares
	pool, _ = k.getPool(ctx, chainID, denom)
	newFunder, newFunderExists := k.getFunder(ctx, chainID, denom, "funder-new")

	require.True(t, newFunderExists, "new funder should exist")
	require.Equal(t, int64(50000), pool.Balance.Amount.Int64(), "pool balance should be new deposit")
	require.Equal(t, "50000.000000000000000000", pool.Shares.String(), "pool shares should equal new deposit (1:1)")
	require.Equal(t, "50000.000000000000000000", newFunder.Shares.String(), "new funder should have 100% of shares")

	// Verify new funder owns 100% of the pool
	ownershipRatio := newFunder.Shares.Quo(pool.Shares)
	require.Equal(t, "1.000000000000000000", ownershipRatio.String(), "new funder should own 100% of pool")

	// Verify old funders cannot claim anything (they don't exist)
	_, aliceExists = k.getFunder(ctx, chainID, denom, "funder-alice")
	_, bobExists = k.getFunder(ctx, chainID, denom, "funder-bob")
	require.False(t, aliceExists, "alice should not exist after fix")
	require.False(t, bobExists, "bob should not exist after fix")
}

func TestDepositIntoPoolWithZeroShares(t *testing.T) {
	// Test edge case: Balance > 0 but Shares = 0
	// This shouldn't happen through normal operations but could occur
	// due to migration/genesis issues. Without the defensive check,
	// new depositors would get 0 shares (stolen deposit).

	k, ctx := testKeeper(t)

	chainID := "zero-shares-chain"
	denom := "utoken"

	// Setup: Invalid state where balance exists but no shares
	invalidPool := types.DenomPool{
		ChainId: chainID,
		Denom:   denom,
		Balance: sdk.NewCoin(denom, math.NewInt(1000)), // Has balance
		Shares:  math.LegacyZeroDec(),                  // But no shares!
	}
	k.setPool(ctx, invalidPool)
	k.setChainlet(ctx, types.ChainletAccount{ChainId: chainID})

	// New depositor adds funds using the FIXED deposit logic
	newDepositAmount := math.NewInt(50000)
	var newShares math.LegacyDec

	pool, _ := k.getPool(ctx, chainID, denom)

	// FIXED LOGIC: checks both balance AND shares
	if pool.Balance.IsPositive() && pool.Shares.IsPositive() {
		newShares = pool.Shares.MulInt(newDepositAmount).QuoInt(pool.Balance.Amount)
	} else {
		// Bootstrap 1:1 because shares is zero
		newShares = math.LegacyNewDecFromInt(newDepositAmount)
	}

	pool.Shares = pool.Shares.Add(newShares)
	pool.Balance = pool.Balance.Add(sdk.NewCoin(denom, newDepositAmount))
	k.setFunder(ctx, chainID, denom, "new-funder", types.Funder{Shares: newShares})
	k.setPool(ctx, pool)

	// Verify new depositor gets proper shares (not zero!)
	pool, _ = k.getPool(ctx, chainID, denom)
	newFunder, _ := k.getFunder(ctx, chainID, denom, "new-funder")

	require.Equal(t, "50000.000000000000000000", newShares.String(), "new funder should get 1:1 shares")
	require.Equal(t, int64(51000), pool.Balance.Amount.Int64(), "pool balance should include both")
	require.Equal(t, "50000.000000000000000000", pool.Shares.String(), "pool shares should be new deposit only")

	// New funder owns 100% of shares
	ownershipRatio := newFunder.Shares.Quo(pool.Shares)
	require.Equal(t, "1.000000000000000000", ownershipRatio.String(), "new funder should own 100% of shares")

	t.Logf("Balance > 0, Shares = 0: New depositor correctly gets 1:1 shares")
	t.Logf("  Pool balance: %d (includes 1000 orphaned + 50000 new)", pool.Balance.Amount.Int64())
	t.Logf("  Pool shares: %s", pool.Shares.String())
	t.Logf("  New funder shares: %s (100%% ownership)", newFunder.Shares.String())
}

func TestDepositIntoPoolWithZeroBalance(t *testing.T) {
	// Test edge case: Balance = 0 but Shares > 0
	// This is the main vulnerability scenario. Without the fix in BillAccount,
	// this state would persist and new depositors would be diluted.

	k, ctx := testKeeper(t)

	chainID := "zero-balance-chain"
	denom := "utoken"

	// Setup: Vulnerable state where shares exist but balance is zero
	vulnerablePool := types.DenomPool{
		ChainId: chainID,
		Denom:   denom,
		Balance: sdk.NewCoin(denom, math.ZeroInt()), // Drained to zero
		Shares:  math.LegacyNewDec(10000),           // But shares still exist!
	}
	k.setPool(ctx, vulnerablePool)
	k.setChainlet(ctx, types.ChainletAccount{ChainId: chainID})

	// Old funder still has shares in this vulnerable state
	k.setFunder(ctx, chainID, denom, "old-funder", types.Funder{
		Shares: math.LegacyNewDec(10000),
	})

	// Simulate the FIX: BillAccount clears funders when balance hits zero
	pool, _ := k.getPool(ctx, chainID, denom)
	if pool.Balance.IsZero() {
		k.clearPoolFunders(ctx, chainID, denom)
		pool.Shares = math.LegacyZeroDec()
		k.setPool(ctx, pool)
	}

	// Verify old funder was cleared
	_, oldFunderExists := k.getFunder(ctx, chainID, denom, "old-funder")
	require.False(t, oldFunderExists, "old funder should be cleared")

	// Now new depositor adds funds
	newDepositAmount := math.NewInt(50000)
	var newShares math.LegacyDec

	pool, _ = k.getPool(ctx, chainID, denom)

	if pool.Balance.IsPositive() && pool.Shares.IsPositive() {
		newShares = pool.Shares.MulInt(newDepositAmount).QuoInt(pool.Balance.Amount)
	} else {
		// Bootstrap 1:1 because pool was reset
		newShares = math.LegacyNewDecFromInt(newDepositAmount)
	}

	pool.Shares = pool.Shares.Add(newShares)
	pool.Balance = pool.Balance.Add(sdk.NewCoin(denom, newDepositAmount))
	k.setFunder(ctx, chainID, denom, "new-funder", types.Funder{Shares: newShares})
	k.setPool(ctx, pool)

	// Verify new depositor gets 100% ownership
	pool, _ = k.getPool(ctx, chainID, denom)
	newFunder, _ := k.getFunder(ctx, chainID, denom, "new-funder")

	require.Equal(t, "50000.000000000000000000", newShares.String(), "new funder should get 1:1 shares")
	require.Equal(t, int64(50000), pool.Balance.Amount.Int64(), "pool balance should be new deposit")
	require.Equal(t, "50000.000000000000000000", pool.Shares.String(), "pool shares should equal new deposit")

	ownershipRatio := newFunder.Shares.Quo(pool.Shares)
	require.Equal(t, "1.000000000000000000", ownershipRatio.String(), "new funder should own 100% of pool")

	// Old funder cannot claim anything
	_, oldFunderExists = k.getFunder(ctx, chainID, denom, "old-funder")
	require.False(t, oldFunderExists, "old funder should still not exist")

	t.Logf("Balance = 0, Shares > 0: Fix clears old funders, new depositor gets 100%%")
	t.Logf("  Pool balance: %d", pool.Balance.Amount.Int64())
	t.Logf("  Pool shares: %s", pool.Shares.String())
	t.Logf("  New funder shares: %s (100%% ownership)", newFunder.Shares.String())
}

