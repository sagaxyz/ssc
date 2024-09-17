package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

func (k *Keeper) loadVersions(ctx sdk.Context) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletStackKey)

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	k.stackVersions = make(map[string]*versions.Versions)
	for ; iterator.Valid(); iterator.Next() {
		var stack types.ChainletStack
		k.cdc.MustUnmarshal(iterator.Value(), &stack)

		if k.stackVersions[stack.DisplayName] == nil {
			k.stackVersions[stack.DisplayName] = versions.New()
		}
		for _, version := range stack.Versions {
			if !version.Enabled {
				continue
			}
			err := k.stackVersions[stack.DisplayName].Add(version.Version)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (k *Keeper) AddVersion(ctx sdk.Context, stackName string, version string) error {
	if k.stackVersions == nil {
		err := k.loadVersions(ctx)
		if err != nil {
			return err
		}
	}

	if k.stackVersions[stackName] == nil {
		k.stackVersions[stackName] = versions.New()
	}
	err := k.stackVersions[stackName].Add(version)
	if err != nil {
		return err
	}

	return nil
}

func (k *Keeper) RemoveVersion(ctx sdk.Context, stackName string, version string) error {
	if k.stackVersions == nil {
		return nil
	}
	if k.stackVersions[stackName] == nil {
		return nil
	}

	err := k.stackVersions[stackName].Remove(version)
	if err != nil {
		return err
	}
	if k.stackVersions[stackName].Empty() {
		delete(k.stackVersions, stackName)
	}

	return nil
}

func (k *Keeper) LatestVersion(ctx sdk.Context, stackName string, version string) (latestVersion string, err error) {
	if k.stackVersions == nil {
		k.stackVersions = make(map[string]*versions.Versions)
		err = k.loadVersions(ctx)
		if err != nil {
			return
		}
	}

	if k.stackVersions[stackName] == nil {
		latestVersion = version
		return
	}

	return k.stackVersions[stackName].LatestCompatible(version)
}
