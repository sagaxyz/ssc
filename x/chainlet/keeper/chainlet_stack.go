package keeper

import (
	"fmt"
	"strings"

	cosmossdkerrors "cosmossdk.io/errors"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

func (k *Keeper) NewChainletStack(ctx sdk.Context, cs types.ChainletStack) error {
	// Get the store
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletStackKey)

	// Our key is the display name e.g. SagaEVM. Associated with this key can be many versions
	// Versions is a slice of object type TemplateMetadata
	byteKey := []byte(cs.DisplayName)
	if store.Has(byteKey) {
		// cannot add a duplicate chainlet stack so return an error
		return fmt.Errorf("cannot add chainlet stack %v as it already exists", cs.DisplayName)
	}

	// Validate versions first (before any cache updates)
	for _, version := range cs.Versions {
		if !versions.Check(version.Version) {
			return fmt.Errorf("version string '%s' invalid", version.Version)
		}
	}

	// Reload caches to ensure consistency before modification
	if err := k.loadVersions(ctx); err != nil {
		return fmt.Errorf("failed to load versions: %w", err)
	}

	// Write to KV FIRST (before updating caches)
	value := k.cdc.MustMarshal(&cs)
	store.Set(byteKey, value)

	// Update caches AFTER successful KV write
	for _, version := range cs.Versions {
		if version.Enabled {
			if err := k.AddVersion(ctx, cs.DisplayName, version); err != nil {
				// Reload caches to rollback partial state
				_ = k.loadVersions(ctx) // Best effort recovery
				return err
			}
		}
	}

	return nil
}

func (k *Keeper) AddChainletStackVersion(ctx sdk.Context, stackName string, version types.ChainletStackParams) error {
	// Reload caches to ensure consistency before modification
	if err := k.loadVersions(ctx); err != nil {
		return fmt.Errorf("failed to load versions: %w", err)
	}

	// Get the store
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletStackKey)

	stack, err := k.getChainletStack(ctx, stackName)
	if err != nil {
		return fmt.Errorf("cannot get chainlet stack %s: %w", stackName, err)
	}

	// Validate that the incoming fields can be updated
	err = validateUpdate(stack, version)
	if err != nil {
		return fmt.Errorf("cannot update chainlet stack %s: %w", stackName, err)
	}

	// Upsert the version in KV FIRST (before updating caches)
	stack.Versions = append(stack.Versions, version)
	updatedValue := k.cdc.MustMarshal(&stack)
	store.Set([]byte(stackName), updatedValue)

	// Update caches AFTER successful KV write
	if version.Enabled {
		if err = k.AddVersion(ctx, stack.DisplayName, version); err != nil {
			// Reload caches to rollback partial state
			_ = k.loadVersions(ctx) // Best effort recovery
			return err
		}
	}

	return nil
}

func validateUpdate(stack types.ChainletStack, version types.ChainletStackParams) error {
	for _, v := range stack.Versions {
		if v.Image == version.Image || v.Version == version.Version || v.Checksum == version.Checksum {
			return fmt.Errorf("cannot update with duplicate values for image, version, or checksum")
		}
	}

	if !versions.Check(version.Version) {
		return fmt.Errorf("version string '%s' invalid", version.Version)
	}

	return nil
}

func (k *Keeper) chainletStackVersionAvailable(ctx sdk.Context, name, version string) (bool, error) {
	params, err := k.getChainletStackVersion(ctx, name, version)
	if err != nil {
		return false, err
	}
	if !params.Enabled {
		return false, nil
	}

	return true, nil
}

func (k *Keeper) getChainletStack(ctx sdk.Context, name string) (stack types.ChainletStack, err error) {
	byteKey := []byte(name)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletStackKey)
	if !store.Has(byteKey) {
		err = fmt.Errorf("stack %s not found", name)
		return
	}

	data := store.Get(byteKey)
	k.cdc.MustUnmarshal(data, &stack)
	return
}

func (k *Keeper) getChainletStackVersion(ctx sdk.Context, name, version string) (types.ChainletStackParams, error) {
	// Ensure caches are loaded
	if k.stackVersionParams == nil || k.stackVersions == nil {
		if err := k.loadVersions(ctx); err != nil {
			return types.ChainletStackParams{}, fmt.Errorf("cannot load versions: %w", err)
		}
	}

	// The stack existence is checked via the caches below.

	// O(1) presence check + params fetch
	verKey := normalizeVer(version)
	if _, ok := k.stackVersions[name]; !ok {
		return types.ChainletStackParams{}, fmt.Errorf("stack %q not found", name)
	}
	pmap := k.stackVersionParams[name]
	if pmap == nil {
		return types.ChainletStackParams{}, fmt.Errorf("no versions indexed for stack %q", name)
	}
	p, ok := pmap[verKey]
	if !ok {
		return types.ChainletStackParams{}, fmt.Errorf("stack version %q is not found", version)
	}
	return p, nil
}

// UpdateChainletStackFees updates the per-stack fees in the exact order submitted.
func (k *Keeper) updateChainletStackFees(ctx sdk.Context, creator sdk.AccAddress, stackName string, fees []types.ChainletStackFees) error {
	// Load stack

	supported := make(map[string]struct{})
	for _, denom := range k.escrowKeeper.GetSupportedDenoms(ctx) {
		supported[denom] = struct{}{}
	}

	for _, fee := range fees {
		if _, ok := supported[fee.Denom]; !ok {
			return cosmossdkerrors.Wrapf(
				types.ErrInvalidDenom,
				"denom %s not supported for escrow deposits",
				fee.Denom,
			)
		}
	}

	stack, err := k.getChainletStack(ctx, stackName)
	if err != nil {
		return fmt.Errorf("cannot get chainlet stack %s: %w", stackName, err)
	}

	if !k.aclKeeper.Allowed(ctx, creator) {
		return cosmossdkerrors.Wrapf(types.ErrUnauthorized, "address %s is not allowed to update fees", creator.String())
	}
	// Persist exact order and strings (copy to avoid caller mutation)
	stack.Fees = append([]types.ChainletStackFees(nil), fees...)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletStackKey)
	store.Set([]byte(stackName), k.cdc.MustMarshal(&stack))

	// Emit event using original strings in original order
	err = ctx.EventManager().EmitTypedEvent(&types.EventUpdateChainletFees{
		StackName: stackName,
		Fees:      joinFeesOriginal(fees),
		By:        creator.String(),
	})
	if err != nil {
		return fmt.Errorf("failed to emit event: %w", err)
	}

	return nil
}

func joinFeesOriginal(fees []types.ChainletStackFees) string {
	if len(fees) == 0 {
		return ""
	}
	out := make([]string, 0, len(fees))
	for _, f := range fees {
		out = append(out, f.EpochFee) // original strings, original order
	}
	return strings.Join(out, ",")
}
