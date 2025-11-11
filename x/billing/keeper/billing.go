package keeper

import (
	"fmt"
	"time"

	cosmossdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sagaxyz/ssc/x/billing/types"
	chainlettypes "github.com/sagaxyz/ssc/x/chainlet/types"
)

func (k Keeper) BillAccount(ctx sdk.Context, amount sdk.Coin, chainlet chainlettypes.Chainlet, memo string) error {
	err := k.escrowkeeper.BillAccount(ctx, amount, chainlet.ChainId, "billing")
	if err != nil {
		ctx.Logger().Info(fmt.Sprintf("failed to bill account %s for %s at epoch %s", chainlet.ChainId, amount.String(), memo))
		//nolint:errcheck // Event emission errors are non-critical
		ctx.EventManager().EmitTypedEvent(&types.BillingEvent{
			ChainId: chainlet.ChainId,
			Amount:  amount.String(),
			Memo:    memo,
			Success: false,
			Debit:   true,
		})
		return err
	}
	ctx.Logger().Info(fmt.Sprintf("successfully billed account %s for %s at epoch %s", chainlet.ChainId, amount.String(), memo))
	//nolint:errcheck // Event emission errors are non-critical
	ctx.EventManager().EmitTypedEvent(&types.BillingEvent{
		ChainId: chainlet.ChainId,
		Amount:  amount.String(),
		Memo:    memo,
		Success: true,
		Debit:   true,
	})
	// Save billing history
	epochIdentifier := k.GetParams(ctx).BillingEpoch
	epochInfo := k.epochskeeper.GetEpochInfo(ctx, epochIdentifier)
	epochEventStartTime := epochInfo.CurrentEpochStartTime.Format(time.RFC3339)

	err = k.SaveBillingHistory(ctx, types.BillingHistory{
		ChainletOwner:     chainlet.Launcher,
		ChainletId:        chainlet.ChainId,
		ChainletName:      chainlet.ChainletName,
		ChainletStackName: chainlet.ChainletStackName,
		EpochIdentifier:   epochIdentifier,
		EpochNumber:       epochInfo.CurrentEpoch,
		EpochStartTime:    epochEventStartTime,
		BilledAmount:      amount.String(),
	})
	if err != nil {
		ctx.Logger().Error("could not save billing history for chainlet " + chainlet.ChainletName + ". Error: " + err.Error())
	}

	return nil
}

func (k Keeper) PayEpochFeeToValidator(ctx sdk.Context, epochFee sdk.Coins, fromModuleName string, valAddr sdk.AccAddress, memo string) (err error) {
	moduleAccount := k.accountkeeper.GetModuleAccount(ctx, fromModuleName) // module account address for the billing module

	err = k.bankkeeper.SendCoinsFromModuleToAccount(ctx, "billing", valAddr, epochFee)
	if err != nil {
		ctx.Logger().Error("could not send funds from billing account address " + moduleAccount.String() + " to validator" + valAddr.String() + ". Error: " + err.Error())

		return cosmossdkerrors.Wrapf(types.ErrInternalBillingFailure, "could not send funds from billing account address %s to validator %s. Error: %s", moduleAccount.String(), valAddr.String(), err.Error())
	} else {
		ctx.Logger().Debug("sent epoch fee of " + epochFee.String() + " to validator address " + valAddr.String())
	}

	newModuleBalance := k.bankkeeper.GetAllBalances(ctx, moduleAccount.GetAddress())
	ctx.Logger().Debug("Updated module balance of account " + moduleAccount.GetAddress().String() + " is " + newModuleBalance.String())

	return nil
}

func (k Keeper) SaveBillingHistory(ctx sdk.Context, billinghistory types.BillingHistory) error {
	// Get the store
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(fmt.Sprintf("%s-%s", types.BillingHistoryKey, billinghistory.ChainletId)))

	uniqueBillingRecord := fmt.Sprintf("%s-%d", billinghistory.EpochIdentifier, billinghistory.EpochNumber)
	byteKey := []byte(uniqueBillingRecord)
	saveBillingHistory := types.SaveBillingHistory{
		ChainletId:      billinghistory.ChainletId,
		EpochIdentifier: billinghistory.EpochIdentifier,
		EpochNumber:     billinghistory.EpochNumber,
		BilledAmount:    billinghistory.BilledAmount,
	}
	value := k.cdc.MustMarshal(&saveBillingHistory)
	if len(value) == 0 {
		return cosmossdkerrors.Wrap(types.ErrInternalFailure, "could not marshal billing history input for appending to the kv store")
	}
	if store.Has(byteKey) {
		// cannot add a duplicate billing record so return an error
		return cosmossdkerrors.Wrap(types.ErrInternalFailure, fmt.Sprintf("cannot add billing record %v as it already exists", uniqueBillingRecord))
	} else {
		store.Set(byteKey, value)
	}

	return nil
}

