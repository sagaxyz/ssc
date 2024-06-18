package keeper

import (
	"context"

	"github.com/sagaxyz/ssc/x/chainlet/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k *Keeper) ListChainletStack(goCtx context.Context, req *types.QueryListChainletStackRequest) (*types.QueryListChainletStackResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var chainletStacks []*types.ChainletStack
	var err error

	ctx := sdk.UnwrapSDKContext(goCtx)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ChainletStackKey))
	pageRes, err := query.Paginate(store, req.Pagination, func(key, value []byte) error {
		var chainletStack types.ChainletStack
		if err := k.cdc.Unmarshal(value, &chainletStack); err != nil {
			return err
		}
		chainletStacks = append(chainletStacks, &chainletStack)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListChainletStackResponse{ChainletStacks: chainletStacks, Pagination: pageRes}, nil
}
