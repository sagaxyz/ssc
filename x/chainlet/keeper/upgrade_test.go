package keeper_test

import (
	"fmt"

	ibcclienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v10/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ccvprovidertypes "github.com/cosmos/interchain-security/v7/x/ccv/provider/types"
	"github.com/golang/mock/gomock"
	sdkchainlettypes "github.com/sagaxyz/saga-sdk/x/chainlet/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (s *TestSuite) TestUpgrade() {
	tests := []struct {
		name   string
		expErr bool
		fn     func(chainID, consumerID, clientID, connectionID, channelID string) error
	}{
		{
			name:   "ok",
			expErr: false,
			fn: func(chainID, consumerID, clientID, connectionID, channelID string) error {
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
							gomock.Any(), // timeout height
							gomock.Any(), // timeout timestamp
							gomock.Any(), // data
							//TODO check any values
						).
						Return(uint64(1337), nil), //TODO
				)
				return nil
			},
		}, {
			name:   "consumer not registered yet",
			expErr: true,
			fn: func(chainID, consumerID, clientID, connectionID, channelID string) error {
				gomock.InOrder(
					s.providerKeeper.EXPECT().
						GetConsumerClientId(gomock.Any(), gomock.Eq(consumerID)).
						Return("", false),
				)
				return nil
			},
		}, {
			name:   "incorrect client id for the provided channel",
			expErr: true,
			fn: func(chainID, consumerID, clientID, connectionID, channelID string) error {
				gomock.InOrder(
					s.providerKeeper.EXPECT().
						GetConsumerClientId(gomock.Any(), gomock.Eq(consumerID)).
						Return("client-123", true),
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
				return nil
			},
		},
	}
	for i, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()

			// Calls we don't care about in this test
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
			s.aclKeeper.EXPECT().
				IsAdmin(gomock.Any(), gomock.Any()).
				Return(false).
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
			chainID := fmt.Sprintf("chain_%d-1", i+1)
			consumerID := fmt.Sprintf("%d", i)
			clientID := fmt.Sprintf("client-%d", i)
			connectionID := fmt.Sprintf("connection-%d", i)
			channelID := fmt.Sprintf("channel-%d", i)

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

			_ = tt.fn(chainID, consumerID, clientID, connectionID, channelID) //TODO remove return value

			// Launch a chainlet
			_, err = s.msgServer.LaunchChainlet(s.ctx, types.NewMsgLaunchChainlet(
				creator.String(), []string{creator.String()}, "test", ver, "test_chainlet", chainID, "asaga", types.ChainletParams{}, nil, false, "",
			))
			s.Require().NoError(err)
			s.chainletKeeper.InitConsumers(s.ctx)

			// Breaking upgrade
			resp, err := s.msgServer.UpgradeChainlet(s.ctx, &types.MsgUpgradeChainlet{
				Creator:      creator.String(),
				ChainId:      chainID,
				StackVersion: "2.0.0",
				HeightDelta:  100,
				ChannelId:    channelID,
			})
			if tt.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(uint64(0xc8), resp.Height) //TODO calculate correct value
			}

			// Check if upgrade is correctly set/unset in the chainlet
			chainlet, err := s.chainletKeeper.Chainlet(s.ctx, chainID)
			s.Require().NoError(err)
			if tt.expErr {
				s.Require().Nil(chainlet.Upgrade)
			} else {
				s.Require().NotNil(chainlet.Upgrade)
				s.Require().Equal("2.0.0", chainlet.Upgrade.Version)
			}
		})
	}
}
