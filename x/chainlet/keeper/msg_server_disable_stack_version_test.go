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
		creator.String(), "test", "test", "test/test:"+ver, ver, "abcd"+ver, fees,
	))
	s.Require().NoError(err)

	// Launch a chainlet
	_, err = s.msgServer.LaunchChainlet(s.ctx, types.NewMsgLaunchChainlet(
		creator.String(), nil, "test", ver, "test_chainlet", "test_12345-1", "asaga", types.ChainletParams{},
	))
	s.Require().NoError(err)

	// Disable the stack version
	_, err = s.msgServer.DisableChainletStackVersion(s.ctx, types.NewMsgDisableChainletStackVersion(creator.String(), "test", ver))
	s.Require().NoError(err)

	// Try and fail to launch another chainlet
	_, err = s.msgServer.LaunchChainlet(s.ctx, types.NewMsgLaunchChainlet(
		creator.String(), nil, "test", ver, "test_chainlet", "test_12346-1", "asaga", types.ChainletParams{},
	))
	s.Require().Error(err)
}

func (s *TestSuite) TestDisabledVersionsUpgrade() {
	s.SetupTest()

	// Create a stack
	ver := "1.2.3"
	_, err := s.msgServer.CreateChainletStack(s.ctx, types.NewMsgCreateChainletStack(
		creator.String(), "test", "test", "test/test:"+ver, ver, "abcd"+ver, fees,
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
		creator.String(), nil, "test", ver, "test_chainlet", "test_12345-1", "asaga", types.ChainletParams{},
	))
	s.Require().NoError(err)

	// Create a newer but disabled stack version
	ver2 := "1.2.4"
	_, err = s.msgServer.UpdateChainletStack(s.ctx, types.NewMsgUpdateChainletStack(
		creator.String(), "test", "test/test:"+ver2, ver2, "abcd"+ver2,
	))
	s.Require().NoError(err)
	_, err = s.msgServer.DisableChainletStackVersion(s.ctx, types.NewMsgDisableChainletStackVersion(creator.String(), "test", ver2))
	s.Require().NoError(err)

	// Try and fail to upgrade the chainlet to a disabled stack version
	_, err = s.msgServer.UpgradeChainlet(s.ctx, types.NewMsgUpgradeChainlet(
		creator.String(), "test_12346-1", ver2,
	))
	s.Require().Error(err)
}

