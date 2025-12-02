package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
)

// ExportChainlets exports all chainlets from the store
func (k *Keeper) ExportChainlets(ctx sdk.Context) []types.Chainlet {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletKey)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var chainlets []types.Chainlet
	for ; iterator.Valid(); iterator.Next() {
		var chainlet types.Chainlet
		k.cdc.MustUnmarshal(iterator.Value(), &chainlet)
		chainlets = append(chainlets, chainlet)
	}
	return chainlets
}

// ExportChainletStacks exports all chainlet stacks from the store
func (k *Keeper) ExportChainletStacks(ctx sdk.Context) []types.ChainletStack {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletStackKey)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var stacks []types.ChainletStack
	for ; iterator.Valid(); iterator.Next() {
		var stack types.ChainletStack
		k.cdc.MustUnmarshal(iterator.Value(), &stack)
		stacks = append(stacks, stack)
	}
	return stacks
}

// ImportChainlet imports a single chainlet into the store (without validation, for genesis import)
func (k *Keeper) ImportChainlet(ctx sdk.Context, chainlet types.Chainlet) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletKey)
	key := []byte(chainlet.ChainId)
	value := k.cdc.MustMarshal(&chainlet)
	store.Set(key, value)
	return nil
}

// ImportChainletStack imports a single chainlet stack into the store (without validation, for genesis import)
func (k *Keeper) ImportChainletStack(ctx sdk.Context, stack types.ChainletStack) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletStackKey)
	key := []byte(stack.DisplayName)
	value := k.cdc.MustMarshal(&stack)
	store.Set(key, value)

	// Also add enabled versions to the version tree
	for _, version := range stack.Versions {
		if version.Enabled {
			if err := k.AddVersion(ctx, stack.DisplayName, version); err != nil {
				return err
			}
		}
	}
	return nil
}
