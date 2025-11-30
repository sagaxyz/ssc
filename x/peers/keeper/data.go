package keeper

import (
	"fmt"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/peers/types"
)

//nolint:unused
func (k Keeper) data(ctx sdk.Context, chainID string, addr string) (data types.Data, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DataKey)
	store = prefix.NewStore(store, types.KeyPrefix(chainID))

	key := []byte(addr)
	if !store.Has(key) {
		err = fmt.Errorf("addr %s has no data for chainlet %s", addr, chainID)
		return
	}

	b := store.Get(key)
	k.cdc.MustUnmarshal(b, &data)
	return
}

func (k Keeper) StoreData(ctx sdk.Context, chainID string, addr string, data types.Data) {
	var new bool
	addrKey := []byte(addr)
	dataStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.DataKey)
	dataStore = prefix.NewStore(dataStore, types.KeyPrefix(chainID))
	if !dataStore.Has(addrKey) {
		new = true
	}
	dataStore.Set(addrKey, k.cdc.MustMarshal(&data))

	// Increment the counter for this chain id
	if !new {
		return
	}
	chainStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainsKey)
	chainKey := []byte(chainID)
	var counter types.Counter
	if chainStore.Has(chainKey) {
		k.cdc.MustUnmarshal(chainStore.Get(chainKey), &counter)
	}
	// If chainKey doesn't exist, counter defaults to Number: 0
	counter.Number++
	chainStore.Set(chainKey, k.cdc.MustMarshal(&counter))
}

func (k Keeper) DeleteValidatorData(ctx sdk.Context, addr string) {
	chainStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainsKey)
	dataStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.DataKey)

	// Delete addr for each chain ID
	iterator := chainStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		chainID := iterator.Key()

		// Delete data for this addr
		var deleted bool
		addrKey := []byte(addr)
		s := prefix.NewStore(dataStore, chainID)
		if s.Has(addrKey) {
			deleted = true
		}
		s.Delete(addrKey)

		// Delete chain ID if there are no more entries for it
		if !deleted {
			continue
		}
		var counter types.Counter
		k.cdc.MustUnmarshal(iterator.Value(), &counter)
		counter.Number--
		if counter.Number > 0 {
			chainStore.Set(chainID, k.cdc.MustMarshal(&counter))
		} else {
			defer chainStore.Delete(chainID)
		}
	}
}

func (k Keeper) GetPeers(ctx sdk.Context, chainID string) (peers []string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DataKey)
	store = prefix.NewStore(store, types.KeyPrefix(chainID))

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	peers = make([]string, 0)
	for ; iterator.Valid(); iterator.Next() {
		var data types.Data
		k.cdc.MustUnmarshal(iterator.Value(), &data)

		peers = append(peers, data.Addresses...)
	}

	return
}

// Test helper
func (k Keeper) Counter(ctx sdk.Context, chainID string) uint32 {
	chainStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainsKey)

	if !chainStore.Has([]byte(chainID)) {
		return 0
	}
	var counter types.Counter
	k.cdc.MustUnmarshal(chainStore.Get([]byte(chainID)), &counter)
	if counter.Number == 0 {
		panic("stored counter at 0")
	}

	return counter.Number
}
