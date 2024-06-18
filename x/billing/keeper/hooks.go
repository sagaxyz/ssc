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

// BeforeEpochStart is the epoch start hook.
func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	ctxx := sdk.WrapSDKContext(ctx)

	stacks, err := k.chainletkeeper.ListChainletStack(ctxx, &chainlettypes.QueryListChainletStackRequest{})
	if err != nil {
		ctx.Logger().Error("could not list chainlet stacks. Error: " + err.Error())
		return cosmossdkerrors.Wrapf(types.ErrInternalFailure, "could not list chainlet stacks. Error: "+err.Error())
	}

	kvs := make(map[string]*chainlettypes.ChainletStack)
	for _, stack := range stacks.ChainletStacks {
		if stack.Fees.EpochLength == epochIdentifier {
			ctx.Logger().Info("billing for chainlets in chainlet stack: " + stack.DisplayName + " as its epoch length is " + epochIdentifier)
			kvs[stack.DisplayName] = stack
		} else {
			ctx.Logger().Debug("skipping billing for chainlets in chainlet stack: " + stack.DisplayName + " as its epoch length is " +
				stack.Fees.EpochLength + " and we only bill for epoch length " + epochIdentifier)
		}
	}
	// do this more efficiently later
	chainlets, err := k.chainletkeeper.ListChainlets(ctxx, &chainlettypes.QueryListChainletsRequest{Pagination: &query.PageRequest{Limit: 5000}})
	if err != nil {
		ctx.Logger().Error("could not list chainlets. Error: " + err.Error())
		return cosmossdkerrors.Wrapf(types.ErrInternalFailure, "could not list chainlets. Error: "+err.Error())
	}

	epochInfo := k.epochskeeper.GetEpochInfo(ctx, epochIdentifier)
	epochEventStartTime := epochInfo.CurrentEpochStartTime.Format(time.RFC3339)
	ctx.Logger().Debug("Current epoch start time is " + epochEventStartTime)

	ctx.Logger().Info("attempting billing of chainlets for epoch " + fmt.Sprintf("%d", epochNumber) + " with epoch identifier " + epochIdentifier)

	skipped := []string{}
	billed := []string{}
	failed := []string{}

	if chainlets.Chainlets == nil {
		ctx.Logger().Info("no active chainlets to be billed this epoch: " + fmt.Sprintf("%d", epochNumber))
		return nil
	}

	for _, chainlet := range chainlets.Chainlets {

		// Skip inactive (OFFLINE) chainlets
		if chainlet.Status == chainlettypes.Status_STATUS_OFFLINE {
			ctx.Logger().Debug("skipping billing for inactive chainlet: " + chainlet.ChainId)
			skipped = append(skipped, chainlet.ChainId)
			continue
		}

		// here we limit the billing for chainlets we find in our kvs map. If it does not exist there, we do not bill it
		stack, ok := kvs[chainlet.ChainletStackName]
		if !ok {
			ctx.Logger().Debug("skipping billing for chainlet: " + chainlet.ChainId)
			skipped = append(skipped, chainlet.ChainId)
			continue
		}

		epochfee, err := sdk.ParseCoinNormalized(stack.Fees.EpochFee)

		if err != nil {
			ctx.Logger().Error("could not calculate epoch fee for chainlet " + chainlet.ChainId + ". Error: " + err.Error())
			err = k.chainletkeeper.StopChainlet(ctx, chainlet.ChainId)
			if err != nil {
				ctx.Logger().Error("could not stop chainlet " + chainlet.ChainId + ". Error: " + err.Error())
			}
			failed = append(failed, chainlet.ChainId)
			continue
		}

		err = k.BillAccount(ctx, epochfee, *chainlet, epochIdentifier, "epoch-start-billing")
		if err != nil {
			ctx.Logger().Error("could not bill account " + chainlet.ChainId + ". Error: " + err.Error())
			err = k.chainletkeeper.StopChainlet(ctx, chainlet.ChainId)
			if err != nil {
				ctx.Logger().Error("could not stop chainlet " + chainlet.ChainId + ". Error: " + err.Error())
			}
			failed = append(failed, chainlet.ChainId)
			continue
		}

		billed = append(billed, chainlet.ChainId)
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

	validators := k.stakingkeeper.GetValidators(ctx, 100)
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
