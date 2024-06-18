package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dactypes "github.com/sagaxyz/saga-sdk/x/acl/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (k msgServer) UpdateChainletStack(goCtx context.Context, msg *types.MsgUpdateChainletStack) (*types.MsgUpdateChainletStackResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return &types.MsgUpdateChainletStackResponse{}, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	p := k.GetParams(ctx)
	if p.ChainletStackProtections {
		dacAddr := dactypes.NewAddress(dactypes.AddressFormat_ADDRESS_BECH32, msg.Creator)
		if !k.dacKeeper.Allowed(ctx, dacAddr) {
			return nil, fmt.Errorf("address %s not allowed to create chainlet stacks", msg.Creator)
		}
	}

	version := types.ChainletStackParams{
		Image:    msg.Image,
		Version:  msg.Version,
		Checksum: msg.Checksum,
		Enabled:  true,
	}
	err = k.AddChainletStackVersion(ctx, msg.DisplayName, version)
	if err != nil {
		return nil, fmt.Errorf("error while adding chainlet stack version: %w", err)
	}

	return &types.MsgUpdateChainletStackResponse{}, ctx.EventManager().EmitTypedEvent(&types.EventNewChainletStackVersion{
		Name:    msg.DisplayName,
		Version: msg.Version,
	})
}
