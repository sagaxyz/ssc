package keeper_test

import (
	ccvprovidertypes "github.com/cosmos/interchain-security/v7/x/ccv/provider/types"
	"github.com/golang/mock/gomock"

	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (s *TestSuite) TestConsumerVSC() {
	s.SetupTest()

	const chainID = "test_12345-1"
	const consumerID = "0"

	// Enable CCV consumer logic for the test
	params := s.chainletKeeper.GetParams(s.ctx)
	params.EnableCCV = true
	s.chainletKeeper.SetParams(s.ctx, params)

	// Set up all mock expectations first
	s.escrowKeeper.EXPECT().
		NewChainletAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	s.billingKeeper.EXPECT().
		BillAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	s.aclKeeper.EXPECT().
		IsAdmin(gomock.Any(), gomock.Any()).
		Return(false).
		AnyTimes()

	// Set up expectations in order (only one round)
	gomock.InOrder(
		// CreateConsumer call
		s.providerMsgServer.EXPECT().
			CreateConsumer(gomock.Any(), gomock.Any()).
			Return(&ccvprovidertypes.MsgCreateConsumerResponse{
				ConsumerId: consumerID,
			}, nil),
		s.providerKeeper.EXPECT().
			GetValidatorSetUpdateId(gomock.Any()).
			Return(uint64(1)),
		s.providerKeeper.EXPECT().
			AppendPendingVSCPackets(gomock.Any(), gomock.Eq(consumerID), gomock.Any()),
		s.providerKeeper.EXPECT().
			IncrementValidatorSetUpdateId(gomock.Any()),

		// First InitConsumers call
		s.providerKeeper.EXPECT().
			GetConsumerPhase(gomock.Any(), gomock.Eq(consumerID)).
			Return(ccvprovidertypes.CONSUMER_PHASE_LAUNCHED),
		s.providerKeeper.EXPECT().
			GetConsumerIdToChannelId(gomock.Any(), gomock.Eq(consumerID)).
			Return("", false),

		// Second InitConsumers call
		s.providerKeeper.EXPECT().
			GetConsumerPhase(gomock.Any(), gomock.Eq(consumerID)).
			Return(ccvprovidertypes.CONSUMER_PHASE_LAUNCHED),
		s.providerKeeper.EXPECT().
			GetConsumerIdToChannelId(gomock.Any(), gomock.Eq(consumerID)).
			Return("channel-42", true),
		s.providerKeeper.EXPECT().
			SendVSCPacketsToChain(gomock.Any(), gomock.Eq(consumerID), gomock.Eq("channel-42")),
	)

	// Create a stack
	ver := "1.2.3"
	_, err := s.msgServer.CreateChainletStack(s.ctx, types.NewMsgCreateChainletStack(
		creator.String(), "test", "test", "test/test:"+ver, ver, "abcd"+ver, fees, true,
	))
	s.Require().NoError(err)

	// Launch a chainlet
	_, err = s.msgServer.LaunchChainlet(s.ctx, types.NewMsgLaunchChainlet(
		creator.String(), nil, "test", ver, "test_chainlet", chainID, "asaga", types.ChainletParams{}, nil, false, "",
	))
	s.Require().NoError(err)

	// VSC not sent without an open channel
	s.chainletKeeper.InitConsumers(s.ctx)

	// VSC sent
	s.chainletKeeper.InitConsumers(s.ctx)
}
