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

func (s *TestSuite) TestVersionExistsInCache() {
	s.SetupTest()

	stackName := "test"
	version1 := "1.2.3"
	version2 := "2.0.0"
	nonExistentVersion := "3.0.0"

	// Create a chainlet stack with version 1.2.3
	_, err := s.msgServer.CreateChainletStack(s.ctx, types.NewMsgCreateChainletStack(
		creator.String(), stackName, "test", "test/test:"+version1, version1, "abcd"+version1, fees, true,
	))
	s.Require().NoError(err)

	// Add version 2.0.0
	_, err = s.msgServer.UpdateChainletStack(s.ctx, types.NewMsgUpdateChainletStack(
		creator.String(), stackName, "test/test:"+version2, version2, "abcd"+version2, true,
	))
	s.Require().NoError(err)

	// Test 1: Version exists in cache (1.2.3)
	exists := s.chainletKeeper.VersionExistsInCache(s.ctx, stackName, version1)
	s.Require().True(exists, "version 1.2.3 should exist in cache")

	// Test 2: Version exists in cache (2.0.0)
	exists = s.chainletKeeper.VersionExistsInCache(s.ctx, stackName, version2)
	s.Require().True(exists, "version 2.0.0 should exist in cache")

	// Test 3: Version doesn't exist
	exists = s.chainletKeeper.VersionExistsInCache(s.ctx, stackName, nonExistentVersion)
	s.Require().False(exists, "version 3.0.0 should not exist in cache")

	// Test 4: Stack doesn't exist
	exists = s.chainletKeeper.VersionExistsInCache(s.ctx, "nonexistent", version1)
	s.Require().False(exists, "version should not exist for nonexistent stack")

	// Test 5: Version with 'v' prefix normalization
	exists = s.chainletKeeper.VersionExistsInCache(s.ctx, stackName, "v"+version1)
	s.Require().True(exists, "version with 'v' prefix should be normalized and found")

	// Test 6: Version with 'V' prefix normalization
	exists = s.chainletKeeper.VersionExistsInCache(s.ctx, stackName, "V"+version2)
	s.Require().True(exists, "version with 'V' prefix should be normalized and found")

	// Test 7: Cache is nil - should load and check
	s.chainletKeeper.DeleteVersions()
	exists = s.chainletKeeper.VersionExistsInCache(s.ctx, stackName, version1)
	s.Require().True(exists, "should load cache and find version after DeleteVersions")

	// Test 8: Disabled version should not exist in cache
	// First add an enabled version, then disable it
	enabledVersion := "4.0.0"
	_, err = s.msgServer.UpdateChainletStack(s.ctx, types.NewMsgUpdateChainletStack(
		creator.String(), stackName, "test/test:"+enabledVersion, enabledVersion, "abcd"+enabledVersion, true, // enabled
	))
	s.Require().NoError(err)
	exists = s.chainletKeeper.VersionExistsInCache(s.ctx, stackName, enabledVersion)
	s.Require().True(exists, "enabled version should exist in cache")

	// Now disable it - this removes it from cache
	_, err = s.msgServer.DisableChainletStackVersion(s.ctx, types.NewMsgDisableChainletStackVersion(creator.String(), stackName, enabledVersion))
	s.Require().NoError(err)
	exists = s.chainletKeeper.VersionExistsInCache(s.ctx, stackName, enabledVersion)
	s.Require().False(exists, "disabled version should not exist in cache after disabling")
}
