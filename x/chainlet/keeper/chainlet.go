package keeper

import (
	"encoding/binary"
	"fmt"
	"reflect"

	cosmossdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (k *Keeper) Chainlet(ctx sdk.Context, chainId string) (chainlet types.Chainlet, err error) {
	byteKey := []byte(chainId)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletKey)
	if !store.Has(byteKey) {
		err = fmt.Errorf("key %s not found", chainId)
		return
	}

	chainletBytes := store.Get(byteKey)
	if len(chainletBytes) == 0 {
		panic(fmt.Sprintf("no data at chainlet %s", chainId))
	}
	k.cdc.MustUnmarshal(chainletBytes, &chainlet)

	return
}

func (k *Keeper) NewChainlet(ctx sdk.Context, chainlet types.Chainlet) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletKey)

	key := []byte(chainlet.ChainId)
	if store.Has(key) {
		return cosmossdkerrors.Wrapf(types.ErrChainletExists, "chainlet with chainId %s already exists", chainlet.ChainId)
	}

	avail, err := k.chainletStackVersionAvailable(ctx, chainlet.ChainletStackName, chainlet.ChainletStackVersion)
	if err != nil {
		return cosmossdkerrors.Wrapf(types.ErrInvalidChainletStack, "cannot use stack %s version %s: %s", chainlet.ChainletStackName, chainlet.ChainletStackVersion, err)
	}
	if !avail {
		return cosmossdkerrors.Wrapf(types.ErrInvalidChainletStack, "stack %s version %s not available", chainlet.ChainletStackName, chainlet.ChainletStackVersion)
	}

	value := k.cdc.MustMarshal(&chainlet)
	store.Set(key, value)
	k.incrementChainletCount(ctx)
	return nil
}

func (k *Keeper) UpgradeChainletStackVersion(ctx sdk.Context, chainId, stackVersion string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletKey)

	key := []byte(chainId)
	if !store.Has(key) {
		return cosmossdkerrors.Wrapf(types.ErrInvalidChainId, "chainlet with chainId %s not found", chainId)
	}

	chainlet, err := k.Chainlet(ctx, chainId)
	if err != nil {
		return err
	}

	avail, err := k.chainletStackVersionAvailable(ctx, chainlet.ChainletStackName, stackVersion)
	if err != nil {
		return cosmossdkerrors.Wrapf(types.ErrInvalidChainletStack, "cannot upgrade to stack %s version %s: %s", chainlet.ChainletStackName, stackVersion, err)
	}
	if !avail {
		return cosmossdkerrors.Wrapf(types.ErrInvalidChainletStack, "stack %s version %s not available", chainlet.ChainletStackName, chainlet.ChainletStackVersion)
	}

	chainlet.ChainletStackVersion = stackVersion

	updatedValue := k.cdc.MustMarshal(&chainlet)
	store.Set(key, updatedValue)

	return nil
}

func updateChainletParams(curParams *types.ChainletParams, params *types.ChainletParams) error { //nolint:unused
	curElem := reflect.ValueOf(curParams).Elem()
	newElem := reflect.ValueOf(params).Elem()
	for i := 0; i < curElem.NumField(); i++ {
		newValue := newElem.Field(i)
		// reflection messed up the structs fields order, fast bail out
		if curElem.Type().Field(i).Name != newElem.Type().Field(i).Name {
			return fmt.Errorf("cannot update chainlet parameters")
		}
		// ensure to prevent unwanted wipes of current state params
		if !valueIsNullOrBlank(newValue) {
			curElem.Field(i).Set(newValue)
		}
	}
	return nil
}

func valueIsNullOrBlank(val reflect.Value) bool { //nolint:unused
	return val.Interface() == nil || val.Interface() == "" || val.IsZero()
}

func (k *Keeper) IsChainletStarted(ctx sdk.Context, chainId string) (bool, error) {
	c, err := k.GetChainletInfo(ctx, chainId)
	if err != nil {
		return false, err
	}

	if c.Status == types.Status_STATUS_ONLINE {
		return true, nil
	}
	return false, nil
}

func (k *Keeper) StartExistingChainlet(ctx sdk.Context, chainId string) error {
	c, err := k.GetChainletInfo(ctx, chainId)
	if err != nil {
		return fmt.Errorf("cannot start existing chainlet %s: %v", chainId, err)
	}

	c.Status = types.Status_STATUS_ONLINE
	k.setChainletInfo(ctx, c)

	return nil
}

