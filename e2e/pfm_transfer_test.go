package e2e

import (
	"context"
	"encoding/json"
	"testing"

	"cosmossdk.io/math"
	pfmtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	e2eutils "github.com/sagaxyz/ssc/e2e/utils"
	interchaintest "github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/stretchr/testify/require"
)

// TestPFMTransfer tests a packet-forward-middleware transfer across 3 chains using 2 relayers
func TestPFMTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()
	ctx := context.Background()

	// NOTE: we want to connect chains A-B and B-C here
	pathAB := e2eutils.RelayerPath{0, 1}
	pathBC := e2eutils.RelayerPath{1, 2}

	icn, err := e2eutils.CreateAndStartFullyConnectedNetwork(t, ctx,
		e2eutils.WithNChains(3),
		e2eutils.WithRelayerPaths(pathAB, pathBC),
	)
	require.NoError(t, err)

	chainA, err := icn.GetChain(0)
	require.NoError(t, err)
	chainB, err := icn.GetChain(1)
	require.NoError(t, err)
	chainC, err := icn.GetChain(2)
	require.NoError(t, err)

	// Fund users on all chains
	fundAmount := math.NewInt(10_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", fundAmount, chainA, chainB, chainC)
	userA := users[0]
	userB := users[1]
	userC := users[2]

	// Wait for a few blocks
	err = testutil.WaitForBlocks(ctx, 3, chainA, chainB, chainC)
	require.NoError(t, err)

	// Verify initial balances
	balA, err := chainA.GetBalance(ctx, userA.FormattedAddress(), chainA.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, fundAmount.String(), balA.String(), "expected different initial balance for chain A")

	balB, err := chainB.GetBalance(ctx, userB.FormattedAddress(), chainB.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, fundAmount.String(), balB.String(), "expected different initial balance for chain B")

	balC, err := chainC.GetBalance(ctx, userC.FormattedAddress(), chainC.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, fundAmount.String(), balC.String(), "expected different initial balance for chain C")

	// Retrieve relayer channels
	channelAB, err := icn.GetChannelInfo(ctx, pathAB)
	require.NoError(t, err, "failed to get channel AB")
	channelBC, err := icn.GetChannelInfo(ctx, pathBC)
	require.NoError(t, err, "failed to get channel BC")

	// Record chain heights before transfer
	chainAHeight, err := chainA.Height(ctx)
	require.NoError(t, err)

	// Prepare PFM transfer from A->C through B
	amountToSend := math.NewInt(1_000_000)
	transfer := ibc.WalletAmount{
		Address: "pfm",
		Denom:   chainA.Config().Denom,
		Amount:  amountToSend,
	}

	firstHopMetadata := &pfmtypes.PacketMetadata{
		Forward: &pfmtypes.ForwardMetadata{
			Receiver: userC.FormattedAddress(),
			Channel:  channelBC.ChannelID,
			Port:     channelBC.PortID,
		},
	}

	memo, err := json.Marshal(firstHopMetadata)
	require.NoError(t, err)

	// Execute PFM transfer
	tx, err := chainA.SendIBCTransfer(ctx, channelAB.ChannelID, userA.KeyName(), transfer, ibc.TransferOptions{
		Memo: string(memo),
	})
	require.NoError(t, err)
	require.NoError(t, tx.Validate())

	// Wait for packet processing
	_, err = testutil.PollForAck(ctx, chainA, chainAHeight, chainAHeight+50, tx.Packet)
	require.NoError(t, err)

	// NOTE: for now we're waiting a bunch of blocks here to account for the PFM forwarding from chain B to C.
	//
	// TODO: we should rather check for ack on chain B that the PFM transfer was completed.
	err = testutil.WaitForBlocks(ctx, 30, chainA, chainB, chainC)
	require.NoError(t, err)

	// Verify final balances
	// Chain A: Initial balance - sent amount
	newBalA, err := chainA.GetBalance(ctx, userA.FormattedAddress(), chainA.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, balA.Sub(amountToSend).String(), newBalA.String(), "expected different balance for chain A after PFM transfer")

	// Chain C: Should have received the forwarded amount with appropriate IBC denom
	firstHopDenom := transfertypes.GetPrefixedDenom(channelAB.Counterparty.PortID, channelAB.Counterparty.ChannelID, chainA.Config().Denom)
	secondHopDenom := transfertypes.GetPrefixedDenom(channelBC.Counterparty.PortID, channelBC.Counterparty.ChannelID, firstHopDenom)
	firstHopIBCDenom := transfertypes.ParseDenomTrace(firstHopDenom).IBCDenom()
	secondHopIBCDenom := transfertypes.ParseDenomTrace(secondHopDenom).IBCDenom()

	newBalB, err := chainB.GetBalance(ctx, userB.FormattedAddress(), firstHopIBCDenom)
	require.NoError(t, err)
	require.Equal(t, "0", newBalB.String(), "expected zero balance for chain B after PFM transfer")

	newBalC, err := chainC.GetBalance(ctx, userC.FormattedAddress(), secondHopIBCDenom)
	require.NoError(t, err)
	require.Equal(t, amountToSend.String(), newBalC.String(), "expected different balance for chain C after PFM transfer")
}
