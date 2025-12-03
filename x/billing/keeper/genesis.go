package keeper

import (
	"fmt"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/billing/types"
)

// ExportBillingHistory exports all billing history records from the store
func (k Keeper) ExportBillingHistory(ctx sdk.Context) []types.SaveBillingHistory {
	// Use prefix store to efficiently iterate only over billing history keys
	billingStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.BillingHistoryKey+"-"))
	iterator := billingStore.Iterator(nil, nil)
	defer iterator.Close()

	var records []types.SaveBillingHistory
	for ; iterator.Valid(); iterator.Next() {
		var record types.SaveBillingHistory
		k.cdc.MustUnmarshal(iterator.Value(), &record)
		records = append(records, record)
	}
	return records
}

// ExportValidatorPayoutHistory exports all validator payout history records from the store
func (k Keeper) ExportValidatorPayoutHistory(ctx sdk.Context) []types.ValidatorPayoutHistory {
	// Use prefix store to efficiently iterate only over validator payout history keys
	payoutStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ValidatorPayoutHistoryKey+"-"))
	iterator := payoutStore.Iterator(nil, nil)
	defer iterator.Close()

	var records []types.ValidatorPayoutHistory
	for ; iterator.Valid(); iterator.Next() {
		var record types.ValidatorPayoutHistory
		k.cdc.MustUnmarshal(iterator.Value(), &record)
		records = append(records, record)
	}
	return records
}

// ImportBillingHistory imports a single billing history record into the store
func (k Keeper) ImportBillingHistory(ctx sdk.Context, record types.SaveBillingHistory) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(fmt.Sprintf("%s-%s", types.BillingHistoryKey, record.ChainletId)))
	uniqueKey := fmt.Sprintf("%s-%d", record.EpochIdentifier, record.EpochNumber)
	value := k.cdc.MustMarshal(&record)
	store.Set([]byte(uniqueKey), value)
}

// ImportValidatorPayoutHistory imports a single validator payout history record into the store
func (k Keeper) ImportValidatorPayoutHistory(ctx sdk.Context, record types.ValidatorPayoutHistory) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(fmt.Sprintf("%s-%s", types.ValidatorPayoutHistoryKey, record.ValidatorAddress)))
	uniqueKey := fmt.Sprintf("%s-%d", record.EpochIdentifier, record.EpochNumber)
	value := k.cdc.MustMarshal(&record)
	store.Set([]byte(uniqueKey), value)
}
