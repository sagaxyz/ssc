package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmdb "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/sagaxyz/ssc/x/billing/keeper"
	"github.com/sagaxyz/ssc/x/billing/testutil"
	"github.com/sagaxyz/ssc/x/billing/types"
	chainlettypes "github.com/sagaxyz/ssc/x/chainlet/types"
	epochstypes "github.com/sagaxyz/ssc/x/epochs/types"
)

// setupKeeperWithMocks creates a billing keeper with mocked dependencies
func setupKeeperWithMocks(t *testing.T) (*keeper.Keeper, sdk.Context, *testutil.MockChainletKeeper, *testutil.MockEpochsKeeper) {
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
		"BillingParams",
	)

	ctrl := gomock.NewController(t)
	mockChainletKeeper := testutil.NewMockChainletKeeper(ctrl)
	mockEpochsKeeper := testutil.NewMockEpochsKeeper(ctrl)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		paramsSubspace,
		nil, // bankkeeper
		nil, // escrowkeeper
		nil, // accountkeeper
		nil, // stakingkeeper
		mockChainletKeeper,
		mockEpochsKeeper,
		"",
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params with default billing epoch
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx, mockChainletKeeper, mockEpochsKeeper
}

// TestBeforeEpochStart_SkipsNonBillingEpoch verifies that BeforeEpochStart
// skips processing when the epoch identifier doesn't match BillingEpoch
func TestBeforeEpochStart_SkipsNonBillingEpoch(t *testing.T) {
	k, ctx, mockChainletKeeper, mockEpochsKeeper := setupKeeperWithMocks(t)

	// Set billing epoch to "day"
	params := k.GetParams(ctx)
	params.BillingEpoch = "day"
	k.SetParams(ctx, params)

	// Test with different epoch identifiers that should be skipped
	epochIdentifiers := []string{"minute", "hour", "week"}

	for _, epochID := range epochIdentifiers {
		t.Run("skips_"+epochID, func(t *testing.T) {
			// Mock should NOT be called since we return early
			// No expectations set means any call would fail the test

			// Call BeforeEpochStart with non-matching epoch identifier
			err := k.BeforeEpochStart(ctx, epochID, 1)

			// Should return nil without error
			require.NoError(t, err)

			// Verify no chainlet keeper methods were called
			// (gomock will fail if any unexpected calls were made)
		})
	}

	// Ensure mocks weren't called
	mockChainletKeeper.EXPECT().ListChainletStack(gomock.Any(), gomock.Any()).Times(0)
	mockEpochsKeeper.EXPECT().GetEpochInfo(gomock.Any(), gomock.Any()).Times(0)
}

// TestBeforeEpochStart_ProcessesBillingEpoch verifies that BeforeEpochStart
// processes billing when the epoch identifier matches BillingEpoch
func TestBeforeEpochStart_ProcessesBillingEpoch(t *testing.T) {
	k, ctx, mockChainletKeeper, mockEpochsKeeper := setupKeeperWithMocks(t)

	// Set billing epoch to "day"
	params := k.GetParams(ctx)
	params.BillingEpoch = "day"
	k.SetParams(ctx, params)

	// Mock chainlet keeper GetParams (called before ListChainlets)
	mockChainletKeeper.EXPECT().
		GetParams(gomock.Any()).
		Return(chainlettypes.Params{MaxChainlets: 100}).
		Times(1)

	// Mock chainlet stack list (empty to avoid further processing)
	mockChainletKeeper.EXPECT().
		ListChainletStack(gomock.Any(), gomock.Any()).
		Return(&chainlettypes.QueryListChainletStackResponse{
			ChainletStacks: []*chainlettypes.ChainletStack{},
		}, nil).
		Times(1)

	// Mock chainlet list (empty to avoid further processing)
	mockChainletKeeper.EXPECT().
		ListChainlets(gomock.Any(), gomock.Any()).
		Return(&chainlettypes.QueryListChainletsResponse{
			Chainlets: []*chainlettypes.Chainlet{},
		}, nil).
		Times(1)

	// Mock epoch info
	mockEpochsKeeper.EXPECT().
		GetEpochInfo(gomock.Any(), "day").
		Return(epochstypes.EpochInfo{
			Identifier:              "day",
			CurrentEpoch:            1,
			CurrentEpochStartTime:   time.Now(),
			CurrentEpochStartHeight: 1,
		}).
		Times(1)

	// Call BeforeEpochStart with matching epoch identifier
	err := k.BeforeEpochStart(ctx, "day", 1)

	// Should return nil without error
	require.NoError(t, err)
}

