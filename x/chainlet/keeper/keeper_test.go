package keeper_test

import (
	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/sagaxyz/ssc/x/chainlet"
	"github.com/sagaxyz/ssc/x/chainlet/keeper"
	chainlettestutil "github.com/sagaxyz/ssc/x/chainlet/testutil"
	"github.com/sagaxyz/ssc/x/chainlet/types"
)

var (
	fees = types.ChainletStackFees{
		EpochFee:    "10utsaga",
		EpochLength: "minute",
		SetupFee:    "10utsaga",
	}
	addrs = []sdk.AccAddress{
		sdk.AccAddress("test1"),
		sdk.AccAddress("test2"),
	}
	creator = addrs[0]
)

type TestSuite struct {
	suite.Suite

	chainletKeeper *keeper.Keeper
	ctx            sdk.Context
	msgServer      types.MsgServer
	stakingKeeper  *chainlettestutil.MockStakingKeeper
	providerKeeper *chainlettestutil.MockProviderKeeper
	aclKeeper      *chainlettestutil.MockAclKeeper
	escrowKeeper   *chainlettestutil.MockEscrowKeeper
	billingKeeper  *chainlettestutil.MockBillingKeeper
}

/*
	func TestKeeperTestSuite(t *testing.T) {
		suite.Run(t, new(TestSuite))
	}
*/
func (s *TestSuite) SetupTest() {
	encCfg := moduletestutil.MakeTestEncodingConfig(chainlet.AppModuleBasic{})
	key := storetypes.NewKVStoreKey(types.StoreKey)
	paramsKey := storetypes.NewKVStoreKey(paramstypes.StoreKey)
	paramsTKey := storetypes.NewTransientStoreKey(paramstypes.TStoreKey)
	ctx := testutil.DefaultContextWithKeys(
		map[string]*storetypes.KVStoreKey{
			types.StoreKey:       key,
			paramstypes.StoreKey: paramsKey,
		},
		map[string]*storetypes.TransientStoreKey{
			paramstypes.TStoreKey: paramsTKey,
		},
		nil,
	)
	s.ctx = ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})

	ctrl := gomock.NewController(s.T())
	s.stakingKeeper = chainlettestutil.NewMockStakingKeeper(ctrl)
	s.providerKeeper = chainlettestutil.NewMockProviderKeeper(ctrl)
	s.aclKeeper = chainlettestutil.NewMockAclKeeper(ctrl)
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

	paramsKeeper := paramskeeper.NewKeeper(encCfg.Codec, encCfg.Amino, paramsKey, paramsTKey)
	paramsKeeper.Subspace(paramstypes.ModuleName)
	paramsKeeper.Subspace(types.ModuleName)
	sub, _ := paramsKeeper.GetSubspace(types.ModuleName)

	s.chainletKeeper = keeper.NewKeeper(
		encCfg.Codec, key, sub,
		s.stakingKeeper,
		s.providerKeeper,
		s.billingKeeper,
		s.escrowKeeper,
		s.aclKeeper,
	)
	s.msgServer = keeper.NewMsgServerImpl(s.chainletKeeper)

	s.Require().Equal(s.ctx.Logger().With("module", "x/"+types.ModuleName),
		s.chainletKeeper.Logger(s.ctx))

	s.chainletKeeper.SetParams(s.ctx, types.DefaultParams())
	s.chainletKeeper.InitializeChainletCount(s.ctx)
}
