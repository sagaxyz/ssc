package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	keepertest "github.com/sagaxyz/ssc/testutil/keeper"

	"github.com/sagaxyz/ssc/x/chainlet/keeper"
	"github.com/sagaxyz/ssc/x/chainlet/types"
)

//nolint:unused
func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.ChainletKeeper(t)
	return keeper.NewMsgServerImpl(k), ctx
}

// TestLaunchChainlet_CustomLauncherEvent verifies that when a custom launcher is set,
// the event emission correctly uses the custom launcher instead of the message creator.
func (s *TestSuite) TestLaunchChainlet_CustomLauncherEvent() {
	// Setup: Create a chainlet stack with fees
	stackName := "test-stack-custom"
	stackVersion := "1.0.0"

	_, err := s.msgServer.CreateChainletStack(s.ctx, types.NewMsgCreateChainletStack(
		creator.String(), stackName, "test description", "test/test:"+stackVersion, stackVersion, "abcd"+stackVersion, fees, false,
	))
	s.Require().NoError(err)

	// Setup: Configure mocks for billing and escrow
	customLauncher := sdk.AccAddress("custom_launcher").String()
	chainID := "test_12345-1" // Valid chain ID format: lowercase_letters_numbers-numbers

	s.aclKeeper.EXPECT().
		IsAdmin(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()

	s.billingKeeper.EXPECT().
		BillAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	s.escrowKeeper.EXPECT().
		NewChainletAccount(gomock.Any(), gomock.Any(), gomock.Eq(chainID), gomock.Any()).
		Return(nil).
		AnyTimes()

	// Clear any existing events
	s.ctx = s.ctx.WithEventManager(sdk.NewEventManager())

	// Test: Launch chainlet with custom launcher
	msg := types.NewMsgLaunchChainlet(
		creator.String(),
		[]string{creator.String()},
		stackName,
		stackVersion,
		"test-chainlet",
		chainID,
		"utsaga",
		types.ChainletParams{},
		nil,
		false,
		customLauncher, // Custom launcher
	)

	_, err = s.msgServer.LaunchChainlet(s.ctx, msg)
	s.Require().NoError(err)

	// Verify: Check that the event was emitted with the custom launcher
	events := s.ctx.EventManager().Events()
	var foundLauncherEvent bool
	for _, event := range events {
		if event.Type == "ssc.chainlet.EventLaunchChainlet" {
			// Find the launcher attribute
			for _, attr := range event.Attributes {
				if string(attr.Key) == "launcher" {
					launcherValue := string(attr.Value)
					// Remove JSON quotes if present
					if len(launcherValue) > 0 && launcherValue[0] == '"' && launcherValue[len(launcherValue)-1] == '"' {
						launcherValue = launcherValue[1 : len(launcherValue)-1]
					}
					s.Require().Equal(customLauncher, launcherValue,
						"Event should contain custom launcher (%s), not message creator (%s). Got: %s", customLauncher, creator.String(), launcherValue)
					foundLauncherEvent = true
					break
				}
			}
		}
	}
	s.Require().True(foundLauncherEvent, "LaunchChainlet event should be emitted with launcher attribute")

	// Verify: Also check the chainlet was created with correct launcher
	chainlet, err := s.chainletKeeper.GetChainletInfo(s.ctx, chainID)
	s.Require().NoError(err)
	s.Require().Equal(customLauncher, chainlet.Launcher,
		"Chainlet should be stored with custom launcher (%s), got %s", customLauncher, chainlet.Launcher)
}