// TestBeforeEpochStart_WithDifferentBillingEpochs tests that the hook
// correctly filters based on the configured BillingEpoch parameter
func TestBeforeEpochStart_WithDifferentBillingEpochs(t *testing.T) {
	testCases := []struct {
		name              string
		billingEpoch      string
		callEpochID       string
		shouldProcess     bool
		expectedCallCount int
	}{
		{
			name:              "billing_day_called_with_day",
			billingEpoch:      "day",
			callEpochID:       "day",
			shouldProcess:     true,
			expectedCallCount: 1,
		},
		{
			name:              "billing_day_called_with_hour",
			billingEpoch:      "day",
			callEpochID:       "hour",
			shouldProcess:     false,
			expectedCallCount: 0,
		},
		{
			name:              "billing_hour_called_with_hour",
			billingEpoch:      "hour",
			callEpochID:       "hour",
			shouldProcess:     true,
			expectedCallCount: 1,
		},
		{
			name:              "billing_hour_called_with_minute",
			billingEpoch:      "hour",
			callEpochID:       "minute",
			shouldProcess:     false,
			expectedCallCount: 0,
		},
		{
			name:              "billing_week_called_with_week",
			billingEpoch:      "week",
			callEpochID:       "week",
			shouldProcess:     true,
			expectedCallCount: 1,
		},
		{
			name:              "billing_week_called_with_day",
			billingEpoch:      "week",
			callEpochID:       "day",
			shouldProcess:     false,
			expectedCallCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			k, ctx, mockChainletKeeper, mockEpochsKeeper := setupKeeperWithMocks(t)

			// Set billing epoch
			params := k.GetParams(ctx)
			params.BillingEpoch = tc.billingEpoch
			k.SetParams(ctx, params)

			if tc.shouldProcess {
				// Mock chainlet keeper GetParams (called before ListChainlets)
				mockChainletKeeper.EXPECT().
					GetParams(gomock.Any()).
					Return(chainlettypes.Params{MaxChainlets: 100}).
					Times(1)

				// Mock chainlet stack list
				mockChainletKeeper.EXPECT().
					ListChainletStack(gomock.Any(), gomock.Any()).
					Return(&chainlettypes.QueryListChainletStackResponse{
						ChainletStacks: []*chainlettypes.ChainletStack{},
					}, nil).
					Times(1)

				// Mock chainlet list
				mockChainletKeeper.EXPECT().
					ListChainlets(gomock.Any(), gomock.Any()).
					Return(&chainlettypes.QueryListChainletsResponse{
						Chainlets: []*chainlettypes.Chainlet{},
					}, nil).
					Times(1)

				// Mock epoch info
				mockEpochsKeeper.EXPECT().
					GetEpochInfo(gomock.Any(), tc.callEpochID).
					Return(epochstypes.EpochInfo{
						Identifier:              tc.callEpochID,
						CurrentEpoch:            1,
						CurrentEpochStartTime:   time.Now(),
						CurrentEpochStartHeight: 1,
					}).
					Times(1)
			} else {
				// No mocks should be called when skipping
				mockChainletKeeper.EXPECT().
					ListChainletStack(gomock.Any(), gomock.Any()).
					Times(0)
				mockChainletKeeper.EXPECT().
					ListChainlets(gomock.Any(), gomock.Any()).
					Times(0)
				mockEpochsKeeper.EXPECT().
					GetEpochInfo(gomock.Any(), gomock.Any()).
					Times(0)
			}

			// Call BeforeEpochStart
			err := k.BeforeEpochStart(ctx, tc.callEpochID, 1)
			require.NoError(t, err)
		})
	}
}

// TestBeforeEpochStart_DefaultBillingEpoch verifies that the default
// billing epoch (day) works correctly
func TestBeforeEpochStart_DefaultBillingEpoch(t *testing.T) {
	k, ctx, mockChainletKeeper, mockEpochsKeeper := setupKeeperWithMocks(t)

	// Use default params (BillingEpoch should be "day")
	params := k.GetParams(ctx)
	require.Equal(t, types.SAGA_EPOCH_IDENTIFIER, params.BillingEpoch)

	// Mock chainlet keeper GetParams (called before ListChainlets)
	mockChainletKeeper.EXPECT().
		GetParams(gomock.Any()).
		Return(chainlettypes.Params{MaxChainlets: 100}).
		Times(1)

	// Mock chainlet stack list
	mockChainletKeeper.EXPECT().
		ListChainletStack(gomock.Any(), gomock.Any()).
		Return(&chainlettypes.QueryListChainletStackResponse{
			ChainletStacks: []*chainlettypes.ChainletStack{},
		}, nil).
		Times(1)

	// Mock chainlet list
	mockChainletKeeper.EXPECT().
		ListChainlets(gomock.Any(), gomock.Any()).
		Return(&chainlettypes.QueryListChainletsResponse{
			Chainlets: []*chainlettypes.Chainlet{},
		}, nil).
		Times(1)

	// Mock epoch info
	mockEpochsKeeper.EXPECT().
		GetEpochInfo(gomock.Any(), "day").
		Return(epochstypes.EpochInfo{
			Identifier:              "day",
			CurrentEpoch:            1,
			CurrentEpochStartTime:   time.Now(),
			CurrentEpochStartHeight: 1,
		}).
		Times(1)

	// Call with matching epoch identifier
	err := k.BeforeEpochStart(ctx, "day", 1)
	require.NoError(t, err)

	// Call with non-matching epoch identifier (should skip)
	err = k.BeforeEpochStart(ctx, "hour", 1)
	require.NoError(t, err)
}

// TestBeforeEpochStart_MirrorsAfterEpochEndPattern verifies that
// BeforeEpochStart follows the same pattern as AfterEpochEnd
func TestBeforeEpochStart_MirrorsAfterEpochEndPattern(t *testing.T) {
	k, ctx, _, _ := setupKeeperWithMocks(t)

	// Set billing epoch to "day"
	params := k.GetParams(ctx)
	params.BillingEpoch = "day"
	k.SetParams(ctx, params)

	// Test that BeforeEpochStart returns early for non-matching epochs
	// (similar to AfterEpochEnd's early return pattern)
	err := k.BeforeEpochStart(ctx, "minute", 1)
	require.NoError(t, err)

	err = k.BeforeEpochStart(ctx, "hour", 1)
	require.NoError(t, err)

	err = k.BeforeEpochStart(ctx, "week", 1)
	require.NoError(t, err)

	// All should return nil without processing
}

