package keeper_test

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v10/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ccvprovidertypes "github.com/cosmos/interchain-security/v7/x/ccv/provider/types"
	"github.com/golang/mock/gomock"
	chainlettypes "github.com/sagaxyz/saga-sdk/x/chainlet/types"
	sdkchainlettypes "github.com/sagaxyz/saga-sdk/x/chainlet/types"

	"github.com/sagaxyz/ssc/x/chainlet/keeper"
	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (s *TestSuite) ibcSetup(chainID, consumerID, channelID string) {
	// Calls we don't care about in these tests
	s.escrowKeeper.EXPECT().
		NewChainletAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	s.billingKeeper.EXPECT().
		BillAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	s.providerKeeper.EXPECT().
		GetValidatorSetUpdateId(gomock.Any()).
		Return(uint64(1)).
		AnyTimes()

	// Create stacks
	ver := "1.2.3"
	_, err := s.msgServer.CreateChainletStack(s.ctx, types.NewMsgCreateChainletStack(
		creator.String(), "test", "test", "test/test:"+ver, ver, "abcd"+ver, fees, true,
	))
	s.Require().NoError(err)
	_, err = s.msgServer.UpdateChainletStack(s.ctx, types.NewMsgUpdateChainletStack(
		creator.String(), "test", "test/test:2.0.0", "2.0.0", "xyz", true,
	))
	s.Require().NoError(err)

	// Setup mocks with the correct chain ID and consumer ID
	s.providerMsgServer.EXPECT().
		CreateConsumer(gomock.Any(), gomock.Any()).
		Return(&ccvprovidertypes.MsgCreateConsumerResponse{
			ConsumerId: consumerID,
		}, nil)
	s.providerKeeper.EXPECT().
		AppendPendingVSCPackets(gomock.Any(), gomock.Eq(consumerID), gomock.Any()).
		AnyTimes()
	s.providerKeeper.EXPECT().
		IncrementValidatorSetUpdateId(gomock.Any()).
		AnyTimes()
	s.providerKeeper.EXPECT().
		GetConsumerIdToChannelId(gomock.Any(), gomock.Eq(consumerID)).
		Return(channelID, true).
		AnyTimes()
	s.providerKeeper.EXPECT().
		SendVSCPacketsToChain(gomock.Any(), gomock.Eq(consumerID), gomock.Eq(channelID)).
		AnyTimes()
	s.providerKeeper.EXPECT().
		GetConsumerPhase(gomock.Any(), gomock.Eq(consumerID)).
		Return(ccvprovidertypes.CONSUMER_PHASE_LAUNCHED).
		AnyTimes()

	// Launch a chainlet
	_, err = s.msgServer.LaunchChainlet(s.ctx, types.NewMsgLaunchChainlet(
		creator.String(),
		[]string{creator.String()},
		"test",
		ver,
		"test_chainlet",
		chainID,
		"asaga",
		types.ChainletParams{},
		nil, false, "",
	))
	s.Require().NoError(err)
	s.chainletKeeper.InitConsumers(s.ctx)
}
func (s *TestSuite) breakingUpgrade(chainID, consumerID, clientID, connectionID, channelID string) {
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
	_, err := s.msgServer.UpgradeChainlet(s.ctx, &types.MsgUpgradeChainlet{
		Creator:      creator.String(),
		ChainId:      chainID,
		StackVersion: "2.0.0",
		HeightDelta:  100,
		ChannelId:    channelID,
	})
	s.Require().NoError(err)

	// Check if upgrade is correctly set in the chainlet
	chainlet, err := s.chainletKeeper.Chainlet(s.ctx, chainID)
	s.Require().NoError(err)
	s.Require().NotNil(chainlet.Upgrade)
}

func (s *TestSuite) packetVerificationMocks(consumerID, consumerClientID, clientID, connectionID, channelID string) {
	gomock.InOrder(
		s.providerKeeper.EXPECT().
			GetConsumerClientId(gomock.Any(), gomock.Eq(consumerID)).
			Return(consumerClientID, true),
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
	)
}

