package keeper

import (
	"fmt"
	"time"

	query "github.com/cosmos/cosmos-sdk/types/query"

	cosmossdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sagaxyz/ssc/x/billing/types"
	chainlettypes "github.com/sagaxyz/ssc/x/chainlet/types"
	epochstypes "github.com/sagaxyz/ssc/x/epochs/types"
)

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	stacks, err := k.chainletkeeper.ListChainletStack(ctx, &chainlettypes.QueryListChainletStackRequest{})

	if err != nil {
		ctx.Logger().Error("could not list chainlet stacks. Error: " + err.Error())
		return cosmossdkerrors.Wrapf(types.ErrInternalFailure, "could not list chainlet stacks. Error: %s", err.Error())
	}
	kvs := make(map[string]*chainlettypes.ChainletStack)
	for _, stack := range stacks.ChainletStacks {
		kvs[stack.DisplayName] = stack
	}
	resp, err := k.chainletkeeper.ListChainlets(ctx, &chainlettypes.QueryListChainletsRequest{
		Pagination: &query.PageRequest{Limit: k.chainletkeeper.GetParams(ctx).MaxChainlets},
	})
	if err != nil {
		ctx.Logger().Error("could not list chainlets. Error: " + err.Error())
		return cosmossdkerrors.Wrapf(types.ErrInternalFailure, "could not list chainlets. Error: %s", err.Error())
	}

	epochInfo := k.epochskeeper.GetEpochInfo(ctx, epochIdentifier)
	ctx.Logger().Debug("Current epoch start time is " + epochInfo.CurrentEpochStartTime.Format(time.RFC3339))
	ctx.Logger().Info("attempting billing of chainlets for epoch " + fmt.Sprintf("%d", epochNumber) + " with epoch identifier " + epochIdentifier)

	var skipped, billed, failed []string

	if len(resp.Chainlets) == 0 {
		ctx.Logger().Info("no active chainlets to be billed this epoch: " + fmt.Sprintf("%d", epochNumber))
		return nil
	}

	for _, ch := range resp.Chainlets {
		// Skip service or offline
		if ch.IsServiceChainlet {
			ctx.Logger().Debug("skipping billing for service chainlet: " + ch.ChainId)
			skipped = append(skipped, ch.ChainId)
			continue
		}
		if ch.Status == chainlettypes.Status_STATUS_OFFLINE {
			ctx.Logger().Debug("skipping billing for inactive chainlet: " + ch.ChainId)
			skipped = append(skipped, ch.ChainId)
			continue
		}

		// Only bill chainlets that appear in kvs (as per your comment)
		stack, ok := kvs[ch.ChainletStackName]
		if !ok {
			ctx.Logger().Debug("skipping billing for chainlet (no stack in kvs): " + ch.ChainId)
			skipped = append(skipped, ch.ChainId)
			continue
		}

		// Try multiple fee options until one works
		succeeded := false
		var errs []string

		for i, fee := range stack.Fees {
			epochFee, perr := sdk.ParseCoinNormalized(fee.EpochFee)
			if perr != nil {
				msg := fmt.Sprintf("fee[%d] parse failed: %q err=%v", i, fee.EpochFee, perr)
				ctx.Logger().Error("billing parse error for " + ch.ChainId + ": " + msg)
				errs = append(errs, msg)
				continue // try next fee option
			}

			// Attempt billing with this coin option
			berr := k.BillAccount(ctx, epochFee, *ch, "epoch-start-billing")
			if berr != nil {
				msg := fmt.Sprintf("fee[%d] %s billing failed: %v", i, epochFee.String(), berr)
				ctx.Logger().Error("billing error for " + ch.ChainId + ": " + msg)
				errs = append(errs, msg)
				continue // try next fee option
			}

			// Success on this fee option; stop trying others
			succeeded = true
			ctx.Logger().Info(fmt.Sprintf("billed %s successfully with %s", ch.ChainId, epochFee.String()))
			break
		}

		if succeeded {
			billed = append(billed, ch.ChainId)
			continue
		}

		// All fee options failed -> stop chainlet and record failure
		stopErr := k.chainletkeeper.StopChainlet(ctx, ch.ChainId)
		if stopErr != nil {
			ctx.Logger().Error("could not stop chainlet " + ch.ChainId + ". Error: " + stopErr.Error())
			errs = append(errs, "stop failed: "+stopErr.Error())
		}
		ctx.Logger().Error(fmt.Sprintf("all billing options failed for %s; reasons: %v", ch.ChainId, errs))
		failed = append(failed, ch.ChainId)
	}

	ctx.Logger().Info("skipped billing for chainlets: " + fmt.Sprintf("%v", skipped))
	ctx.Logger().Info("failed billing for chainlets: " + fmt.Sprintf("%v", failed))
	ctx.Logger().Info("successful billing for chainlets: " + fmt.Sprintf("%v", billed))
	return nil
}

