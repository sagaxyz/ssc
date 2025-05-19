package keeper

import (
	"fmt"
	"strings"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/peers/types"
)

// Data returns the data for a given address and chain ID
func (k Keeper) Data(ctx sdk.Context, chainId string, addr string) (data types.Data, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DataKey)
	store = prefix.NewStore(store, types.KeyPrefix(chainId))

	key := []byte(addr)
	if !store.Has(key) {
		err = fmt.Errorf("addr %s has no data for chainlet %s", addr, chainId)
		return
	}

	b := store.Get(key)
	k.cdc.MustUnmarshal(b, &data)
	return
}

// StoreData stores data for a given address and chain ID
func (k Keeper) StoreData(ctx sdk.Context, chainId string, addr string, data types.Data) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DataKey)
	store = prefix.NewStore(store, types.KeyPrefix(chainId))
	key := []byte(addr)
	value := k.cdc.MustMarshal(&data)
	store.Set(key, value)

	// Update reverse index
	kvStore := ctx.KVStore(k.storeKey)
	chainIdsKey := validatorChainIdsKey(addr)

	// Get existing chain IDs
	var chainIds []string
	if kvStore.Has(chainIdsKey) {
		chainIds = strings.Split(string(kvStore.Get(chainIdsKey)), ",")
	}

	// Add new chain ID if not present
	found := false
	for _, id := range chainIds {
		if id == chainId {
			found = true
			break
		}
	}
	if !found {
		chainIds = append(chainIds, chainId)
		kvStore.Set(chainIdsKey, []byte(strings.Join(chainIds, ",")))
	}

	return nil
}

// DeleteData deletes data for a given address and chain ID
func (k Keeper) DeleteData(ctx sdk.Context, chainId string, addr string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DataKey)
	store = prefix.NewStore(store, types.KeyPrefix(chainId))
	key := []byte(addr)
	store.Delete(key)

	// Update reverse index
	kvStore := ctx.KVStore(k.storeKey)
	chainIdsKey := validatorChainIdsKey(addr)

	if kvStore.Has(chainIdsKey) {
		chainIds := strings.Split(string(kvStore.Get(chainIdsKey)), ",")

		// Remove chain ID from list
		newChainIds := make([]string, 0, len(chainIds))
		for _, id := range chainIds {
			if id != chainId {
				newChainIds = append(newChainIds, id)
			}
		}

		if len(newChainIds) == 0 {
			// If no more chain IDs, delete the entry
			kvStore.Delete(chainIdsKey)
		} else {
			// Otherwise update with remaining chain IDs
			kvStore.Set(chainIdsKey, []byte(strings.Join(newChainIds, ",")))
		}
	}

	return nil
}

// DeleteAllValidatorData deletes all data for a given validator address
func (k Keeper) DeleteAllValidatorData(ctx sdk.Context, addr string) error {
	kvStore := ctx.KVStore(k.storeKey)
	chainIdsKey := validatorChainIdsKey(addr)

	if !kvStore.Has(chainIdsKey) {
		return nil // No data to delete
	}

	chainIds := strings.Split(string(kvStore.Get(chainIdsKey)), ",")

	// Delete data for each chain ID
	for _, chainId := range chainIds {
		err := k.DeleteData(ctx, chainId, addr)
		if err != nil {
			return err
		}
	}

	// Delete the reverse index entry
	kvStore.Delete(chainIdsKey)
	return nil
}

// validatorChainIdsKey returns the key for storing a validator's chain IDs
func validatorChainIdsKey(addr string) []byte {
	return append(types.ValidatorChainsPrefix, []byte(addr)...)
}

func (k Keeper) peers(ctx sdk.Context, chainId string) (peers []string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DataKey)
	store = prefix.NewStore(store, types.KeyPrefix(chainId))

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var data types.Data
		k.cdc.MustUnmarshal(iterator.Value(), &data)

		peers = append(peers, data.Addresses...)
	}

	return
}
