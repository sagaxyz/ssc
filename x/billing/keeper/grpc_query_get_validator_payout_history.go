package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sagaxyz/ssc/x/billing/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetValidatorPayoutHistory(goCtx context.Context, req *types.QueryGetValidatorPayoutHistoryRequest) (*types.QueryGetValidatorPayoutHistoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(fmt.Sprintf("%s-%s", types.ValidatorPayoutHistoryKey, req.ValidatorAddress)))

	var vph []*types.ValidatorPayoutHistory
	pageRes, err := query.Paginate(store, req.Pagination, func(key, value []byte) error {
		var vphr types.ValidatorPayoutHistory
		k.cdc.MustUnmarshal(value, &vphr)
		vph = append(vph, &vphr)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetValidatorPayoutHistoryResponse{Validatorpayouthistory: vph, Pagination: pageRes}, nil
}
