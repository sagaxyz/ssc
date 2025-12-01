package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/peers/types"
)

// ExportPeerData exports all peer data from the store
func (k Keeper) ExportPeerData(ctx sdk.Context) []types.GenesisPeerData {
	dataStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.DataKey)
	chainStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainsKey)

	var peerData []types.GenesisPeerData

	// Iterate over all chain IDs
	chainIterator := chainStore.Iterator(nil, nil)
	defer chainIterator.Close()

	for ; chainIterator.Valid(); chainIterator.Next() {
		chainID := string(chainIterator.Key())

		// For each chain, iterate over all validators
		chainDataStore := prefix.NewStore(dataStore, types.KeyPrefix(chainID))
		dataIterator := chainDataStore.Iterator(nil, nil)

		for ; dataIterator.Valid(); dataIterator.Next() {
			validatorAddr := string(dataIterator.Key())
			var data types.Data
			k.cdc.MustUnmarshal(dataIterator.Value(), &data)

			peerData = append(peerData, types.GenesisPeerData{
				ChainId:          chainID,
				ValidatorAddress: validatorAddr,
				Data:             data,
			})
		}
		dataIterator.Close()
	}

	return peerData
}

// ExportChainCounters exports all chain counters from the store
func (k Keeper) ExportChainCounters(ctx sdk.Context) []types.GenesisChainCounter {
	chainStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainsKey)

	var counters []types.GenesisChainCounter

	iterator := chainStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var counter types.Counter
		k.cdc.MustUnmarshal(iterator.Value(), &counter)

		counters = append(counters, types.GenesisChainCounter{
			ChainId: string(iterator.Key()),
			Counter: counter,
		})
	}

	return counters
}

// ImportPeerData imports a single peer data entry into the store
func (k Keeper) ImportPeerData(ctx sdk.Context, chainID, validatorAddr string, data types.Data) {
	dataStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.DataKey)
	chainDataStore := prefix.NewStore(dataStore, types.KeyPrefix(chainID))
	chainDataStore.Set([]byte(validatorAddr), k.cdc.MustMarshal(&data))
}

// ImportChainCounter imports a single chain counter into the store
func (k Keeper) ImportChainCounter(ctx sdk.Context, chainID string, counter types.Counter) {
	chainStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainsKey)
	chainStore.Set([]byte(chainID), k.cdc.MustMarshal(&counter))
}