func (s *TestSuite) TestCreateUpgradePacket() {
	tests := []struct {
		name string
		fn   func(chainID, consumerID, clientID, connectionID, channelID string)
	}{
		{
			name: "success ack",
			fn: func(chainID, consumerID, clientID, connectionID, channelID string) {
				packet := channeltypes.Packet{}
				data := chainlettypes.CreateUpgradePacketData{
					ChainId: chainID,
					Name:    "xxx",
					Height:  123,
					Info:    "xyz",
				}

				// Success ack without an upgrade in progress
				packetAck := chainlettypes.CreateUpgradePacketAck{}
				packetAckBytes, err := types.ModuleCdc.MarshalJSON(&packetAck)
				s.Require().NoError(err)
				ack := channeltypes.NewResultAcknowledgement(sdk.MustSortJSON(packetAckBytes))
				err = s.chainletKeeper.OnAcknowledgementCreateUpgradePacket(s.ctx, packet, data, ack)
				s.Require().NoError(err)

				// Error ack without an upgrade in progress
				ack = channeltypes.NewErrorAcknowledgement(errors.New("error"))
				err = s.chainletKeeper.OnAcknowledgementCreateUpgradePacket(s.ctx, packet, data, ack)
				s.Require().NoError(err)

				// Check chainlet is unaffected
				chainlet, err := s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().Nil(chainlet.Upgrade)
				s.Require().Equal("1.2.3", chainlet.ChainletStackVersion)

				// Upgrade
				s.breakingUpgrade(chainID, consumerID, clientID, connectionID, channelID)

				// Check if upgrade is correctly set/unset in the chainlet
				chainlet, err = s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)

				// Success ack
				packetAck = chainlettypes.CreateUpgradePacketAck{}
				packetAckBytes, err = types.ModuleCdc.MarshalJSON(&packetAck)
				s.Require().NoError(err)
				ack = channeltypes.NewResultAcknowledgement(sdk.MustSortJSON(packetAckBytes))
				err = s.chainletKeeper.OnAcknowledgementCreateUpgradePacket(s.ctx, packet, data, ack)
				s.Require().NoError(err)
			},
		}, {
			name: "valid error ack",
			fn: func(chainID, consumerID, clientID, connectionID, channelID string) {
				s.breakingUpgrade(chainID, consumerID, clientID, connectionID, channelID)

				// Get correct upgrade plan name
				chainlet, err := s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)
				planName, err := keeper.UpgradePlanName(chainlet.ChainletStackVersion, chainlet.Upgrade.Version)
				s.Require().NoError(err)

				// Error ack
				gomock.InOrder(
					// Verification for the source of the packet
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
				)
				packet := channeltypes.Packet{
					SourceChannel: channelID,
				}
				data := chainlettypes.CreateUpgradePacketData{
					ChainId: chainID,
					Name:    planName,
					Height:  123,
					Info:    "xyz",
				}
				ack := channeltypes.NewErrorAcknowledgement(errors.New("error"))
				err = s.chainletKeeper.OnAcknowledgementCreateUpgradePacket(s.ctx, packet, data, ack)
				s.Require().NoError(err)

				// Upgrade removed
				chainlet, err = s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().Nil(chainlet.Upgrade)
			},
		}, {
			name: "valid timeout",
			fn: func(chainID, consumerID, clientID, connectionID, channelID string) {
				// Upgrade it
				s.breakingUpgrade(chainID, consumerID, clientID, connectionID, channelID)

				// Get correct upgrade plan name
				chainlet, err := s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)
				planName, err := keeper.UpgradePlanName(chainlet.ChainletStackVersion, chainlet.Upgrade.Version)
				s.Require().NoError(err)

				packet := channeltypes.Packet{
					SourceChannel: channelID,
				}
				data := chainlettypes.CreateUpgradePacketData{
					ChainId: chainID,
					Name:    planName,
					Height:  123,
					Info:    "xyz",
				}

				s.packetVerificationMocks(consumerID, clientID, clientID, connectionID, channelID)
				err = s.chainletKeeper.OnTimeoutCreateUpgradePacket(s.ctx, packet, data)
				s.Require().NoError(err)

				// Upgrade removed
				chainlet, err = s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().Nil(chainlet.Upgrade)
			},
		}, {
			name: "error ack or timeout for different upgrade plan",
			fn: func(chainID, consumerID, clientID, connectionID, channelID string) {
				// Upgrade it
				s.breakingUpgrade(chainID, consumerID, clientID, connectionID, channelID)

				packet := channeltypes.Packet{
					SourceChannel: channelID,
				}
				data := chainlettypes.CreateUpgradePacketData{
					ChainId: chainID,
					Name:    "xxx", // incorrect
					Height:  123,
					Info:    "xyz",
				}

				// Error ack
				s.packetVerificationMocks(consumerID, clientID, clientID, connectionID, channelID)
				ack := channeltypes.NewErrorAcknowledgement(errors.New("error"))
				err := s.chainletKeeper.OnAcknowledgementCreateUpgradePacket(s.ctx, packet, data, ack)
				s.Require().NoError(err)

				// Upgrade NOT removed
				chainlet, err := s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)

				// Timeout
				s.packetVerificationMocks(consumerID, clientID, clientID, connectionID, channelID)
				err = s.chainletKeeper.OnTimeoutCreateUpgradePacket(s.ctx, packet, data)
				s.Require().NoError(err)

				// Upgrade NOT removed
				chainlet, err = s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)
			},
		}, {
			name: "error ack or timeout from incorrect client ID",
			fn: func(chainID, consumerID, clientID, connectionID, channelID string) {
				// Upgrade it
				s.breakingUpgrade(chainID, consumerID, clientID, connectionID, channelID)

				// Get correct upgrade plan name
				chainlet, err := s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)
				planName, err := keeper.UpgradePlanName(chainlet.ChainletStackVersion, chainlet.Upgrade.Version)
				s.Require().NoError(err)

				packet := channeltypes.Packet{
					SourceChannel: "channel-42",
				}
				data := chainlettypes.CreateUpgradePacketData{
					ChainId: chainID,
					Name:    planName,
					Height:  123,
					Info:    "xyz",
				}

				// Error ack
				s.packetVerificationMocks(consumerID, clientID, "bad-client", "bad-connection", "channel-42")
				ack := channeltypes.NewErrorAcknowledgement(errors.New("error"))
				err = s.chainletKeeper.OnAcknowledgementCreateUpgradePacket(s.ctx, packet, data, ack)
				s.Require().Error(err)

				// Upgrade NOT removed
				chainlet, err = s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)

				// Timeout
				s.packetVerificationMocks(consumerID, clientID, "bad-client", "bad-connection", "channel-42")
				err = s.chainletKeeper.OnTimeoutCreateUpgradePacket(s.ctx, packet, data)
				s.Require().Error(err)

				// Upgrade NOT removed
				chainlet, err = s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)
			},
		},
	}
	for i, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()

			chainID := fmt.Sprintf("chain_%d-1", i+1)
			consumerID := fmt.Sprintf("%d", i)
			clientID := fmt.Sprintf("client-%d", i)
			connectionID := fmt.Sprintf("connection-%d", i)
			channelID := fmt.Sprintf("channel-%d", i)

			s.ibcSetup(chainID, consumerID, channelID)

			tt.fn(chainID, consumerID, clientID, connectionID, channelID)
		})
	}
}

