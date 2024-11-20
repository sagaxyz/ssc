package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/sagaxyz/ssc/testutil/keeper"
	"github.com/sagaxyz/ssc/x/billing/keeper"
	"github.com/sagaxyz/ssc/x/billing/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) { //nolint:unused
	k, ctx := keepertest.BillingKeeper(t)
	return keeper.NewMsgServerImpl(*k), ctx
}
