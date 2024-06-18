package keeper_test

import (
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/sagaxyz/ssc/x/chainlet"
	"github.com/sagaxyz/ssc/x/chainlet/keeper"
	chainlettestutil "github.com/sagaxyz/ssc/x/chainlet/testutil"
	"github.com/sagaxyz/ssc/x/chainlet/types"
)

var (
	fees = types.ChainletStackFees{
		EpochFee:    "10upsaga",
		EpochLength: "minute",
		SetupFee:    "10upsaga",
	}
	addrs = []sdk.AccAddress{
		sdk.AccAddress("test1"),
		sdk.AccAddress("test2"),
	}
	creator = addrs[0]
)

func DefaultContextWithDB(t *testing.T, keys []storetypes.StoreKey, tkeys []storetypes.StoreKey) testutil.TestContext { //TODO remove
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	for _, key := range keys {
		cms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	}
	for _, tkey := range tkeys {
		cms.MountStoreWithDB(tkey, storetypes.StoreTypeTransient, db)
	}
	err := cms.LoadLatestVersion()
	assert.NoError(t, err)

	ctx := sdk.NewContext(cms, tmproto.Header{}, false, log.NewNopLogger())

	return testutil.TestContext{
		Ctx: ctx,
		DB:  db,
		CMS: cms,
	}
}

type TestSuite struct {
	suite.Suite

	chainletKeeper *keeper.Keeper
	ctx            sdk.Context
	msgServer      types.MsgServer

	providerKeeper *chainlettestutil.MockProviderKeeper
	aclKeeper      *chainlettestutil.MockDacKeeper
	escrowKeeper   *chainlettestutil.MockEscrowKeeper
	billingKeeper  *chainlettestutil.MockBillingKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupTest() {
	encCfg := moduletestutil.MakeTestEncodingConfig(chainlet.AppModuleBasic{})
	key := sdk.NewKVStoreKey(types.StoreKey)
	tkey := sdk.NewTransientStoreKey(types.MemStoreKey)
	paramsKey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	paramsTkey := sdk.NewTransientStoreKey("params_tkey")
	testCtx := DefaultContextWithDB(s.T(),
		[]storetypes.StoreKey{paramsKey, key},
		[]storetypes.StoreKey{paramsTkey, tkey},
	)
	s.ctx = testCtx.Ctx

	ctrl := gomock.NewController(s.T())
	s.providerKeeper = chainlettestutil.NewMockProviderKeeper(ctrl)
	s.aclKeeper = chainlettestutil.NewMockDacKeeper(ctrl)
	s.billingKeeper = chainlettestutil.NewMockBillingKeeper(ctrl)
	s.escrowKeeper = chainlettestutil.NewMockEscrowKeeper(ctrl)

	s.providerKeeper.EXPECT().
		HandleConsumerAdditionProposal(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	s.providerKeeper.EXPECT().
		GetValidatorSetUpdateId(gomock.Any()).
		Return(uint64(0)).
		AnyTimes()
	s.providerKeeper.EXPECT().
		AppendPendingVSCPackets(gomock.Any(), gomock.Any(), gomock.Any()).
		Return().
		AnyTimes()
	s.providerKeeper.EXPECT().
		IncrementValidatorSetUpdateId(gomock.Any()).
		Return().
		AnyTimes()
	s.providerKeeper.EXPECT().
		GetConsumerClientId(gomock.Any(), gomock.Any()).
		Return("abcd", true).
		AnyTimes()

	paramsKeeper := paramskeeper.NewKeeper(encCfg.Codec, encCfg.Amino, paramsKey, paramsTkey)
	paramsKeeper.Subspace(paramstypes.ModuleName)
	paramsKeeper.Subspace(types.ModuleName)
	sub, _ := paramsKeeper.GetSubspace(types.ModuleName)

	s.chainletKeeper = keeper.NewKeeper(
		encCfg.Codec, key, tkey, sub,
		s.providerKeeper,
		s.billingKeeper,
		s.escrowKeeper,
		s.aclKeeper,
	)
	s.msgServer = keeper.NewMsgServerImpl(s.chainletKeeper)

	s.Require().Equal(testCtx.Ctx.Logger().With("module", "x/"+types.ModuleName),
		s.chainletKeeper.Logger(testCtx.Ctx))

	s.chainletKeeper.SetParams(s.ctx, types.DefaultParams())
}
