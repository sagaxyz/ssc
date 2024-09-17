package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sagaxyz/ssc/x/billing/types"
	chainlettypes "github.com/sagaxyz/ssc/x/chainlet/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetBillingHistory(goCtx context.Context, req *types.QueryGetBillingHistoryRequest) (*types.QueryGetBillingHistoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(fmt.Sprintf("%s-%s", types.BillingHistoryKey, req.ChainId)))

	var bh []*types.BillingHistory
	// get chainlet info to fill in the details for returning to the user
	chainletReq := &chainlettypes.QueryGetChainletRequest{ChainId: req.ChainId}

	chainletRes, err := k.chainletkeeper.GetChainlet(ctx, chainletReq)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve chainlet info for chain %s. Error: %v", req.ChainId, err)
	}

	pageRes, err := query.Paginate(store, req.Pagination, func(key, value []byte) error {
		var sbhr types.SaveBillingHistory
		k.cdc.MustUnmarshal(value, &sbhr)
		// get epoch info
		epochInfo := k.epochskeeper.GetEpochInfo(ctx, sbhr.EpochIdentifier)
		epochSince := (epochInfo.CurrentEpoch - int64(sbhr.EpochNumber))
		epochEventStartTime := epochInfo.CurrentEpochStartTime.Add(-time.Duration(epochSince * int64(epochInfo.Duration)))
		bhr := types.BillingHistory{
			ChainletId:        sbhr.ChainletId,
			ChainletName:      chainletRes.Chainlet.ChainletName,
			ChainletOwner:     chainletRes.Chainlet.Launcher,
			ChainletStackName: chainletRes.Chainlet.ChainletStackName,
			EpochIdentifier:   sbhr.EpochIdentifier,
			EpochNumber:       sbhr.EpochNumber,
			EpochStartTime:    epochEventStartTime.Format(time.RFC3339),
			BilledAmount:      sbhr.BilledAmount,
		}
		bh = append(bh, &bhr)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetBillingHistoryResponse{Billhistory: bh, Pagination: pageRes}, nil
}
