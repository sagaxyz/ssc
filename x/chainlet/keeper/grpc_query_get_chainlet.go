package keeper

import (
	"context"

	"github.com/sagaxyz/ssc/x/chainlet/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k *Keeper) GetChainlet(goCtx context.Context, req *types.QueryGetChainletRequest) (*types.QueryGetChainletResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	chainlet, err := k.Chainlet(ctx, req.ChainId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &types.QueryGetChainletResponse{
		Chainlet: chainlet,
	}, nil
}
