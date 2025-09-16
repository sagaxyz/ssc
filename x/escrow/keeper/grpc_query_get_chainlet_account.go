package keeper

// import (
// 	"context"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/sagaxyz/ssc/x/escrow/types"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// func (k Keeper) GetChainletAccount(goCtx context.Context, req *types.QueryGetChainletAccountRequest) (*types.QueryGetChainletAccountResponse, error) {
// 	if req == nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid request")
// 	}

// 	ctx := sdk.UnwrapSDKContext(goCtx)

// 	// TODO: Process the query

// 	acc, err := k.GetKprChainletAccount(ctx, req.ChainId)
// 	if err != nil {
// 		return &types.QueryGetChainletAccountResponse{}, err
// 	}
// 	return &types.QueryGetChainletAccountResponse{
// 		Account: acc,
// 	}, nil
// }
