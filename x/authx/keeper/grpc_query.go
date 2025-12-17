package keeper

import (
	"bytes"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/authx"
)

var _ authx.QueryServer = Keeper{}

// Grants implements the Query/Grants gRPC method.
// It returns grants for a granter-grantee pair. If msg type URL is set, it returns grants only for that msg type.
func (k Keeper) Grants(ctx context.Context, req *authx.QueryGrantsRequest) (*authx.QueryGrantsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	granter, err := k.authKeeper.AddressCodec().StringToBytes(req.Granter)
	if err != nil {
		return nil, err
	}

	grantee, err := k.authKeeper.AddressCodec().StringToBytes(req.Grantee)
	if err != nil {
		return nil, err
	}

	if req.MsgTypeUrl != "" {
		grant, found := k.getGrant(ctx, grantStoreKey(grantee, granter, req.MsgTypeUrl))
		if !found {
			return nil, errors.Wrapf(authx.ErrNoAuthorizationFound, "authorization not found for %s type", req.MsgTypeUrl)
		}

		authorization, err := grant.GetAuthorization()
		if err != nil {
			return nil, err
		}

		authorizationAny, err := codectypes.NewAnyWithValue(authorization)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return &authx.QueryGrantsResponse{
			Grants: []*authx.Grant{{
				Authorization: authorizationAny,
				Expiration:    grant.Expiration,
			}},
		}, nil
	}

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := grantStoreKey(grantee, granter, "")
	grantsStore := prefix.NewStore(store, key)

	authorizations, pageRes, err := query.GenericFilteredPaginate(k.cdc, grantsStore, req.Pagination, func(key []byte, auth *authx.Grant) (*authx.Grant, error) {
		auth1, err := auth.GetAuthorization()
		if err != nil {
			return nil, err
		}

		authorizationAny, err := codectypes.NewAnyWithValue(auth1)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return &authx.Grant{
			Authorization: authorizationAny,
			Expiration:    auth.Expiration,
		}, nil
	}, func() *authx.Grant {
		return &authx.Grant{}
	})
	if err != nil {
		return nil, err
	}

	return &authx.QueryGrantsResponse{
		Grants:     authorizations,
		Pagination: pageRes,
	}, nil
}

// GranterGrants implements the Query/GranterGrants gRPC method.
func (k Keeper) GranterGrants(ctx context.Context, req *authx.QueryGranterGrantsRequest) (*authx.QueryGranterGrantsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	granter, err := k.authKeeper.AddressCodec().StringToBytes(req.Granter)
	if err != nil {
		return nil, err
	}

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	authxStore := prefix.NewStore(store, grantStoreKey(nil, granter, ""))

	grants, pageRes, err := query.GenericFilteredPaginate(k.cdc, authxStore, req.Pagination, func(key []byte, auth *authx.Grant) (*authx.GrantAuthorization, error) {
		auth1, err := auth.GetAuthorization()
		if err != nil {
			return nil, err
		}

		any, err := codectypes.NewAnyWithValue(auth1)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		grantee := firstAddressFromGrantStoreKey(key)
		return &authx.GrantAuthorization{
			Granter:       req.Granter,
			Grantee:       grantee.String(),
			Authorization: any,
			Expiration:    auth.Expiration,
		}, nil
	}, func() *authx.Grant {
		return &authx.Grant{}
	})
	if err != nil {
		return nil, err
	}

	return &authx.QueryGranterGrantsResponse{
		Grants:     grants,
		Pagination: pageRes,
	}, nil
}

// GranteeGrants implements the Query/GranteeGrants gRPC method.
func (k Keeper) GranteeGrants(ctx context.Context, req *authx.QueryGranteeGrantsRequest) (*authx.QueryGranteeGrantsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	grantee, err := k.authKeeper.AddressCodec().StringToBytes(req.Grantee)
	if err != nil {
		return nil, err
	}

	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), GrantKey)

	authorizations, pageRes, err := query.GenericFilteredPaginate(k.cdc, store, req.Pagination, func(key []byte, auth *authx.Grant) (*authx.GrantAuthorization, error) {
		auth1, err := auth.GetAuthorization()
		if err != nil {
			return nil, err
		}

		granter, g, _ := parseGrantStoreKey(append(GrantKey, key...))
		if !bytes.Equal(g, grantee) {
			return nil, nil
		}

		authorizationAny, err := codectypes.NewAnyWithValue(auth1)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		return &authx.GrantAuthorization{
			Authorization: authorizationAny,
			Expiration:    auth.Expiration,
			Granter:       granter.String(),
			Grantee:       req.Grantee,
		}, nil
	}, func() *authx.Grant {
		return &authx.Grant{}
	})
	if err != nil {
		return nil, err
	}

	return &authx.QueryGranteeGrantsResponse{
		Grants:     authorizations,
		Pagination: pageRes,
	}, nil
}
