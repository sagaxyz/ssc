package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"

	chainlettypes "github.com/sagaxyz/ssc/x/chainlet/types"
	types "github.com/sagaxyz/ssc/x/peers/types"
)

func (s *TestSuite) TestAfterValidatorRemoved() {
	require := s.Require()

	_, _, addr := testdata.KeyTestPubAddr()
	valAddr := sdk.ValAddress(addr)
	accAddr := sdk.AccAddress(addr)
	consAddr := sdk.ConsAddress(addr)

	s.chainletKeeper.EXPECT().Chainlet(gomock.Any(), chainIDs[0]).Return(chainlettypes.Chainlet{ChainId: chainIDs[0]}, nil).AnyTimes()
	s.chainletKeeper.EXPECT().Chainlet(gomock.Any(), chainIDs[1]).Return(chainlettypes.Chainlet{ChainId: chainIDs[1]}, nil).AnyTimes()

	_, err := s.msgServer.SetPeers(s.ctx, &types.MsgSetPeers{
		Validator: accAddr.String(),
		ChainId:   chainIDs[0],
		Peers:     addrs[chainIDs[0]],
	})
	require.NoError(err)
	_, err = s.msgServer.SetPeers(s.ctx, &types.MsgSetPeers{
		Validator: accAddr.String(),
		ChainId:   chainIDs[1],
		Peers:     addrs[chainIDs[1]],
	})
	require.NoError(err)

	resp, err := s.queryClient.Peers(s.ctx, &types.QueryPeersRequest{
		ChainId: chainIDs[0],
	})
	require.NoError(err)
	require.Equal(addrs[chainIDs[0]], resp.Peers)

	err = s.peersKeeper.Hooks().AfterValidatorRemoved(s.ctx, consAddr, valAddr)
	require.NoError(err)

	resp, err = s.queryClient.Peers(s.ctx, &types.QueryPeersRequest{
		ChainId: chainIDs[0],
	})
	require.NoError(err)
	require.Equal([]string(nil), resp.Peers)
	resp, err = s.queryClient.Peers(s.ctx, &types.QueryPeersRequest{
		ChainId: chainIDs[1],
	})
	require.NoError(err)
	require.Equal([]string(nil), resp.Peers)
}