func (s *TestSuite) TestCancelUpgradePacket() {
	tests := []struct {
		name string
		fn   func(chainID, consumerID, clientID, connectionID, channelID string)
	}{
		{
			name: "success ack",
			fn: func(chainID, consumerID, clientID, connectionID, channelID string) {
				// Upgrade
				s.breakingUpgrade(chainID, consumerID, clientID, connectionID, channelID)

				// Get correct upgrade plan name
				chainlet, err := s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)
				planName, err := keeper.UpgradePlanName(chainlet.ChainletStackVersion, chainlet.Upgrade.Version)
				s.Require().NoError(err)

				packet := channeltypes.Packet{
					SourceChannel: channelID,
				}
				data := chainlettypes.CancelUpgradePacketData{
					ChainId: chainID,
					Plan:    planName,
				}

				// Success ack
				s.packetVerificationMocks(consumerID, clientID, clientID, connectionID, channelID)
				packetAck := chainlettypes.CancelUpgradePacketAck{}
				packetAckBytes, err := types.ModuleCdc.MarshalJSON(&packetAck)
				s.Require().NoError(err)
				ack := channeltypes.NewResultAcknowledgement(sdk.MustSortJSON(packetAckBytes))
				err = s.chainletKeeper.OnAcknowledgementCancelUpgradePacket(s.ctx, packet, data, ack)
				s.Require().NoError(err)

				// Check chainlet upgrade is removed
				chainlet, err = s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().Nil(chainlet.Upgrade)
			},
		}, {
			name: "valid error ack",
			fn: func(chainID, consumerID, clientID, connectionID, channelID string) {
				// Upgrade it
				s.breakingUpgrade(chainID, consumerID, clientID, connectionID, channelID)

				// Get correct upgrade plan name
				chainlet, err := s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)
				planName, err := keeper.UpgradePlanName(chainlet.ChainletStackVersion, chainlet.Upgrade.Version)
				s.Require().NoError(err)

				// Error ack
				packet := channeltypes.Packet{
					SourceChannel: channelID,
				}
				data := chainlettypes.CancelUpgradePacketData{
					ChainId: chainID,
					Plan:    planName,
				}
				ack := channeltypes.NewErrorAcknowledgement(errors.New("error"))
				err = s.chainletKeeper.OnAcknowledgementCancelUpgradePacket(s.ctx, packet, data, ack)
				s.Require().NoError(err)

				// Upgrade NOT removed
				chainlet, err = s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)
			},
		}, {
			name: "valid timeout",
			fn: func(chainID, consumerID, clientID, connectionID, channelID string) {
				// Upgrade it
				s.breakingUpgrade(chainID, consumerID, clientID, connectionID, channelID)

				// Get correct upgrade plan name
				chainlet, err := s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)
				planName, err := keeper.UpgradePlanName(chainlet.ChainletStackVersion, chainlet.Upgrade.Version)
				s.Require().NoError(err)

				packet := channeltypes.Packet{
					SourceChannel: channelID,
				}
				data := chainlettypes.CancelUpgradePacketData{
					ChainId: chainID,
					Plan:    planName,
				}

				err = s.chainletKeeper.OnTimeoutCancelUpgradePacket(s.ctx, packet, data)
				s.Require().NoError(err)

				// Upgrade NOT removed
				chainlet, err = s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)
			},
		}, {
			name: "ack for different upgrade plan",
			fn: func(chainID, consumerID, clientID, connectionID, channelID string) {
				// Upgrade it
				s.breakingUpgrade(chainID, consumerID, clientID, connectionID, channelID)

				packet := channeltypes.Packet{
					SourceChannel: channelID,
				}
				data := chainlettypes.CancelUpgradePacketData{
					ChainId: chainID,
					Plan:    "xxx", // incorrect
				}

				// Success ack
				s.packetVerificationMocks(consumerID, clientID, clientID, connectionID, channelID)
				packetAck := chainlettypes.CancelUpgradePacketAck{}
				packetAckBytes, err := types.ModuleCdc.MarshalJSON(&packetAck)
				s.Require().NoError(err)
				ack := channeltypes.NewResultAcknowledgement(sdk.MustSortJSON(packetAckBytes))
				err = s.chainletKeeper.OnAcknowledgementCancelUpgradePacket(s.ctx, packet, data, ack)
				s.Require().NoError(err)

				// Upgrade NOT removed
				chainlet, err := s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)
			},
		}, {
			name: "ack from incorrect client ID",
			fn: func(chainID, consumerID, clientID, connectionID, channelID string) {
				// Upgrade it
				s.breakingUpgrade(chainID, consumerID, clientID, connectionID, channelID)

				// Get correct upgrade plan name
				chainlet, err := s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)
				planName, err := keeper.UpgradePlanName(chainlet.ChainletStackVersion, chainlet.Upgrade.Version)
				s.Require().NoError(err)

				packet := channeltypes.Packet{
					SourceChannel: "channel-42",
				}
				data := chainlettypes.CancelUpgradePacketData{
					ChainId: chainID,
					Plan:    planName,
				}

				// Success ack
				s.packetVerificationMocks(consumerID, clientID, "bad-client", "bad-connection", "channel-42")
				packetAck := chainlettypes.CancelUpgradePacketAck{}
				packetAckBytes, err := types.ModuleCdc.MarshalJSON(&packetAck)
				s.Require().NoError(err)
				ack := channeltypes.NewResultAcknowledgement(sdk.MustSortJSON(packetAckBytes))
				err = s.chainletKeeper.OnAcknowledgementCancelUpgradePacket(s.ctx, packet, data, ack)
				s.Require().Error(err)

				// Upgrade NOT removed
				chainlet, err = s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)
			},
		},
	}
	for i, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()

			chainID := fmt.Sprintf("chain_%d-1", i+1)
			consumerID := fmt.Sprintf("%d", i)
			clientID := fmt.Sprintf("client-%d", i)
			connectionID := fmt.Sprintf("connection-%d", i)
			channelID := fmt.Sprintf("channel-%d", i)

			s.ibcSetup(chainID, consumerID, channelID)

			tt.fn(chainID, consumerID, clientID, connectionID, channelID)
		})
	}
}

