package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/sagaxyz/ssc/testutil/keeper"

	"github.com/sagaxyz/ssc/x/peers/keeper"
	"github.com/sagaxyz/ssc/x/peers/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//nolint:unused
func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.PeersKeeper(t)
	return keeper.NewMsgServerImpl(k), sdk.WrapSDKContext(ctx)
}
