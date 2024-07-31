package keeper

import (
	"fmt"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/peers/types"
)

//nolint:unused
func (k Keeper) data(ctx sdk.Context, chainId string, addr string) (data types.Data, err error) {
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

func (k Keeper) storeData(ctx sdk.Context, chainId string, addr string, data types.Data) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DataKey)
	store = prefix.NewStore(store, types.KeyPrefix(chainId))
	key := []byte(addr)
	value := k.cdc.MustMarshal(&data)
	store.Set(key, value)
	return nil
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
