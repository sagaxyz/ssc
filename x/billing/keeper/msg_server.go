package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sagaxyz/ssc/x/billing/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) SetPlatformValidators(goCtx context.Context, msg *types.MsgSetPlatformValidators) (*types.MsgSetPlatformValidatorsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Creator != m.Keeper.GetAuthority() {
		return nil, types.ErrUnauthorized
	}

	err := m.Keeper.SetPlatformValidators(ctx, msg.PlatformValidators)
	if err != nil {
		return nil, err
	}

	return &types.MsgSetPlatformValidatorsResponse{}, nil
}
