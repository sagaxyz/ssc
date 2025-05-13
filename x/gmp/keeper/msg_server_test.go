package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/sagaxyz/ssc/testutil/keeper"
	"github.com/sagaxyz/ssc/x/gmp/keeper"
	"github.com/sagaxyz/ssc/x/gmp/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.GmpKeeper(t)
	return keeper.NewMsgServerImpl(*k), ctx
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}
