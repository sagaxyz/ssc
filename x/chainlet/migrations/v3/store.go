package v2

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/chainlet/exported"
	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, legacySubspace exported.Subspace) error {
	chainletStore := prefix.NewStore(ctx.KVStore(storeKey), []byte(types.ChainletKey))

	// Add default values of new params
	legacySubspace.Set(ctx, []byte("AutomaticChainletUpgrades"), true)
	legacySubspace.Set(ctx, []byte("AutomaticChainletUpgradeInterval"), int64(100))
	var p types.Params
	legacySubspace.GetParamSet(ctx, &p)
	err := p.Validate()
	if err != nil {
		return err
	}
	legacySubspace.SetParamSet(ctx, &p)

	// Set all existing chainlets to have auto-upgrade enabled
	iterator := chainletStore.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		var chainlet types.Chainlet
		cdc.MustUnmarshal(iterator.Value(), &chainlet)

		chainlet.AutoUpgradeStack = true

		updated := cdc.MustMarshal(&chainlet)
		key := iterator.Key()
		defer chainletStore.Set(key, updated)

		ctx.Logger().Info(fmt.Sprintf("enabled auto-upgrade for chainlet %s", chainlet.ChainId))
	}
	iterator.Close()

	// Set all existing stack versions as enabled
	stackStore := prefix.NewStore(ctx.KVStore(storeKey), []byte(types.ChainletStackKey))
	iterator = stackStore.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		var stack types.ChainletStack
		cdc.MustUnmarshal(iterator.Value(), &stack)

		for _, version := range stack.Versions {
			version.Enabled = true
		}

		updated := cdc.MustMarshal(&stack)
		key := iterator.Key()
		defer stackStore.Set(key, updated)

		ctx.Logger().Info(fmt.Sprintf("enabled all chainlet stack versions of %s", stack.DisplayName))
	}
	iterator.Close()

	return nil
}
