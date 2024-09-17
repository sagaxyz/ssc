package keeper

import (
	"testing"

	tmdb "github.com/cosmos/cosmos-db"
	"cosmossdk.io/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"cosmossdk.io/store/metrics"
	"github.com/stretchr/testify/require"
	
	"github.com/sagaxyz/ssc/x/epochs"
	"github.com/sagaxyz/ssc/x/epochs/keeper"
	"github.com/sagaxyz/ssc/x/epochs/types"
)

func EpochsKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	k := keeper.NewKeeper(
		storeKey,
	)

	k.SetHooks(
		types.NewMultiEpochHooks(),
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	epochs.InitGenesis(ctx, *k, *types.DefaultGenesis())

	return k, ctx
}
