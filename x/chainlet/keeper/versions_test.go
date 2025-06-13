package keeper_test

import (
	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (s *TestSuite) TestVersionsLoading() {
	tests := []struct {
		addedVersions    []string
		disabledVersions []string
		expectedState    []string
	}{
		{
			[]string{"0.1.2"},
			[]string{"0.1.2"},
			[]string{},
		},
		{
			[]string{"0.1.2", "1.2.3", "2.0.0"},
			[]string{},
			[]string{"0.1.2", "1.2.3", "2.0.0"},
		},
		{
			[]string{"0.1.2", "0.1.3", "1.2.3", "1.2.4", "2.0.0", "2.1.0", "2.1.1"},
			[]string{},
			[]string{"0.1.2", "0.1.3", "1.2.3", "1.2.4", "2.0.0", "2.1.0", "2.1.1"},
		},
		{
			[]string{"0.1.2", "0.1.3"},
			[]string{"0.1.2"},
			[]string{"0.1.3"},
		},
		{
			[]string{"0.1.2", "1.2.3", "2.0.0"},
			[]string{"0.1.2", "1.2.3", "2.0.0"},
			[]string{},
		},
		{
			[]string{"0.1.2", "1.2.3", "2.0.0"},
			[]string{"0.1.2", "2.0.0"},
			[]string{"1.2.3"},
		},
		{
			[]string{"0.1.2", "0.1.3", "1.2.3", "1.2.4", "2.0.0", "2.1.0", "2.1.1"},
			[]string{"0.1.2", "0.1.3", "1.2.3", "2.0.0", "2.1.0", "2.1.1"},
			[]string{"1.2.4"},
		},
	}

	for _, tt := range tests {
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

		// Disable specified stack versions
		for _, ver := range tt.disabledVersions {
			_, err = s.msgServer.DisableChainletStackVersion(s.ctx, types.NewMsgDisableChainletStackVersion(creator.String(), "test", ver))
			s.Require().NoError(err)
		}

		s.chainletKeeper.DeleteVersions()
		// Force re-load
		_, err = s.msgServer.CreateChainletStack(s.ctx, types.NewMsgCreateChainletStack(
			creator.String(), "xxx", "xxx", "xxx", "1.2.3", "abcd", fees, false,
		))
		s.Require().NoError(err)
		versions := s.chainletKeeper.Versions("test")
		s.Require().Equal(tt.expectedState, versions)

		s.chainletKeeper.DeleteVersions()
		// Force re-load
		_, err = s.chainletKeeper.LatestVersion(s.ctx, "test", "1.2.3")
		s.Require().NoError(err)
		versions = s.chainletKeeper.Versions("test")
		s.Require().Equal(tt.expectedState, versions)
	}
}
