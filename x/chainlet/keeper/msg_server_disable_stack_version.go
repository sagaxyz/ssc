package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dactypes "github.com/sagaxyz/saga-sdk/x/acl/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (k msgServer) DisableChainletStackVersion(goCtx context.Context, msg *types.MsgDisableChainletStackVersion) (resp *types.MsgDisableChainletStackVersionResponse, err error) {
	err = msg.ValidateBasic()
	if err != nil {
		return
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	p := k.GetParams(ctx)
	if p.ChainletStackProtections {
		dacAddr := dactypes.NewAddress(dactypes.AddressFormat_ADDRESS_BECH32, msg.Creator)
		if !k.dacKeeper.Allowed(ctx, dacAddr) {
			err = fmt.Errorf("address %s not allowed to modify chainlet stacks", msg.Creator)
			return
		}
	}

	stack, err := k.getChainletStack(ctx, msg.DisplayName)
	if err != nil {
		err = fmt.Errorf("cannot get chainlet stack %s: %w", msg.DisplayName, err)
		return
	}

	//TODO avoid loop
	var found bool
	for i, version := range stack.Versions {
		if version.Version != msg.Version {
			continue
		}

		if !version.Enabled {
			return // Already disabled
		}
		stack.Versions[i].Enabled = false
		found = true
		break
	}
	if !found {
		err = fmt.Errorf("cannot find chainlet stack %s version %s", msg.DisplayName, msg.Version)
		return
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletStackKey)
	value := k.cdc.MustMarshal(&stack)
	store.Set([]byte(msg.DisplayName), value)

	err = k.RemoveVersion(ctx, msg.DisplayName, msg.Version)
	if err != nil {
		return
	}

	return &types.MsgDisableChainletStackVersionResponse{}, ctx.EventManager().EmitTypedEvent(&types.EventChainletStackVersionDisabled{
		Name:    msg.DisplayName,
		Version: msg.Version,
	})
}