func (k Keeper) GetChainletBillingHistory(ctx sdk.Context, chainId string) ([]*types.BillingHistory, error) {
	var bh []*types.BillingHistory

	// Get the store
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(fmt.Sprintf("%s-%s", types.BillingHistoryKey, chainId)))
	it := store.Iterator(nil, nil)
	if !it.Valid() {
		return nil, cosmossdkerrors.Wrapf(types.ErrNoRecords, "no billing history found for chain %s", chainId)
	}

	// get chainlet info to fill in the details for returning to the user
	chainletReq := &chainlettypes.QueryGetChainletRequest{ChainId: chainId}

	chainletRes, err := k.chainletkeeper.GetChainlet(ctx, chainletReq)
	if err != nil {
		return nil, cosmossdkerrors.Wrapf(types.ErrInternalFailure, "could not retrieve chainlet info for chain %s. Error: %v", chainId, err)
	}

	for val := it.Value(); it.Valid(); it.Next() {
		var sbhr types.SaveBillingHistory
		k.cdc.MustUnmarshal(val, &sbhr)
		// get epoch info
		epochInfo := k.epochskeeper.GetEpochInfo(ctx, sbhr.EpochIdentifier)
		epochSince := (epochInfo.CurrentEpoch - sbhr.EpochNumber)
		epochEventStartTime := epochInfo.CurrentEpochStartTime.Add(-time.Duration(epochSince * int64(epochInfo.Duration)))
		bhr := types.BillingHistory{
			ChainletId:        sbhr.ChainletId,
			ChainletName:      chainletRes.Chainlet.ChainletName,
			ChainletOwner:     chainletRes.Chainlet.Launcher,
			ChainletStackName: chainletRes.Chainlet.ChainletStackName,
			EpochIdentifier:   sbhr.EpochIdentifier,
			EpochNumber:       sbhr.EpochNumber,
			EpochStartTime:    epochEventStartTime.Format(time.RFC3339),
			BilledAmount:      sbhr.BilledAmount,
		}
		bh = append(bh, &bhr)
	}
	return bh, nil
}

func (k Keeper) SaveValidatorPayoutHistory(ctx sdk.Context, payouthistory types.ValidatorPayoutHistory) error {
	// Get the store
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(fmt.Sprintf("%s-%s", types.ValidatorPayoutHistoryKey, payouthistory.ValidatorAddress)))

	uniqueBillingRecord := fmt.Sprintf("%s-%d", payouthistory.EpochIdentifier, payouthistory.EpochNumber)
	byteKey := []byte(uniqueBillingRecord)
	value := k.cdc.MustMarshal(&payouthistory)
	if len(value) == 0 {
		return cosmossdkerrors.Wrapf(types.ErrJSONMarhsal, "could not marshal validator payout history input for appending to the kv store")
	}
	if store.Has(byteKey) {
		// cannot add a duplicate billing record so return an error
		return cosmossdkerrors.Wrapf(types.ErrDuplicateRecord, "cannot add validator payout record %v as it already exists", uniqueBillingRecord)
	} else {
		store.Set(byteKey, value)
	}

	return nil
}

func (k Keeper) GetKprValidatorPayoutHistory(ctx sdk.Context, validatorAddress string) ([]*types.ValidatorPayoutHistory, error) {
	var vph []*types.ValidatorPayoutHistory

	// Get the store
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(fmt.Sprintf("%s-%s", types.ValidatorPayoutHistoryKey, validatorAddress)))
	it := store.Iterator(nil, nil)
	if !it.Valid() {
		return nil, cosmossdkerrors.Wrapf(types.ErrNoRecords, "no validator payout history found for validator %s", validatorAddress)
	}

	for val := it.Value(); it.Valid(); it.Next() {
		var vphr types.ValidatorPayoutHistory
		k.cdc.MustUnmarshal(val, &vphr)
		vph = append(vph, &vphr)

	}
	return vph, nil
}

func (k Keeper) BillAndRestartChainlet(ctx sdk.Context, chainId string) error {
	started, err := k.chainletkeeper.IsChainletStarted(ctx, chainId)
	if err != nil {
		return err
	}
	if started {
		return nil
	}

	stack, err := k.chainletkeeper.GetChainletStackInfo(ctx, chainId)
	if err != nil {
		return err
	}

	billed := false
	for _, feeOption := range stack.Fees {
		epochfee, err := sdk.ParseCoinNormalized(feeOption.EpochFee)
		if err != nil {
			return err
		}

		chainlet, err := k.chainletkeeper.GetChainletInfo(ctx, chainId)
		if err != nil {
			return err
		}

		// Check if there is enough funds to restart the chainlet
		err = k.BillAccount(ctx, epochfee, *chainlet, "restarting chainlet")
		if err == nil {
			billed = true
			break
		}
	}

	if !billed {
		return cosmossdkerrors.Wrapf(types.ErrInternalBillingFailure, "could not bill account for chainlet %s", chainId)
	}

	err = k.chainletkeeper.StartExistingChainlet(ctx, chainId)
	if err != nil {
		return err
	}
	return nil
}
