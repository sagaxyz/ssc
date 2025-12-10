package keeper_test

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v10/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ccvprovidertypes "github.com/cosmos/interchain-security/v7/x/ccv/provider/types"
	"github.com/golang/mock/gomock"
	sdkchainlettypes "github.com/sagaxyz/saga-sdk/x/chainlet/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (s *TestSuite) TestUpgradeChainlet() {
	var (
		clientID     = "client-123"
		channelID    = "channel-42"
		connectionID = "connection-0"
		consumerID   = "0"
	)

	testCases := []struct {
		name        string
		fromVersion string
		fromCCV     bool
		toVersion   string
		toCCV       bool
		sender      sdk.AccAddress
		mocks       func(s *TestSuite)
		expErr      string
	}{
		{
			"ok - legacy-to-legacy",
			"1.2.3", false,
			"2.0.0", false,
			admin,
			nil,
			"",
		},
		{
			"ok - legacy-to-CCV",
			"1.2.3", false,
			"2.0.0", true,
			admin,
			func(s *TestSuite) {
				// Added to consumers during upgrade
				gomock.InOrder(
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
				)
			},
			"",
		},
		{
			"ok - CCV-to-CCV",
			"1.2.3", true,
			"2.0.0", true,
			maintainer,
			func(s *TestSuite) {
				gomock.InOrder(
					s.providerKeeper.EXPECT().
						GetConsumerClientId(gomock.Any(), gomock.Eq(consumerID)).
						Return(clientID, true),
					s.channelKeeper.EXPECT().
						GetChannel(gomock.Any(), sdkchainlettypes.PortID, gomock.Eq(channelID)).
						Return(ibcchanneltypes.Channel{
							ConnectionHops: []string{connectionID},
						}, true),
					s.connectionKeeper.EXPECT().
						GetConnection(gomock.Any(), gomock.Eq(connectionID)).
						Return(ibcconnectiontypes.ConnectionEnd{
							ClientId:     clientID,
							Versions:     []*ibcconnectiontypes.Version{},
							State:        0,
							Counterparty: ibcconnectiontypes.Counterparty{},
							DelayPeriod:  0,
						}, true),
					s.clientKeeper.EXPECT().
						GetClientLatestHeight(gomock.Any(), gomock.Eq(clientID)).
						Return(ibcclienttypes.Height{}),
					s.channelKeeper.EXPECT().
						SendPacket(
							gomock.Any(),
							gomock.Eq(sdkchainlettypes.PortID),
							gomock.Eq(channelID),
							gomock.Any(),
							gomock.Any(),
							gomock.Any(),
						).
						Return(uint64(1337), nil),
				)
			},
			"",
		},
		{
			"fail - CCV-to-legacy",
			"1.2.3", true,
			"2.0.0", false,
			maintainer,
			nil,
			"cannot be disabled",
		},
		{
			"fail - skip upgrade",
			"1.2.3", false,
			"3.0.0", false,
			admin,
			nil,
			"increments of",
		},
		{
			"fail - legacy upgrade as maintainer",
			"1.2.3", false,
			"2.0.0", false,
			maintainer,
			nil,
			"not allowed",
		},
		{
			"fail - legacy upgrade as maintainer",
			"1.2.3", false,
			"2.0.0", true,
			maintainer,
			nil,
			"not allowed",
		},
		{
			"fail - not maintainer",
			"1.2.3", true,
			"2.0.0", true,
			creator,
			nil,
			"not a chainlet maintainer",
		},
		{
			"fail - enable CCV with non-breaking upgrade",
			"1.0.2", false,
			"1.1.0", true,
			admin,
			nil,
			"requires a breaking upgrade",
		},
		{
			"fail - disable CCV with non-breaking upgrade",
			"1.0.2", true,
			"1.1.0", false,
			admin,
			nil,
			"requires a breaking upgrade",
		},
	}
	for i, tc := range testCases {
		s.Run(fmt.Sprintf("%d: %s", i, tc.name), func() {
			s.SetupTest()

			// Mocks we do not care about
			s.escrowKeeper.EXPECT().
				NewChainletAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()
			s.billingKeeper.EXPECT().
				BillAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()

			// Mock IsAdmin so our admin address is the (only) admin
			s.aclKeeper.EXPECT().
				IsAdmin(gomock.Any(), gomock.Eq(tc.sender)).
				Return(tc.sender.String() == admin.String()).
				AnyTimes()
			s.aclKeeper.EXPECT().
				IsAdmin(gomock.Any(), gomock.Eq(creator)).
				Return(false).
				AnyTimes()

			// Create stack versions
			_, err := s.msgServer.CreateChainletStack(s.ctx, types.NewMsgCreateChainletStack(
				creator.String(), "test", "test", "test/test:"+tc.fromVersion, tc.fromVersion, "abcd"+tc.fromVersion, fees, tc.fromCCV,
			))
			s.Require().NoError(err)
			_, err = s.msgServer.UpdateChainletStack(s.ctx, types.NewMsgUpdateChainletStack(
				creator.String(), "test", "test/test:"+tc.toVersion, tc.toVersion, "abcd"+tc.toVersion, tc.toCCV,
			))
			s.Require().NoError(err)

			// Launch a chainlet
			if tc.fromCCV {
				// Implies adding to consumers when launching
				s.providerMsgServer.EXPECT().
					CreateConsumer(gomock.Any(), gomock.Any()).
					Return(&ccvprovidertypes.MsgCreateConsumerResponse{
						ConsumerId: consumerID,
					}, nil)
				s.providerKeeper.EXPECT().
					GetValidatorSetUpdateId(gomock.Any()).
					Return(uint64(1))
				s.providerKeeper.EXPECT().
					AppendPendingVSCPackets(gomock.Any(), gomock.Eq(consumerID), gomock.Any())
				s.providerKeeper.EXPECT().
					IncrementValidatorSetUpdateId(gomock.Any())
			}
			chainID := fmt.Sprintf("test_%d-1", i+1)
			_, err = s.msgServer.LaunchChainlet(s.ctx, types.NewMsgLaunchChainlet(
				creator.String(), []string{maintainer.String()}, "test", tc.fromVersion, "test_chainlet", chainID, "asaga", types.ChainletParams{}, nil, false, "",
			))
			s.Require().NoError(err)

			// Upgrade the chainlet
			if tc.mocks != nil {
				tc.mocks(s)
			}
			_, err = s.msgServer.UpgradeChainlet(s.ctx, types.NewMsgUpgradeChainlet(
				tc.sender.String(), chainID, tc.toVersion, 0, channelID, nil,
			))
			if tc.expErr == "" {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
				if !strings.Contains(err.Error(), tc.expErr) {
					s.Require().Fail(fmt.Sprintf("err '%s' does not contain '%s'", err.Error(), tc.expErr))
				}
			}
		})
	}
}