func (s *TestSuite) TestDisabledVersionAutoUpgrade() {
	type chainlet struct {
		chainletStackName string
		chainletName      string
		chainletID        string

		addedVersions    []string
		disabledVersions []string

		current        string
		expectedLatest string
	}

	tests := []struct {
		chainlets []chainlet
		name      string
	}{
		{
			name: "Add one version 0.1.2 and than disabled it",
			chainlets: []chainlet{
				{
					chainletStackName: "saga-stack",
					chainletName:      "saga-chain",
					chainletID:        "test_12345-42",
					addedVersions:     []string{"0.1.2"},
					disabledVersions:  []string{"0.1.2"},
					current:           "0.1.2",
					expectedLatest:    "0.1.2",
				},
				{
					chainletStackName: "xyz-stack",
					chainletName:      "xyz-chain",
					chainletID:        "test_12345-43",
					addedVersions:     []string{"0.1.2"},
					disabledVersions:  []string{"0.1.2"},
					current:           "0.1.2",
					expectedLatest:    "0.1.2",
				},
			},
		},
		{
			name: "update from 0.1.2 to 0.1.3",
			chainlets: []chainlet{
				{
					chainletStackName: "saga-stack",
					chainletName:      "saga-chain",
					chainletID:        "test_12345-44",
					addedVersions:     []string{"0.1.2", "0.1.3"},
					disabledVersions:  []string{"0.1.2"},
					current:           "0.1.2",
					expectedLatest:    "0.1.3",
				},
				{
					chainletStackName: "xyz-stack",
					chainletName:      "xyz-chain",
					chainletID:        "test_12345-45",
					addedVersions:     []string{"0.1.2", "0.1.3"},
					disabledVersions:  []string{"0.1.2"},
					current:           "0.1.2",
					expectedLatest:    "0.1.3",
				},
			},
		},
		{
			name: "Add two versions and than disables them",
			chainlets: []chainlet{
				{
					chainletStackName: "saga-stack",
					chainletName:      "saga-chain",
					chainletID:        "test_12345-46",
					addedVersions:     []string{"0.1.2", "0.1.3"},
					disabledVersions:  []string{"0.1.2", "0.1.3"},
					current:           "0.1.2",
					expectedLatest:    "0.1.2",
				},
				{
					chainletStackName: "xyz-stack",
					chainletName:      "xyz-chain",
					chainletID:        "test_12345-47",
					addedVersions:     []string{"0.1.2", "0.1.3"},
					disabledVersions:  []string{"0.1.2", "0.1.3"},
					current:           "0.1.2",
					expectedLatest:    "0.1.2",
				},
			},
		},
		{
			name: "Add three versions and then disable the first one and the last one",
			chainlets: []chainlet{
				{
					chainletStackName: "saga-stack",
					chainletName:      "saga-chain",
					chainletID:        "test_12345-48",
					addedVersions:     []string{"0.1.2", "0.1.3", "0.1.4"},
					disabledVersions:  []string{"0.1.2", "0.1.4"},
					current:           "0.1.2",
					expectedLatest:    "0.1.3",
				},
				{
					chainletStackName: "xyz-stack",
					chainletName:      "xyz-chain",
					chainletID:        "test_12345-49",
					addedVersions:     []string{"0.1.2", "0.1.3", "0.1.4"},
					disabledVersions:  []string{"0.1.2", "0.1.4"},
					current:           "0.1.2",
					expectedLatest:    "0.1.3",
				},
			},
		},
	}

	for _, ttt := range tests {
		s.Run(fmt.Sprintf(ttt.name), func() {
			s.SetupTest()
			for _, tt := range ttt.chainlets {
				var err error
				// Add all stack versions
				for j, ver := range tt.addedVersions {
					image := fmt.Sprintf("%s/%s:%s", tt.chainletStackName, tt.chainletStackName, ver)
					if j == 0 {
						_, err = s.msgServer.CreateChainletStack(s.ctx, types.NewMsgCreateChainletStack(
							creator.String(), tt.chainletStackName, tt.chainletStackName, image, ver, "abcd"+ver, fees,
						))
						s.Require().NoError(err)
					} else {
						_, err = s.msgServer.UpdateChainletStack(s.ctx, types.NewMsgUpdateChainletStack(
							creator.String(), tt.chainletStackName, image, ver, "abcd"+ver,
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
				_, err = s.msgServer.LaunchChainlet(s.ctx, types.NewMsgLaunchChainlet(
					creator.String(), nil, tt.chainletStackName, tt.current, tt.chainletName, tt.chainletID, "asaga", types.ChainletParams{},
				))
				s.Require().NoError(err)

				// Disable specified stack versions
				for _, ver := range tt.disabledVersions {
					_, err = s.msgServer.DisableChainletStackVersion(s.ctx, types.NewMsgDisableChainletStackVersion(creator.String(), tt.chainletStackName, ver))
					s.Require().NoError(err)
				}

				// Check it directly
				lv, err := s.chainletKeeper.LatestVersion(s.ctx, tt.chainletStackName, tt.current)
				s.Require().NoError(err)
				s.Require().Equal(tt.expectedLatest, lv)

				// Check it with a chainlet auto-upgrade
				err = s.chainletKeeper.AutoUpgradeChainlets(s.ctx)
				s.Require().NoError(err)
				got, err := s.chainletKeeper.Chainlet(s.ctx, tt.chainletID)
				s.Require().NoError(err)
				s.Require().Equal(tt.expectedLatest, got.ChainletStackVersion)
			}
		})
	}
}
