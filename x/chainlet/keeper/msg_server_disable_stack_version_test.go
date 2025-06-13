package keeper_test

import (
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (s *TestSuite) TestDisabledVersionsLaunch() {
	s.SetupTest()

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
		creator.String(), "test", "test", "test/test:"+ver, ver, "abcd"+ver, fees, false,
	))
	s.Require().NoError(err)

	// Launch a chainlet
	_, err = s.msgServer.LaunchChainlet(s.ctx, types.NewMsgLaunchChainlet(
		creator.String(), nil, "test", ver, "test_chainlet", "test_12345-1", "asaga", types.ChainletParams{}, nil, false,
	))
	s.Require().NoError(err)

	// Disable the stack version
	_, err = s.msgServer.DisableChainletStackVersion(s.ctx, types.NewMsgDisableChainletStackVersion(creator.String(), "test", ver))
	s.Require().NoError(err)

	// Try and fail to launch another chainlet
	_, err = s.msgServer.LaunchChainlet(s.ctx, types.NewMsgLaunchChainlet(
		creator.String(), nil, "test", ver, "test_chainlet", "test_12346-1", "asaga", types.ChainletParams{}, nil, false,
	))
	s.Require().Error(err)
}

func (s *TestSuite) TestDisabledVersionsUpgrade() {
	s.SetupTest()

	// Create a stack
	ver := "1.2.3"
	_, err := s.msgServer.CreateChainletStack(s.ctx, types.NewMsgCreateChainletStack(
		creator.String(), "test", "test", "test/test:"+ver, ver, "abcd"+ver, fees, false,
	))
	s.Require().NoError(err)

	// Launch a chainlet
	s.escrowKeeper.EXPECT().
		NewChainletAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)
	s.billingKeeper.EXPECT().
		BillAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)
	_, err = s.msgServer.LaunchChainlet(s.ctx, types.NewMsgLaunchChainlet(
		creator.String(), nil, "test", ver, "test_chainlet", "test_12345-1", "asaga", types.ChainletParams{}, nil, false,
	))
	s.Require().NoError(err)

	// Create a newer but disabled stack version
	ver2 := "1.2.4"
	_, err = s.msgServer.UpdateChainletStack(s.ctx, types.NewMsgUpdateChainletStack(
		creator.String(), "test", "test/test:"+ver2, ver2, "abcd"+ver2, false,
	))
	s.Require().NoError(err)
	_, err = s.msgServer.DisableChainletStackVersion(s.ctx, types.NewMsgDisableChainletStackVersion(creator.String(), "test", ver2))
	s.Require().NoError(err)

	// Try and fail to upgrade the chainlet to a disabled stack version
	_, err = s.msgServer.UpgradeChainlet(s.ctx, types.NewMsgUpgradeChainlet(
		creator.String(), "test_12346-1", ver2, 0,
	))
	s.Require().Error(err)
}

func (s *TestSuite) TestDisabledVersionAutoUpgrade() {
	tests := []struct {
		addedVersions    []string
		disabledVersions []string

		current        string
		expectedLatest string
	}{
		{
			[]string{"0.1.2"},
			[]string{"0.1.2"},
			"0.1.2", "0.1.2",
		},
		{
			[]string{"0.1.2", "0.1.3"},
			[]string{"0.1.2"},
			"0.1.2", "0.1.3",
		},
		{
			[]string{"0.1.2", "0.1.3"},
			[]string{"0.1.2", "0.1.3"},
			"0.1.2", "0.1.2",
		},
		{
			[]string{"0.1.2", "0.1.3", "0.1.4"},
			[]string{"0.1.2", "0.1.4"},
			"0.1.2", "0.1.3",
		},
	}

	for i, tt := range tests {
		s.Run(fmt.Sprintf("%d: %s -> %s", i, tt.current, tt.expectedLatest), func() {
			s.SetupTest()

			var err error
			// Add all stack versions
			for j, ver := range tt.addedVersions {
				if j == 0 {
					_, err = s.msgServer.CreateChainletStack(s.ctx, types.NewMsgCreateChainletStack(
						creator.String(), "test", "test", "test/test:"+ver, ver, "abcd"+ver, fees, false,
					))
					s.Require().NoError(err)
				} else {
					_, err = s.msgServer.UpdateChainletStack(s.ctx, types.NewMsgUpdateChainletStack(
						creator.String(), "test", "test/test:"+ver, ver, "abcd"+ver, false,
					))
					s.Require().NoError(err)
				}
			}
			// Launch a testing chainlet
			s.escrowKeeper.EXPECT().
				NewChainletAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)
			s.billingKeeper.EXPECT().
				BillAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)
			chainId := "test_12345-42"
			_, err = s.msgServer.LaunchChainlet(s.ctx, types.NewMsgLaunchChainlet(
				creator.String(), nil, "test", tt.current, "test_chainlet", chainId, "asaga", types.ChainletParams{}, nil, false,
			))
			s.Require().NoError(err)

			// Disable specified stack versions
			for _, ver := range tt.disabledVersions {
				_, err = s.msgServer.DisableChainletStackVersion(s.ctx, types.NewMsgDisableChainletStackVersion(creator.String(), "test", ver))
				s.Require().NoError(err)
			}

			// Check it directly
			lv, err := s.chainletKeeper.LatestVersion(s.ctx, "test", tt.current)
			s.Require().NoError(err)
			s.Require().Equal(tt.expectedLatest, lv)

			// Check it with a chainlet auto-upgrade
			err = s.chainletKeeper.AutoUpgradeChainlets(s.ctx)
			s.Require().NoError(err)
			chainlet, err := s.chainletKeeper.Chainlet(s.ctx, chainId)
			s.Require().NoError(err)
			s.Require().Equal(tt.expectedLatest, chainlet.ChainletStackVersion)
		})
	}
}
