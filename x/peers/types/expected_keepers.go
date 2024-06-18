package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	chainlettypes "github.com/sagaxyz/ssc/x/chainlet/types"
)

type ChainletKeeper interface {
	Chainlet(ctx sdk.Context, chainId string) (chainlet chainlettypes.Chainlet, err error)
}
