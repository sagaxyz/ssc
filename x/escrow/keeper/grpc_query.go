package keeper

import (
	"context"
	"strings"

	cosmossdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sagaxyz/ssc/x/escrow/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Ensure Keeper implements the generated gRPC interface.
var _ types.QueryServer = Keeper{}

// ---------------------------------------
// Params
// ---------------------------------------

func (k Keeper) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(sdkCtx)
	return &types.QueryParamsResponse{Params: params}, nil
}

// ---------------------------------------
// Chainlet head
// ---------------------------------------

func (k Keeper) GetChainletAccount(ctx context.Context, req *types.QueryGetChainletAccountRequest) (*types.QueryGetChainletAccountResponse, error) {
	if req == nil || req.ChainId == "" {
		return nil, status.Error(codes.InvalidArgument, "chainId required")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	acc, ok := k.getChainlet(sdkCtx, req.ChainId)
	if !ok {
		return nil, status.Error(codes.NotFound, cosmossdkerrors.Wrapf(types.ErrChainletAccountNotFound, "chainlet %s not found", req.ChainId).Error())
	}

	return &types.QueryGetChainletAccountResponse{Account: acc}, nil
}

// ---------------------------------------
// GetPools (per chainlet)
// ---------------------------------------

func (k Keeper) GetPools(ctx context.Context, req *types.QueryPoolsRequest) (*types.QueryPoolsResponse, error) {
	if req == nil || req.ChainId == "" {
		return nil, status.Error(codes.InvalidArgument, "chainId required")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := sdkCtx.KVStore(k.storeKey)
	pfx := prefix.NewStore(store, types.PoolPrefix(req.ChainId))

	var out []*types.DenomPool
	pageRes, err := query.Paginate(pfx, req.Pagination, func(key, value []byte) error {
		var p types.DenomPool
		k.cdc.MustUnmarshal(value, &p)
		out = append(out, &p)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolsResponse{
		Pools:      out,
		Pagination: pageRes,
	}, nil
}

// ---------------------------------------
// GetFunders (per {chainlet, denom})
// ---------------------------------------

func (k Keeper) GetFunders(ctx context.Context, req *types.QueryFundersRequest) (*types.QueryFundersResponse, error) {
	if req == nil || req.ChainId == "" || req.Denom == "" {
		return nil, status.Error(codes.InvalidArgument, "chainId and denom required")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := sdkCtx.KVStore(k.storeKey)
	pfx := prefix.NewStore(store, types.FunderPrefix(req.ChainId, req.Denom))

	var out []*types.FunderEntry
	pageRes, err := query.Paginate(pfx, req.Pagination, func(key, value []byte) error {
		// Under this prefix, key == "{addr}"
		addr := string(key)
		var f types.Funder
		k.cdc.MustUnmarshal(value, &f)
		out = append(out, &types.FunderEntry{
			Address: addr,
			Funder:  f,
		})
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryFundersResponse{
		Funders:    out,
		Pagination: pageRes,
	}, nil
}

// ---------------------------------------
// Single funder shares
// ---------------------------------------

func (k Keeper) GetFunder(ctx context.Context, req *types.QueryFunderSharesRequest) (*types.QueryFunderSharesResponse, error) {
	if req == nil || req.ChainId == "" || req.Denom == "" || req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "chainId, denom, and address required")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Validate bech32; ignore returned address (we keep string keys).
	if _, err := sdk.AccAddressFromBech32(req.Address); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid bech32 address")
	}

	f, ok := k.getFunder(sdkCtx, req.ChainId, req.Denom, req.Address)
	if !ok {
		return nil, status.Error(codes.NotFound, cosmossdkerrors.Wrap(types.ErrFunderNotFound, req.Address).Error())
	}

	return &types.QueryFunderSharesResponse{Shares: &f}, nil
}

// ---------------------------------------
// Positions (reverse index: by funder)
// ---------------------------------------

func (k Keeper) GetFunderBalance(ctx context.Context, req *types.QueryFunderPositionsRequest) (*types.QueryFunderPositionsResponse, error) {
	if req == nil || req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address required")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Validate bech32; ignore returned address (we keep string keys).
	if _, err := sdk.AccAddressFromBech32(req.Address); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid bech32 address")
	}

	store := sdkCtx.KVStore(k.storeKey)
	pfx := prefix.NewStore(store, types.ByFunderPrefix(req.Address))

	var positions []*types.Position
	pageRes, err := query.Paginate(pfx, req.Pagination, func(key, _ []byte) error {
		// key under this prefix is "{chainId}/{denom}"
		rel := string(key)
		parts := strings.SplitN(rel, "/", 2)
		if len(parts) != 2 {
			// skip malformed keys
			return nil
		}
		chainID, denom := parts[0], parts[1]

		// get funder row
		f, ok := k.getFunder(sdkCtx, chainID, denom, req.Address)
		if !ok {
			return nil
		}

		positions = append(positions, &types.Position{
			ChainId: chainID,
			Denom:   denom,
			Shares:  f,
		})
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryFunderPositionsResponse{
		Positions:  positions,
		Pagination: pageRes,
	}, nil
}
