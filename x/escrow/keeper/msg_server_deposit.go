package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sagaxyz/ssc/x/escrow/types"
)

func (k msgServer) Deposit(goCtx context.Context, msg *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return &types.MsgDepositResponse{}, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return &types.MsgDepositResponse{}, err
	}

	coin, err := sdk.ParseCoinNormalized(msg.Amount)
	if err != nil {
		return &types.MsgDepositResponse{}, err
	}

	err = k.deposit(ctx, addr, msg.ChainId, coin)
	if err != nil {
		return &types.MsgDepositResponse{}, err
	}

	return &types.MsgDepositResponse{}, nil
}
