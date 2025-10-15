package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	types "github.com/sagaxyz/ssc/x/peers/types"
)

func (s *KeeperTestSuite) TestDataStorage() {
	_, _, addrA := testdata.KeyTestPubAddr()
	valAddrA := sdk.ValAddress(addrA)
	_, _, addrB := testdata.KeyTestPubAddr()
	valAddrB := sdk.ValAddress(addrB)

	// Sanity check
	s.Require().Equal([]string{}, s.peersKeeper.GetPeers(s.ctx, chainIDs[0]))
	s.Require().Equal(uint32(0), s.peersKeeper.Counter(s.ctx, chainIDs[0]))
	s.Require().Equal([]string{}, s.peersKeeper.GetPeers(s.ctx, chainIDs[1]))
	s.Require().Equal(uint32(0), s.peersKeeper.Counter(s.ctx, chainIDs[1]))

	// Add some other chain peers
	s.peersKeeper.StoreData(s.ctx, chainIDs[1], valAddrA.String(), types.Data{
		Updated:   s.ctx.BlockTime(),
		Addresses: []string{"x", "y"},
	})
	s.peersKeeper.StoreData(s.ctx, chainIDs[1], valAddrB.String(), types.Data{
		Updated:   s.ctx.BlockTime(),
		Addresses: []string{"X", "Y"},
	})
	s.Require().Equal(4, len(s.peersKeeper.GetPeers(s.ctx, chainIDs[1])))
	s.Require().Equal(uint32(2), s.peersKeeper.Counter(s.ctx, chainIDs[1]))

	// Add val A addrs
	s.peersKeeper.StoreData(s.ctx, chainIDs[0], valAddrA.String(), types.Data{
		Updated:   s.ctx.BlockTime(),
		Addresses: []string{"a", "b"},
	})
	s.Require().Equal([]string{"a", "b"}, s.peersKeeper.GetPeers(s.ctx, chainIDs[0]))
	s.Require().Equal(uint32(1), s.peersKeeper.Counter(s.ctx, chainIDs[0]))

	// Add val B addrs
	s.peersKeeper.StoreData(s.ctx, chainIDs[0], valAddrB.String(), types.Data{
		Updated:   s.ctx.BlockTime(),
		Addresses: []string{"c", "d"},
	})
	s.Require().Equal(4, len(s.peersKeeper.GetPeers(s.ctx, chainIDs[0])))
	s.Require().Equal(uint32(2), s.peersKeeper.Counter(s.ctx, chainIDs[0]))

	// Remove val A addrs
	s.peersKeeper.DeleteValidatorData(s.ctx, valAddrA.String())
	s.Require().Equal([]string{"c", "d"}, s.peersKeeper.GetPeers(s.ctx, chainIDs[0]))
	s.Require().Equal(uint32(1), s.peersKeeper.Counter(s.ctx, chainIDs[0]))

	// Remove val B addrs
	s.peersKeeper.DeleteValidatorData(s.ctx, valAddrB.String())
	s.Require().Equal([]string{}, s.peersKeeper.GetPeers(s.ctx, chainIDs[0]))
	s.Require().Equal(uint32(0), s.peersKeeper.Counter(s.ctx, chainIDs[0]))

	// We removed all addresses for both validators
	s.Require().Equal(0, len(s.peersKeeper.GetPeers(s.ctx, chainIDs[1])))
	s.Require().Equal(uint32(0), s.peersKeeper.Counter(s.ctx, chainIDs[1]))
}