func (k *Keeper) GetChainletStackInfo(ctx sdk.Context, chainId string) (*types.ChainletStack, error) {
	c, err := k.GetChainletInfo(ctx, chainId)
	if err != nil {
		return nil, err
	}

	// Get the chainlet stack store
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletStackKey)
	byteKey := []byte(c.ChainletStackName)

	if !store.Has(byteKey) {
		return nil, cosmossdkerrors.Wrapf(types.ErrInvalidChainletStack, "chainlet stack with name %s not found", c.ChainletStackName)
	}
	var stack types.ChainletStack
	storeChainletStackData := store.Get(byteKey)
	k.cdc.MustUnmarshal(storeChainletStackData, &stack)

	return &stack, nil
}

func (k *Keeper) StopChainlet(ctx sdk.Context, chainId string) error {
	c, err := k.GetChainletInfo(ctx, chainId)
	if err != nil {
		return fmt.Errorf("cannot stop chainlet %s: %v", chainId, err)
	}
	c.Status = types.Status_STATUS_OFFLINE
	k.setChainletInfo(ctx, c)
	ctx.Logger().Info(fmt.Sprintf("Successfully stopped chainlet %s", chainId))
	return nil
}

func (k *Keeper) ChainletExists(ctx sdk.Context, chainId string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ChainletKey))
	key := []byte(chainId)
	return store.Has(key)
}

func (k *Keeper) GetChainletInfo(ctx sdk.Context, chainId string) (*types.Chainlet, error) {
	// Get the store
	lcStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletKey)
	byteLCKey := []byte(chainId)

	if !lcStore.Has(byteLCKey) {
		return nil, fmt.Errorf("cannot get info for chainlet %s", chainId)
	}

	var c types.Chainlet

	storeData := lcStore.Get(byteLCKey)
	k.cdc.MustUnmarshal(storeData, &c)
	return &c, nil
}

func (k *Keeper) setChainletInfo(ctx sdk.Context, chainlet *types.Chainlet) {
	lcStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletKey)
	byteLCKey := []byte(chainlet.ChainId)
	updatedValue := k.cdc.MustMarshal(chainlet)
	lcStore.Set(byteLCKey, updatedValue)
}

func (k Keeper) InitializeChainletCount(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(0))
	store.Set(types.ChainletCountKey, bz)
}

func (k Keeper) GetChainletCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ChainletCountKey)
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetChainletCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.ChainletCountKey, bz)
}

func (k Keeper) incrementChainletCount(ctx sdk.Context) {
	count := k.GetChainletCount(ctx)
	k.SetChainletCount(ctx, count+1)
}

func (k *Keeper) AutoUpgradeChainlets(ctx sdk.Context) error {
	iter := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletKey).Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var chainlet types.Chainlet
		k.cdc.MustUnmarshal(iter.Value(), &chainlet)

		if !chainlet.AutoUpgradeStack {
			ctx.Logger().Debug(fmt.Sprintf("skipping auto-upgrade for chainlet %s\n", chainlet.ChainId))
			continue
		}

		latestVersion, err := k.LatestVersion(ctx, chainlet.ChainletStackName, chainlet.ChainletStackVersion)
		if err != nil {
			iter.Close()
			return err
		}
		if latestVersion == chainlet.ChainletStackVersion {
			iter.Close()
			return nil
		}

		available, err := k.chainletStackVersionAvailable(ctx, chainlet.ChainletStackName, latestVersion)
		if err != nil || !available {
			iter.Close()
			//TODO change to panic in the future, should never happen if the loaded versions are consistent with the state
			return fmt.Errorf("chainlet stack %s has unavailable version %s loaded", chainlet.ChainletStackName, latestVersion)
		}

		if chainlet.ChainletStackVersion == latestVersion {
			ctx.Logger().Debug(fmt.Sprintf("chainlet %s: %s is at its latest available version\n", chainlet.ChainId, chainlet.ChainletStackVersion))
			continue
		}

		ctx.Logger().Info(fmt.Sprintf("upgrading chainlet %s: %s to %s\n", chainlet.ChainId, chainlet.ChainletStackVersion, latestVersion))
		chainlet.ChainletStackVersion = latestVersion
		defer k.setChainletInfo(ctx, &chainlet)
	}
	iter.Close()

	return nil
}
