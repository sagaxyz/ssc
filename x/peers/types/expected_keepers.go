package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	chainlettypes "github.com/sagaxyz/ssc/x/chainlet/types"
)

type ChainletKeeper interface {
	Chainlet(ctx sdk.Context, chainId string) (chainlet chainlettypes.Chainlet, err error)
}

type StakingHooks interface {
	AfterValidatorRemoved(ctx context.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error // Must be called when a validator is deleted
}
