package keeper_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	testutil "github.com/sagaxyz/ssc/x/peers/testutil"
	"github.com/sagaxyz/ssc/x/peers/types"
	"github.com/sagaxyz/ssc/x/peers/keeper"
)

var (
	chainIDs = []string{"chain_1-1", "chain_2-1", "chain_3-1"}
)

type KeeperTestSuite struct {
	suite.Suite

	ctx            sdk.Context
	chainletKeeper *testutil.MockChainletKeeper
	peersKeeper    keeper.Keeper
	queryClient    types.QueryClient
	msgServer      types.MsgServer
}

func (s *KeeperTestSuite) SetupTest() {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeTKey := storetypes.NewTransientStoreKey("transient_test")
	testCtx := sdktestutil.DefaultContextWithDB(s.T(), storeKey, storeTKey)
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()

	// gomock initializations
	ctrl := gomock.NewController(s.T())
	s.chainletKeeper = testutil.NewMockChainletKeeper(ctrl)

	//nolint:staticcheck
	paramsKey := storetypes.NewKVStoreKey(paramstypes.StoreKey)
	paramsTKey := storetypes.NewTransientStoreKey(paramstypes.TStoreKey)
	paramsKeeper := paramskeeper.NewKeeper(encCfg.Codec, encCfg.Amino, paramsKey, paramsTKey) //nolint:staticcheck
	paramsKeeper.Subspace(paramstypes.ModuleName)
	paramsKeeper.Subspace(types.ModuleName)
	sub, _ := paramsKeeper.GetSubspace(types.ModuleName)

	s.ctx = ctx
	s.peersKeeper = keeper.New(
		encCfg.Codec,
		storeKey,
		sub,
		s.chainletKeeper,
	)

	types.RegisterInterfaces(encCfg.InterfaceRegistry)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, s.peersKeeper)

	s.queryClient = types.NewQueryClient(queryHelper)
	s.msgServer = keeper.NewMsgServerImpl(s.peersKeeper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
