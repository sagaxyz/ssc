package keeper_test

import (
	"testing"

	testkeeper "github.com/sagaxyz/ssc/testutil/keeper"
	"github.com/sagaxyz/ssc/x/peers/types"
	"github.com/stretchr/testify/require"
)

func TestStoreData(t *testing.T) {
	k, ctx := testkeeper.PeersKeeper(t)

	// Test storing data
	chainId := "test-chain"
	addr := "test-addr"
	data := types.Data{
		Addresses: []string{"peer1", "peer2"},
	}

	// Store data
	err := k.StoreData(ctx, chainId, addr, data)
	require.NoError(t, err)

	// Verify data was stored
	storedData, err := k.Data(ctx, chainId, addr)
	require.NoError(t, err)
	require.Equal(t, data.Addresses, storedData.Addresses)

	// Verify reverse index was created
	kvStore := ctx.KVStore(k.GetStoreKey())
	chainIdsKey := append(types.ValidatorChainsPrefix, []byte(addr)...)
	require.True(t, kvStore.Has(chainIdsKey))
	require.Equal(t, chainId, string(kvStore.Get(chainIdsKey)))

	// Test storing data for same validator in different chain
	chainId2 := "test-chain-2"
	err = k.StoreData(ctx, chainId2, addr, data)
	require.NoError(t, err)

	// Verify reverse index was updated
	chainIds := string(kvStore.Get(chainIdsKey))
	require.Contains(t, chainIds, chainId)
	require.Contains(t, chainIds, chainId2)
}

func TestDeleteData(t *testing.T) {
	k, ctx := testkeeper.PeersKeeper(t)

	// Setup test data
	chainId := "test-chain"
	addr := "test-addr"
	data := types.Data{
		Addresses: []string{"peer1", "peer2"},
	}

	// Store data
	err := k.StoreData(ctx, chainId, addr, data)
	require.NoError(t, err)

	// Delete data
	err = k.DeleteData(ctx, chainId, addr)
	require.NoError(t, err)

	// Verify data was deleted
	_, err = k.Data(ctx, chainId, addr)
	require.Error(t, err)

	// Verify reverse index was updated
	kvStore := ctx.KVStore(k.GetStoreKey())
	chainIdsKey := append(types.ValidatorChainsPrefix, []byte(addr)...)
	require.False(t, kvStore.Has(chainIdsKey))

	// Test deleting non-existent data
	err = k.DeleteData(ctx, "non-existent", addr)
	require.NoError(t, err)
}

func TestDeleteAllValidatorData(t *testing.T) {
	k, ctx := testkeeper.PeersKeeper(t)

	// Setup test data
	addr := "test-addr"
	data := types.Data{
		Addresses: []string{"peer1", "peer2"},
	}

	// Store data in multiple chains
	chainIds := []string{"chain1", "chain2", "chain3"}
	for _, chainId := range chainIds {
		err := k.StoreData(ctx, chainId, addr, data)
		require.NoError(t, err)
	}

	// Delete all validator data
	err := k.DeleteAllValidatorData(ctx, addr)
	require.NoError(t, err)

	// Verify all data was deleted
	for _, chainId := range chainIds {
		_, err := k.Data(ctx, chainId, addr)
		require.Error(t, err)
	}

	// Verify reverse index was deleted
	kvStore := ctx.KVStore(k.GetStoreKey())
	chainIdsKey := append(types.ValidatorChainsPrefix, []byte(addr)...)
	require.False(t, kvStore.Has(chainIdsKey))

	// Test deleting non-existent validator
	err = k.DeleteAllValidatorData(ctx, "non-existent")
	require.NoError(t, err)
}