// AfterEpochEnd is the epoch end hook.
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {

	validatorPayoutEpoch := k.GetParams(ctx).ValidatorPayoutEpoch

	if epochIdentifier != validatorPayoutEpoch {
		ctx.Logger().Info("skipping distribution of epoch fees to validators as epoch length is " + epochIdentifier + " and we only process epoch fees at epoch length " + validatorPayoutEpoch)
		return nil
	}

	validators, err := k.stakingkeeper.GetValidators(ctx, 100)
	if err != nil {
		return err
	}
	numValidators := len(validators)                                  // number of validators
	moduleAccount := k.accountkeeper.GetModuleAccount(ctx, "billing") // module account address for the billing module
	moduleAccountBalance := k.bankkeeper.GetAllBalances(ctx, moduleAccount.GetAddress())
	validatorDepositAmount := moduleAccountBalance.QuoInt(math.NewIntFromUint64(uint64(numValidators)))
	ctx.Logger().Debug("Validator deposit amount is " + validatorDepositAmount.String())

	ctx.Logger().Debug("Module account address is " + moduleAccount.GetAddress().String() + " and module account balance is " + moduleAccountBalance.String())

	epochInfo := k.epochskeeper.GetEpochInfo(ctx, epochIdentifier)
	epochEventStartTime := epochInfo.CurrentEpochStartTime.Format(time.RFC3339)

	for _, v := range validators {
		var valAddr sdk.ValAddress
		var err error

		ctx.Logger().Debug("Validator being processed is " + v.OperatorAddress)

		if validatorDepositAmount.IsValid() && validatorDepositAmount.IsAllPositive() {
			valAddr, err = sdk.ValAddressFromBech32(v.OperatorAddress)
			if err != nil {
				ctx.Logger().Error("could not get the validator address from operator address: " + v.OperatorAddress + ". Error: " + err.Error())
				continue
			}

			ctx.Logger().Debug("Validator hex address is: " + valAddr.String())
		} else {
			ctx.Logger().Error("funds in billing module with address " + moduleAccount.GetAddress().String() + " could not be validated, or no funds exist, for distribution to validators: " + v.OperatorAddress)
			continue
		}

		err = k.PayEpochFeeToValidator(ctx, validatorDepositAmount, "billing", sdk.AccAddress(valAddr), "epoch fee reward")
		if err != nil {
			ctx.Logger().Error("could not pay epoch fee to validator " + v.OperatorAddress + ". Error: " + err.Error())
			continue
		}
		err = k.SaveValidatorPayoutHistory(ctx, types.ValidatorPayoutHistory{
			ValidatorAddress: sdk.AccAddress(valAddr).String(),
			EpochIdentifier:  epochIdentifier,
			EpochNumber:      int32(epochNumber),
			EpochStartTime:   epochEventStartTime,
			RewardAmount:     validatorDepositAmount.String(),
		})
		if err != nil {
			ctx.Logger().Error("could not save validator payout history for validator " + sdk.AccAddress(valAddr).String() + ". Error: " + err.Error())
		}
	}

	return nil
}

// ___________________________________________________________________________________________________

// Hooks is the wrapper struct for the incentives keeper.
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

// Hooks returns the hook wrapper struct.
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// BeforeEpochStart is the epoch start hook.
func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

// AfterEpochEnd is the epoch end hook.
func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
