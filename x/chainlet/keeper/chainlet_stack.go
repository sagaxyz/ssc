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

	for _, version := range cs.Versions {
		if !versions.Check(version.Version) {
			return fmt.Errorf("version string '%s' invalid", version.Version)
		}
		if version.Enabled {
			err := k.AddVersion(ctx, cs.DisplayName, version.Version)
			if err != nil {
				return err
			}
		}
	}

	value := k.cdc.MustMarshal(&cs)
	store.Set(byteKey, value)

	return nil
}

func (k *Keeper) AddChainletStackVersion(ctx sdk.Context, stackName string, version types.ChainletStackParams) error {
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

	// Upsert the version
	stack.Versions = append(stack.Versions, version)
	updatedValue := k.cdc.MustMarshal(&stack)
	store.Set([]byte(stackName), updatedValue)

	// Store in the version tree for automatic updates
	if version.Enabled {
		err = k.AddVersion(ctx, stack.DisplayName, version.Version)
		if err != nil {
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

func (k *Keeper) getChainletStackVersion(ctx sdk.Context, name, version string) (params types.ChainletStackParams, err error) {
	stack, err := k.getChainletStack(ctx, name)
	if err != nil {
		err = fmt.Errorf("cannot get chainlet stack with name %s: %w", name, err)
		return
	}

	//TODO avoid loop
	for _, v := range stack.Versions {
		if v.Version == version {
			params = v
			return
		}
	}

	err = fmt.Errorf("stack version %s is not found", version)
	return
}

// UpdateChainletStackFees updates the per-stack fees in the exact order submitted.
func (k *Keeper) updateChainletStackFees(ctx sdk.Context, creator sdk.AccAddress, stackName string, fees []types.ChainletStackFees) error {
	// Load stack
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
	ctx.EventManager().EmitTypedEvent(&types.EventUpdateChainletFees{
		StackName: stackName,
		Fees:      joinFeesOriginal(fees),
		By:        creator.String(),
	})

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
