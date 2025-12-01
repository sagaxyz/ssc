package keeper

import (
	"bytes"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/escrow/types"
)

// ExportChainletAccounts exports all chainlet accounts from the store
func (k Keeper) ExportChainletAccounts(ctx sdk.Context) []types.ChainletAccount {
	store := ctx.KVStore(k.storeKey)
	iterator := prefix.NewStore(store, types.KeyChainletPrefix).Iterator(nil, nil)
	defer iterator.Close()

	var accounts []types.ChainletAccount
	for ; iterator.Valid(); iterator.Next() {
		var acc types.ChainletAccount
		k.cdc.MustUnmarshal(iterator.Value(), &acc)
		accounts = append(accounts, acc)
	}
	return accounts
}

// ExportPools exports all denomination pools from the store
func (k Keeper) ExportPools(ctx sdk.Context) []types.DenomPool {
	store := ctx.KVStore(k.storeKey)
	iterator := prefix.NewStore(store, types.KeyPoolPrefix).Iterator(nil, nil)
	defer iterator.Close()

	var pools []types.DenomPool
	for ; iterator.Valid(); iterator.Next() {
		var pool types.DenomPool
		k.cdc.MustUnmarshal(iterator.Value(), &pool)
		pools = append(pools, pool)
	}
	return pools
}

// ExportFunders exports all funder positions from the store
func (k Keeper) ExportFunders(ctx sdk.Context) []types.GenesisFunder {
	store := ctx.KVStore(k.storeKey)
	iterator := prefix.NewStore(store, types.KeyFunderPrefix).Iterator(nil, nil)
	defer iterator.Close()

	var funders []types.GenesisFunder
	for ; iterator.Valid(); iterator.Next() {
		var funder types.Funder
		k.cdc.MustUnmarshal(iterator.Value(), &funder)

		// Parse the key to extract chainId, denom, and address
		// Key format: {chainId}/{denom}/{addr}
		key := iterator.Key()
		parts := bytes.SplitN(key, []byte{'/'}, 3)
		if len(parts) != 3 {
			continue
		}

		funders = append(funders, types.GenesisFunder{
			ChainId: string(parts[0]),
			Denom:   string(parts[1]),
			Address: string(parts[2]),
			Funder:  funder,
		})
	}
	return funders
}

// ImportChainletAccount imports a single chainlet account into the store
func (k Keeper) ImportChainletAccount(ctx sdk.Context, acc types.ChainletAccount) {
	k.setChainlet(ctx, acc)
}

// ImportPool imports a single pool into the store
func (k Keeper) ImportPool(ctx sdk.Context, pool types.DenomPool) {
	k.setPool(ctx, pool)
}

// ImportFunder imports a single funder into the store (includes reverse index)
func (k Keeper) ImportFunder(ctx sdk.Context, chainID, denom, addr string, funder types.Funder) {
	k.setFunder(ctx, chainID, denom, addr, funder)
}