func (s *TestSuite) TestConfirmUpgradePacket() {
	tests := []struct {
		name string
		fn   func(chainID, consumerID, clientID, connectionID, channelID string)
	}{
		{
			name: "ok",
			fn: func(chainID, consumerID, clientID, connectionID, channelID string) {
				s.breakingUpgrade(chainID, consumerID, clientID, connectionID, channelID)

				packet := channeltypes.Packet{
					DestinationChannel: channelID,
				}
				data := chainlettypes.ConfirmUpgradePacketData{
					ChainId: chainID,
					Height:  123,
					Plan:    "xyz",
				}

				// Upgrade confirmation
				s.packetVerificationMocks(consumerID, clientID, clientID, connectionID, channelID)
				_, err := s.chainletKeeper.OnRecvConfirmUpgradePacket(s.ctx, packet, data)
				s.Require().NoError(err)

				// Check if the upgrade is finished
				chainlet, err := s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().Nil(chainlet.Upgrade)                         // not upgrading anymore
				s.Require().Equal("2.0.0", chainlet.ChainletStackVersion) // new version set
			},
		}, {
			name: "incorrect client ID",
			fn: func(chainID, consumerID, clientID, connectionID, channelID string) {
				s.breakingUpgrade(chainID, consumerID, clientID, connectionID, channelID)

				packet := channeltypes.Packet{
					DestinationChannel: "channel-42",
				}
				data := chainlettypes.ConfirmUpgradePacketData{
					ChainId: chainID,
					Height:  123,
					Plan:    "xyz",
				}

				// Upgrade confirmation
				s.packetVerificationMocks(consumerID, clientID, "bad-client", "bad-connection", "channel-42")
				_, err := s.chainletKeeper.OnRecvConfirmUpgradePacket(s.ctx, packet, data)
				s.Require().Error(err)

				// Check if the upgrade is not removed
				chainlet, err := s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)
			},
		}, {
			name: "incorrect chain ID",
			fn: func(chainID, consumerID, clientID, connectionID, channelID string) {
				s.breakingUpgrade(chainID, consumerID, clientID, connectionID, channelID)

				packet := channeltypes.Packet{
					DestinationChannel: channelID,
				}
				data := chainlettypes.ConfirmUpgradePacketData{
					ChainId: "abcd", // incorrect
					Height:  123,
					Plan:    "xyz",
				}
				_, err := s.chainletKeeper.OnRecvConfirmUpgradePacket(s.ctx, packet, data)
				s.Require().Error(err)

				// Check if the upgrade is not removed
				chainlet, err := s.chainletKeeper.Chainlet(s.ctx, chainID)
				s.Require().NoError(err)
				s.Require().NotNil(chainlet.Upgrade)
			},
		},
	}
	for i, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()

			chainID := fmt.Sprintf("chain_%d-1", i+1)
			consumerID := fmt.Sprintf("%d", i)
			clientID := fmt.Sprintf("client-%d", i)
			connectionID := fmt.Sprintf("connection-%d", i)
			channelID := fmt.Sprintf("channel-%d", i)

			s.ibcSetup(chainID, consumerID, channelID)

			tt.fn(chainID, consumerID, clientID, connectionID, channelID)
		})
	}
}
