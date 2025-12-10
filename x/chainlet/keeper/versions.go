package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

func (k *Keeper) loadVersions(ctx sdk.Context) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletStackKey)

	it := store.Iterator(nil, nil)
	defer it.Close()

	k.stackVersions = make(map[string]*versions.Versions)
	k.stackVersionParams = make(map[string]map[string]types.ChainletStackParams)

	for ; it.Valid(); it.Next() {
		var stack types.ChainletStack
		k.cdc.MustUnmarshal(it.Value(), &stack)

		// init caches for this stack
		if k.stackVersions[stack.DisplayName] == nil {
			k.stackVersions[stack.DisplayName] = versions.New()
		}
		if k.stackVersionParams[stack.DisplayName] == nil {
			k.stackVersionParams[stack.DisplayName] = make(map[string]types.ChainletStackParams)
		}

		// single pass to fill both indexes
		for _, sv := range stack.Versions {
			if !sv.Enabled {
				continue
			}
			// presence cache
			if err := k.stackVersions[stack.DisplayName].Add(sv.Version); err != nil {
				return err
			}
			// params cache
			k.stackVersionParams[stack.DisplayName][normalizeVer(sv.Version)] = sv
		}
	}
	return nil
}

func normalizeVer(v string) string {
	if len(v) > 0 && (v[0] == 'v' || v[0] == 'V') {
		return v[1:]
	}
	return v
}

// VersionExistsInCache checks if a version already exists in the cache.
// Loads the cache if it's not initialized to ensure accurate results.
// Note: This function may be called independently (not just before AddVersion),
// so it must load the cache if needed for correctness.
func (k *Keeper) VersionExistsInCache(ctx sdk.Context, stackName, version string) bool {
	// Ensure caches are loaded for accurate results
	if k.stackVersionParams == nil || k.stackVersions == nil {
		if err := k.loadVersions(ctx); err != nil {
			// If loading fails, assume version doesn't exist (safe default)
			return false
		}
	}
	verKey := normalizeVer(version)
	pmap := k.stackVersionParams[stackName]
	if pmap == nil {
		return false
	}
	_, exists := pmap[verKey]
	return exists
}

func (k *Keeper) AddVersion(ctx sdk.Context, stackName string, params types.ChainletStackParams) error {
	version := params.Version
	if k.stackVersions == nil || k.stackVersionParams == nil {
		if err := k.loadVersions(ctx); err != nil {
			return err
		}
	}
	if k.stackVersions[stackName] == nil {
		k.stackVersions[stackName] = versions.New()
	}
	if k.stackVersionParams[stackName] == nil {
		k.stackVersionParams[stackName] = make(map[string]types.ChainletStackParams)
	}
	if err := k.stackVersions[stackName].Add(version); err != nil {
		return err
	}
	k.stackVersionParams[stackName][normalizeVer(version)] = params
	return nil
}

func (k *Keeper) RemoveVersion(ctx sdk.Context, stackName, version string) error {
	if k.stackVersions == nil || k.stackVersionParams == nil {
		return nil
	}
	verKey := normalizeVer(version)
	if s := k.stackVersions[stackName]; s != nil {
		if err := s.Remove(version); err != nil {
			return err
		}
		if s.Empty() {
			delete(k.stackVersions, stackName)
		}
	}
	if m := k.stackVersionParams[stackName]; m != nil {
		delete(m, verKey)
		if len(m) == 0 {
			delete(k.stackVersionParams, stackName)
		}
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

// Testing helpers
func (k *Keeper) Versions(stack string) []string {
	s := k.stackVersions[stack]
	if s == nil {
		return nil
	}

	return s.Export()
}
func (k *Keeper) DeleteVersions() {
	k.stackVersions = nil
}
