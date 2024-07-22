package keeper_test

import (
	"github.com/golang/mock/gomock"

	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (s *TestSuite) TestConsumerVSC() {
	s.SetupTest()

	const chainID = "test_12345-1"

	s.escrowKeeper.EXPECT().
		NewChainletAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	s.billingKeeper.EXPECT().
		BillAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	// Create a stack
	ver := "1.2.3"
	_, err := s.msgServer.CreateChainletStack(s.ctx, types.NewMsgCreateChainletStack(
		creator.String(), "test", "test", "test/test:"+ver, ver, "abcd"+ver, fees,
	))
	s.Require().NoError(err)

	// Launch a chainlet
	_, err = s.msgServer.LaunchChainlet(s.ctx, types.NewMsgLaunchChainlet(
		creator.String(), nil, "test", ver, "test_chainlet", chainID, "asaga", types.ChainletParams{},
	))
	s.Require().NoError(err)

	// VSC not sent without an open channel
	s.providerKeeper.EXPECT().
		GetChainToChannel(gomock.Any(), chainID).
		Return("", false)
	s.providerKeeper.EXPECT().
		SendVSCPacketsToChain(gomock.Any(), chainID, "channel-42").
		Return().
		Times(0)
	s.chainletKeeper.ForcePendingVSC(s.ctx)

	// VSC sent
	s.providerKeeper.EXPECT().
		GetChainToChannel(gomock.Any(), chainID).
		Return("channel-42", true)
	s.providerKeeper.EXPECT().
		SendVSCPacketsToChain(gomock.Any(), chainID, "channel-42").
		Return()
	s.chainletKeeper.ForcePendingVSC(s.ctx)
}
