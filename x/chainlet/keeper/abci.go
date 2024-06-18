package keeper

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	p := k.GetParams(ctx)
	if p.AutomaticChainletUpgrades && req.Header.Height%p.AutomaticChainletUpgradeInterval == 0 {
		ctx.Logger().Debug("checking chainlets for available upgrades")

		err := k.AutoUpgradeChainlets(ctx)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("failed to auto-upgrade all chainlets: %s", err))
			// We can ignore the error because chainlets not being upgraded does not prevent
			// the network from continuing and the function cannot result in an invalid state.
		}
	}

	k.ForcePendingVSC(ctx)
}
