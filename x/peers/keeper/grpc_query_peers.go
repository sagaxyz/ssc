package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sagaxyz/ssc/x/peers/types"
)

func (k Keeper) Peers(goCtx context.Context, req *types.QueryPeersRequest) (resp *types.QueryPeersResponse, err error) {
	if req == nil {
		err = status.Error(codes.InvalidArgument, "invalid request")
		return
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Requesting a non-existing chain ID is different than having an empty peer list
	_, err = k.chainletKeeper.Chainlet(ctx, req.ChainId)
	if err != nil {
		err = errors.New("no such chain ID")
		return
	}

	resp = &types.QueryPeersResponse{
		Peers: k.GetPeers(ctx, req.ChainId),
	}
	return
}
