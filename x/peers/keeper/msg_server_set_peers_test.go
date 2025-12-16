package keeper_test

import (
	"errors"
	"math"
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
		name   string
		peers  []string
		expErr string
	}{
		{
			"ok",
			addrs[chainIDs[0]],
			"",
		},
		{
			"ok",
			addrs[chainIDs[1]],
			"",
		},
		{
			"ok",
			addrs[chainIDs[2]],
			"",
		},
		{
			"ok",
			[]string{"aa@127.0.0.1:1234"},
			"",
		},
		{
			"ok",
			[]string{"aa@example.com:1234"},
			"",
		},
		{
			"missing ID and port",
			[]string{"abcd"},
			"invalid addr",
		},
		{
			"empty string",
			[]string{""},
			"invalid addr",
		},
		{
			"empty set",
			[]string{},
			"no peers provided",
		},
		{
			"missing ID",
			[]string{"127.0.0.1:1234"},
			"invalid addr",
		},
		{
			"missing ID and host:port",
			[]string{"@"},
			"invalid addr",
		},
		{
			"missing port",
			[]string{"aa@b"},
			"invalid addr",
		},
		{
			"missing host:port",
			[]string{"aa@"},
			"invalid addr",
		},
		{
			"missing ID",
			[]string{"@127.0.0.1"},
			"invalid addr",
		},
		{
			"missing ID",
			[]string{"@127.0.0.1:1234"},
			"invalid addr",
		},
		{
			"invalid port",
			[]string{"aa@127.0.0.1:y"},
			"invalid addr",
		},
		{
			"invalid hex ID",
			[]string{"a@127.0.0.1:1234"},
			"invalid addr",
		},
		{
			"invalid hex ID",
			[]string{"xx@127.0.0.1:1234"},
			"invalid addr",
		},
		{
			"invalid character",
			[]string{"aa'@127.0.0.1:1234"},
			"invalid addr",
		},
		{
			"invalid character",
			[]string{"aa@'127.0.0.1:1234"},
			"invalid addr",
		},
		{
			"invalid character",
			[]string{"aa@127.0.0.1:1234'"},
			"invalid addr",
		},
		{
			"invalid character",
			[]string{"aa@127.0.0.1:1234 "},
			"invalid addr",
		},
		{
			"invalid character",
			[]string{"aa\"@127.0.0.1:1234"},
			"invalid addr",
		},
		{
			"invalid character",
			[]string{"aa@\"127.0.0.1:1234"},
			"invalid addr",
		},
		{
			"invalid character",
			[]string{"aa@127.0.0.1:1234\""},
			"invalid addr",
		},
		// Test size limit
		{
			"single entry size",
			[]string{strings.Repeat("aa",
				types.DefaultMaxData) + "@example.com:1234"},
			"exceeded maximum size"},
		{
			"verify single entry works for the maximum data size test",
			[]string{strings.Repeat("aa",
				types.DefaultMaxData/3) + "@example.com:1234"},
			""},
		{
			"verify single entry works for the maximum data size test",
			[]string{strings.Repeat("bb",
				types.DefaultMaxData/3) + "@example2.com:1234"},
			""},
		{
			"maximum data size",
			[]string{
				strings.Repeat("aa", types.DefaultMaxData/3) + "@example.com:1234",
				strings.Repeat("bb", types.DefaultMaxData/3) + "@example2.com:1234",
			},
			"exceeded maximum size",
		},
		// Test overflow
		{
			"single addr size overflow",
			[]string{strings.Repeat("a", math.MaxUint32) + "@example.com:1234"},
			"addr size exceeds uint32",
		},
		{
			"total data size overflow",
			[]string{
				// First address is small enough to be allowed, but makes the second address overflow the total data counter
				// without it failing the single entry overflow check.
				strings.Repeat("a", types.DefaultMaxData/2) + "@example.com:1234",
				strings.Repeat("a", math.MaxUint32-types.DefaultMaxData/2+1) + "@example.com:1234",
			}, "total size exceeds uint32",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()

			s.chainletKeeper.EXPECT().Chainlet(gomock.Any(), chainIDs[0]).Return(chainlettypes.Chainlet{ChainId: chainIDs[0]}, nil).AnyTimes()

			_, err := s.msgServer.SetPeers(s.ctx, types.NewMsgSetPeers(accounts[0].String(), chainIDs[0], tt.peers...))
			if tt.expErr != "" {
				s.Require().Error(err)
				s.Require().True(strings.Contains(err.Error(), tt.expErr), "error '%s' does not contain '%s'", err.Error(), tt.expErr)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
