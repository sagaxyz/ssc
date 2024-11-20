package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/sagaxyz/ssc/testutil/keeper"
	"github.com/sagaxyz/ssc/x/escrow/keeper"
	"github.com/sagaxyz/ssc/x/escrow/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) { //nolint:unused
	k, ctx := keepertest.EscrowKeeper(t)
	return keeper.NewMsgServerImpl(*k), ctx
}
