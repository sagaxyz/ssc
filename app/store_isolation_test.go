package app_test

import (
	"testing"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	ibctransfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	icahosttypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host/types"
	"github.com/stretchr/testify/require"

	"github.com/sagaxyz/ssc/app"
)

// TestICAHostStoreIsolation verifies that ICA Host keeper uses its dedicated store
// and not the transfer store, preventing state collisions.
func TestICAHostStoreIsolation(t *testing.T) {
	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = app.DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = uint(1)

	logger := log.NewNopLogger()
	db := dbm.NewMemDB()

	bApp := app.New(
		logger,
		db,
		nil,
		true,
		appOptions,
		baseapp.SetChainID("test-chain"),
	)

	// Verify that ICA Host keeper is initialized
	require.NotNil(t, bApp.ICAHostKeeper, "ICAHostKeeper should be initialized")

	// Verify that Transfer keeper is initialized
	require.NotNil(t, bApp.TransferKeeper, "TransferKeeper should be initialized")

	// Get the store keys from the app
	icaHostStoreKey := bApp.GetKey(icahosttypes.StoreKey)
	transferStoreKey := bApp.GetKey(ibctransfertypes.StoreKey)

	// Verify that both store keys exist
	require.NotNil(t, icaHostStoreKey, "ICA Host store key should exist")
	require.NotNil(t, transferStoreKey, "Transfer store key should exist")

	// Verify that they are different store keys
	require.NotEqual(t, icaHostStoreKey, transferStoreKey, "ICA Host and Transfer should use different store keys")
	require.NotEqual(t, icaHostStoreKey.Name(), transferStoreKey.Name(), "Store key names should be different")

	// Verify store key names match expected values
	require.Equal(t, icahosttypes.StoreKey, icaHostStoreKey.Name(), "ICA Host store key name should match")
	require.Equal(t, ibctransfertypes.StoreKey, transferStoreKey.Name(), "Transfer store key name should match")
}

// TestStoreKeysUniqueness verifies that all IBC-related store keys are unique
func TestStoreKeysUniqueness(t *testing.T) {
	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = app.DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = uint(1)

	logger := log.NewNopLogger()
	db := dbm.NewMemDB()

	bApp := app.New(
		logger,
		db,
		nil,
		true,
		appOptions,
		baseapp.SetChainID("test-chain"),
	)

	// Collect all IBC-related store keys
	storeKeys := make(map[string]storetypes.StoreKey)
	
	// Get IBC-related store keys
	if key := bApp.GetKey(ibctransfertypes.StoreKey); key != nil {
		storeKeys[key.Name()] = key
	}
	if key := bApp.GetKey(icahosttypes.StoreKey); key != nil {
		storeKeys[key.Name()] = key
	}

	// Verify that we have both keys
	require.Len(t, storeKeys, 2, "Should have both Transfer and ICA Host store keys")

	// Verify they are different
	transferKey := storeKeys[ibctransfertypes.StoreKey]
	icaHostKey := storeKeys[icahosttypes.StoreKey]
	
	require.NotNil(t, transferKey, "Transfer store key should exist")
	require.NotNil(t, icaHostKey, "ICA Host store key should exist")
	require.NotEqual(t, transferKey, icaHostKey, "Transfer and ICA Host store keys must be different")
}

