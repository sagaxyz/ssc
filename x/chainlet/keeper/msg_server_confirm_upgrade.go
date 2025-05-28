package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (k msgServer) ConfirmUpgradeChainlet(goCtx context.Context, msg *types.MsgConfirmUpgradeChainlet) (*types.MsgConfirmUpgradeChainletResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := msg.ValidateBasic()
	if err != nil {
		return &types.MsgConfirmUpgradeChainletResponse{}, err
	}

	chainlet, err := k.Chainlet(ctx, msg.ChainId)
	if err != nil {
		return &types.MsgConfirmUpgradeChainletResponse{}, err
	}

	//TODO verify source of the message
	//TODO verify upgrade plan matches

	err = k.finishUpgrading(ctx, &chainlet)
	if err != nil {
		return &types.MsgConfirmUpgradeChainletResponse{}, err
	}

	return &types.MsgConfirmUpgradeChainletResponse{}, ctx.EventManager().EmitTypedEvent(&types.EventUpdateChainlet{
		ChainId:      msg.ChainId,
		//TODO
	})
}
