package keeper

import (
	"context"

	"github.com/sagaxyz/ssc/x/chainlet/types"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k *Keeper) ListChainlets(goCtx context.Context, req *types.QueryListChainletsRequest) (*types.QueryListChainletsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var chainlets []*types.Chainlet
	var err error

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletKey)
	pageRes, err := query.Paginate(store, req.Pagination, func(key, value []byte) error {
		var chainlet types.Chainlet
		if err := k.cdc.Unmarshal(value, &chainlet); err != nil {
			return err
		}
		chainlets = append(chainlets, &chainlet)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListChainletsResponse{Chainlets: chainlets, Pagination: pageRes}, nil
}
