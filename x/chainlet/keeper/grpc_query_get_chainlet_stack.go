package keeper

import (
	"context"

	"github.com/sagaxyz/ssc/x/chainlet/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k *Keeper) GetChainletStack(goCtx context.Context, req *types.QueryGetChainletStackRequest) (*types.QueryGetChainletStackResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if req.DisplayName == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	stack, err := k.getChainletStack(ctx, req.DisplayName)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "TODO") //TODO
	}
	return &types.QueryGetChainletStackResponse{
		ChainletStack: stack,
	}, nil
}
