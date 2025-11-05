package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sagaxyz/ssc/x/escrow/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return &types.MsgUpdateParamsResponse{}, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	addr, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return &types.MsgUpdateParamsResponse{}, err
	}

	if !k.aclKeeper.Allowed(ctx, addr) {
		return &types.MsgUpdateParamsResponse{}, types.ErrUnauthorized
	}

	k.SetParams(ctx, *msg.Params)
	if err != nil {
		return &types.MsgUpdateParamsResponse{}, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
