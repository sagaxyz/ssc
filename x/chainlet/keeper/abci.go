package keeper

import (
	"fmt"

	//errorsmod "cosmossdk.io/errors"
	//"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	//"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (k *Keeper) BeginBlock(ctx sdk.Context) error {
	k.ForcePendingVSC(ctx)

	p := k.GetParams(ctx)
	if p.AutomaticChainletUpgrades && ctx.BlockHeight()%p.AutomaticChainletUpgradeInterval == 0 {
		ctx.Logger().Debug("checking chainlets for available upgrades")

		err := k.AutoUpgradeChainlets(ctx)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("failed to auto-upgrade all chainlets: %s", err))
			// We can ignore the error because chainlets not being upgraded does not prevent
			// the network from continuing and the function cannot result in an invalid state.
		}
	}

	//TODO move to a function
	/*iterator := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletUpgradesKey).Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		chainID := iterator.Key()

		chainlet, err := k.Chainlet(ctx, string(chainID))
		if err != nil {
			return err
		}

		if chainlet.Upgrade == nil {
			panic("no upgrade info")
		}

		clientState, ex := k.ibcKeeper.ClientKeeper.GetClientState(ctx, chainlet.IbcClientId)
		if !ex {
			return errorsmod.Wrapf("client state missing for client ID '%s'", chainlet.IbcClientId)
		}

		height := clientState.GetLatestHeight()
		if height > chainlet.Upgrade.Height {
			// Chain failed to stop before the upgrade height => cancel the upgrade
			//TODO remove the upgrade with an error
		}
		if height == chainlet.Upgrade.Height {
			// Upgrade height reached
			//TODO remove the upgrade flag
			//TODO udate the chainlet struct
		}
	}*/

	return nil
}
