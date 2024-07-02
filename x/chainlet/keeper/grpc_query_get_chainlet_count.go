package keeper

import (
	"context"
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sagaxyz/ssc/x/chainlet/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetChainletCount(goCtx context.Context, req *types.QueryGetChainletCountRequest) (*types.QueryGetChainletCountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NumChainletsKey)
	ctx.Logger().Info("GetChainletCount", "count", binary.BigEndian.Uint64(bz))
	return &types.QueryGetChainletCountResponse{Count: binary.BigEndian.Uint64(bz)}, nil
}
