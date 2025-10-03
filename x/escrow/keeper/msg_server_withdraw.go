package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sagaxyz/ssc/x/escrow/types"
)

func (k msgServer) Withdraw(goCtx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return &types.MsgWithdrawResponse{}, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return &types.MsgWithdrawResponse{}, err
	}

	err = k.WithdrawAll(ctx, addr, msg.ChainId)
	if err != nil {
		return &types.MsgWithdrawResponse{}, err
	}

	return &types.MsgWithdrawResponse{}, nil
}
