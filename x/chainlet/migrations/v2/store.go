package v2

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	chainletStore := prefix.NewStore(ctx.KVStore(storeKey), []byte(types.ChainletKey))
	stackStore := prefix.NewStore(ctx.KVStore(storeKey), []byte(types.ChainletStackKey))

	// Remove stack versions with an invalid version string
	iterator := stackStore.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		var stack types.ChainletStack
		cdc.MustUnmarshal(iterator.Value(), &stack)

		newVersions := make([]types.ChainletStackParams, 0, len(stack.Versions))
		for _, version := range stack.Versions {
			ok := versions.Check(version.Version)
			if ok {
				newVersions = append(newVersions, version)
				continue
			}

			ctx.Logger().Info(fmt.Sprintf("deleting invalid version %s in stack %s", version.Version, stack.DisplayName))
		}
		stack.Versions = newVersions

		if len(stack.Versions) == 0 {
			// Completely delete stack with no versions
			ctx.Logger().Info(fmt.Sprintf("deleting stack %s with no versions", stack.DisplayName))

			key := iterator.Key()
			defer stackStore.Delete(key)
		} else {
			updated := cdc.MustMarshal(&stack)
			key := iterator.Key()
			defer stackStore.Set(key, updated)
		}
	}
	iterator.Close()

	// Remove chainlets with now non-existent chainlet stack
	iterator = chainletStore.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		var chainlet types.Chainlet
		cdc.MustUnmarshal(iterator.Value(), &chainlet)

		ok := versions.Check(chainlet.ChainletStackVersion)
		if ok {
			continue
		}

		ctx.Logger().Info(fmt.Sprintf("deleting chainlet %s with deleted stack %s:%s", chainlet.ChainId, chainlet.ChainletStackName, chainlet.ChainletStackVersion))
		key := iterator.Key()
		defer chainletStore.Delete(key)
	}
	iterator.Close()

	return nil
}
