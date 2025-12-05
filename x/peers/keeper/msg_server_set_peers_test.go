package keeper_test

import (
	"errors"
	"strings"

	"github.com/golang/mock/gomock"
	chainlettypes "github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/sagaxyz/ssc/x/peers/types"
)

func (s *TestSuite) TestChainIDValidation() {
	s.SetupTest()

	// Chainlet exists
	s.chainletKeeper.EXPECT().Chainlet(gomock.Any(), chainIDs[0]).Return(chainlettypes.Chainlet{ChainId: chainIDs[0]}, nil)
	_, err := s.msgServer.SetPeers(s.ctx, types.NewMsgSetPeers(accounts[0].String(), chainIDs[0], addrs[chainIDs[0]]...))
	s.Require().NoError(err)

	// Chainlet does not exists
	s.chainletKeeper.EXPECT().Chainlet(gomock.Any(), "something").Return(chainlettypes.Chainlet{}, errors.New("nope"))
	_, err = s.msgServer.SetPeers(s.ctx, types.NewMsgSetPeers(accounts[0].String(), "something", addrs[chainIDs[0]]...))
	s.Require().Error(err)
}
func (s *TestSuite) TestPeersValidation() {
	tests := []struct {
		peers  []string
		expErr bool
	}{
		{addrs[chainIDs[0]], false},
		{addrs[chainIDs[1]], false},
		{addrs[chainIDs[2]], false},
		{[]string{"aa@127.0.0.1:1234"}, false},
		{[]string{"aa@example.com:1234"}, false},
		{[]string{"abcd"}, true},
		{[]string{""}, true},
		{[]string{}, true},
		{[]string{"127.0.0.1:1234"}, true},
		{[]string{"@"}, true},
		{[]string{"aa@b"}, true},
		{[]string{"aa@"}, true},
		{[]string{"@127.0.0.1"}, true},
		{[]string{"@127.0.0.1:1234"}, true},     // missing ID
		{[]string{"aa@127.0.0.1:y"}, true},      // invalid port
		{[]string{"a@127.0.0.1:1234"}, true},    // invalid hex ID
		{[]string{"xx@127.0.0.1:1234"}, true},   // invalid hex ID
		{[]string{"aa'@127.0.0.1:1234"}, true},  // invalid character
		{[]string{"aa@'127.0.0.1:1234"}, true},  // invalid character
		{[]string{"aa@127.0.0.1:1234'"}, true},  // invalid character
		{[]string{"aa@127.0.0.1:1234 "}, true},  // invalid character
		{[]string{"aa\"@127.0.0.1:1234"}, true}, // invalid character
		{[]string{"aa@\"127.0.0.1:1234"}, true}, // invalid character
		{[]string{"aa@127.0.0.1:1234\""}, true}, // invalid character
		// Test size limit
		{[]string{strings.Repeat("aa", 1000) + "@example.com:1234"}, true},
		{[]string{strings.Repeat("aa", 300) + "@example.com:1234"}, false},
		{[]string{strings.Repeat("bb", 300) + "@example2.com:1234"}, false},
		{[]string{
			strings.Repeat("aa", 300) + "@example.com:1234",
			strings.Repeat("bb", 300) + "@example2.com:1234",
		}, true},
	}

	for _, tt := range tests {
		s.Run(strings.Join(tt.peers, ","), func() {
			s.SetupTest()

			s.chainletKeeper.EXPECT().Chainlet(gomock.Any(), chainIDs[0]).Return(chainlettypes.Chainlet{ChainId: chainIDs[0]}, nil).AnyTimes()

			_, err := s.msgServer.SetPeers(s.ctx, types.NewMsgSetPeers(accounts[0].String(), chainIDs[0], tt.peers...))
			if tt.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
