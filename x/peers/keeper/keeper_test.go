package keeper_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/sagaxyz/ssc/x/peers/keeper"
	testutil "github.com/sagaxyz/ssc/x/peers/testutil"
	"github.com/sagaxyz/ssc/x/peers/types"
)

var (
	chainIDs = []string{"chain_1-1", "chain_2-1", "chain_3-1"}
	addrs    = map[string][]string{
		"chain_1-1": {"aa@123.123.123.123:1234", "bb@111.111.111.111:1234"},
		"chain_2-1": {"cc@100.100.100.100:1234", "dd@example.com:1234"},
		"chain_3-1": {"ee@google.com:1234"},
	}
	accounts = []sdk.AccAddress{
		sdk.AccAddress("test1"),
		sdk.AccAddress("test2"),
	}
)

type TestSuite struct {
	suite.Suite

	ctx            sdk.Context
	chainletKeeper *testutil.MockChainletKeeper
	peersKeeper    keeper.Keeper
	queryClient    types.QueryClient
	msgServer      types.MsgServer
}

func (s *TestSuite) SetupTest() {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	paramsKey := storetypes.NewKVStoreKey(paramstypes.StoreKey)
	paramsTKey := storetypes.NewTransientStoreKey(paramstypes.TStoreKey)
	ctx := sdktestutil.DefaultContextWithKeys(
		map[string]*storetypes.KVStoreKey{
			types.StoreKey:       storeKey,
			paramstypes.StoreKey: paramsKey,
		},
		map[string]*storetypes.TransientStoreKey{
			paramstypes.TStoreKey: paramsTKey,
		},
		nil,
	)
	s.ctx = ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()

	// gomock initializations
	ctrl := gomock.NewController(s.T())
	s.chainletKeeper = testutil.NewMockChainletKeeper(ctrl)

	//nolint:staticcheck
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

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
