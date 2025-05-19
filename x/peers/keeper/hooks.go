package keeper

import (
	"context"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// AfterValidatorRemoved implements types.StakingHooks.
func (h Hooks) AfterValidatorRemoved(ctx context.Context, consAddr types.ConsAddress, valAddr types.ValAddress) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	h.k.DeleteAllValidatorData(sdkCtx, valAddr.String())
	return nil
}

// AfterDelegationModified implements types.StakingHooks.
func (h Hooks) AfterDelegationModified(ctx context.Context, delAddr types.AccAddress, valAddr types.ValAddress) error {
	return nil
}

// AfterUnbondingInitiated implements types.StakingHooks.
func (h Hooks) AfterUnbondingInitiated(ctx context.Context, id uint64) error {
	return nil
}

// AfterValidatorBeginUnbonding implements types.StakingHooks.
func (h Hooks) AfterValidatorBeginUnbonding(ctx context.Context, consAddr types.ConsAddress, valAddr types.ValAddress) error {
	return nil
}

// AfterValidatorBonded implements types.StakingHooks.
func (h Hooks) AfterValidatorBonded(ctx context.Context, consAddr types.ConsAddress, valAddr types.ValAddress) error {
	return nil
}

// AfterValidatorCreated implements types.StakingHooks.
func (h Hooks) AfterValidatorCreated(ctx context.Context, valAddr types.ValAddress) error {
	return nil
}

// BeforeDelegationCreated implements types.StakingHooks.
func (h Hooks) BeforeDelegationCreated(ctx context.Context, delAddr types.AccAddress, valAddr types.ValAddress) error {
	return nil
}

// BeforeDelegationRemoved implements types.StakingHooks.
func (h Hooks) BeforeDelegationRemoved(ctx context.Context, delAddr types.AccAddress, valAddr types.ValAddress) error {
	return nil
}

// BeforeDelegationSharesModified implements types.StakingHooks.
func (h Hooks) BeforeDelegationSharesModified(ctx context.Context, delAddr types.AccAddress, valAddr types.ValAddress) error {
	return nil
}

// BeforeValidatorModified implements types.StakingHooks.
func (h Hooks) BeforeValidatorModified(ctx context.Context, valAddr types.ValAddress) error {
	return nil
}

// BeforeValidatorSlashed implements types.StakingHooks.
func (h Hooks) BeforeValidatorSlashed(ctx context.Context, valAddr types.ValAddress, fraction math.LegacyDec) error {
	return nil
}
