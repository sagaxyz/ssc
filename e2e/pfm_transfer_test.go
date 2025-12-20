package e2e

import (
	"context"
	"encoding/json"
	"testing"

	"cosmossdk.io/math"
	pfmtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v10/packetforward/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	interchaintest "github.com/cosmos/interchaintest/v10"
	"github.com/cosmos/interchaintest/v10/ibc"
	"github.com/cosmos/interchaintest/v10/testutil"
	e2eutils "github.com/sagaxyz/ssc/e2e/utils"
	"github.com/stretchr/testify/require"
)

// TestPFMTransfer tests a packet-forward-middleware transfer across 3 chains using 2 relayers
func TestPFMTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	t.Log("🧪 Starting PFM Transfer Test")
	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	t.Parallel()
	ctx := context.Background()

	// NOTE: we want to connect chains A-B and B-C here
	pathAB := e2eutils.RelayerPath{0, 1}
	pathBC := e2eutils.RelayerPath{1, 2}

	t.Log("📡 Step 1: Creating and starting 3-chain network with relayers")
	t.Logf("   - Relayer path A→B: %v", pathAB)
	t.Logf("   - Relayer path B→C: %v", pathBC)

	icn, err := e2eutils.CreateAndStartFullyConnectedNetwork(t, ctx,
		e2eutils.WithNChains(3),
		e2eutils.WithRelayerPaths(pathAB, pathBC),
	)
	if err != nil {
		t.Logf("❌ Failed to create network: %v", err)
	}
	require.NoError(t, err)
	t.Log("✅ Network created successfully")

	chainA, err := icn.GetChain(0)
	require.NoError(t, err)
	chainB, err := icn.GetChain(1)
	require.NoError(t, err)
	chainC, err := icn.GetChain(2)
	require.NoError(t, err)

	t.Logf("   - Chain A: %s (denom: %s)", chainA.Config().Name, chainA.Config().Denom)
	t.Logf("   - Chain B: %s (denom: %s)", chainB.Config().Name, chainB.Config().Denom)
	t.Logf("   - Chain C: %s (denom: %s)", chainC.Config().Name, chainC.Config().Denom)

	// Fund users on all chains
	t.Log("")
	t.Log("💰 Step 2: Funding test users on all chains")
	fundAmount := math.NewInt(10_000_000)
	t.Logf("   - Funding amount: %s", fundAmount.String())

	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", fundAmount, chainA, chainB, chainC)
	userA := users[0]
	userB := users[1]
	userC := users[2]

	t.Logf("   - User A (Chain A): %s", userA.FormattedAddress())
	t.Logf("   - User B (Chain B): %s", userB.FormattedAddress())
	t.Logf("   - User C (Chain C): %s", userC.FormattedAddress())
	t.Log("✅ Users funded")

	// Wait for a few blocks
	t.Log("")
	t.Log("⏳ Step 3: Waiting for blocks to be produced")
	err = testutil.WaitForBlocks(ctx, 3, chainA, chainB, chainC)
	if err != nil {
		t.Logf("❌ Failed to wait for blocks: %v", err)
	}
	require.NoError(t, err)
	t.Log("✅ Blocks produced")

	// Verify initial balances
	t.Log("")
	t.Log("🔍 Step 4: Verifying initial balances")
	balA, err := chainA.GetBalance(ctx, userA.FormattedAddress(), chainA.Config().Denom)
	require.NoError(t, err)
	if balA.String() == fundAmount.String() {
		t.Logf("✅ Chain A balance verified: %s (expected: %s)", balA.String(), fundAmount.String())
	} else {
		t.Logf("❌ Chain A balance mismatch: got %s, expected %s", balA.String(), fundAmount.String())
	}
	require.Equal(t, fundAmount.String(), balA.String(), "expected different initial balance for chain A")

	balB, err := chainB.GetBalance(ctx, userB.FormattedAddress(), chainB.Config().Denom)
	require.NoError(t, err)
	if balB.String() == fundAmount.String() {
		t.Logf("✅ Chain B balance verified: %s (expected: %s)", balB.String(), fundAmount.String())
	} else {
		t.Logf("❌ Chain B balance mismatch: got %s, expected %s", balB.String(), fundAmount.String())
	}
	require.Equal(t, fundAmount.String(), balB.String(), "expected different initial balance for chain B")

	balC, err := chainC.GetBalance(ctx, userC.FormattedAddress(), chainC.Config().Denom)
	require.NoError(t, err)
	if balC.String() == fundAmount.String() {
		t.Logf("✅ Chain C balance verified: %s (expected: %s)", balC.String(), fundAmount.String())
	} else {
		t.Logf("❌ Chain C balance mismatch: got %s, expected %s", balC.String(), fundAmount.String())
	}
	require.Equal(t, fundAmount.String(), balC.String(), "expected different initial balance for chain C")

	// Retrieve relayer channels
	t.Log("")
	t.Log("🔗 Step 5: Retrieving IBC channel information")
	channelAB, err := icn.GetChannelInfo(ctx, pathAB)
	if err != nil {
		t.Logf("❌ Failed to get channel AB: %v", err)
	}
	require.NoError(t, err, "failed to get channel AB")
	t.Logf("   - Channel A→B: %s (port: %s)", channelAB.ChannelID, channelAB.PortID)

	channelBC, err := icn.GetChannelInfo(ctx, pathBC)
	if err != nil {
		t.Logf("❌ Failed to get channel BC: %v", err)
	}
	require.NoError(t, err, "failed to get channel BC")
	t.Logf("   - Channel B→C: %s (port: %s)", channelBC.ChannelID, channelBC.PortID)
	t.Log("✅ Channels retrieved")

	// Record chain heights before transfer
	chainAHeight, err := chainA.Height(ctx)
	require.NoError(t, err)
	t.Logf("   - Chain A height before transfer: %d", chainAHeight)

	// Prepare PFM transfer from A->C through B
	t.Log("")
	t.Log("📦 Step 6: Preparing PFM transfer (A→B→C)")
	amountToSend := math.NewInt(1_000_000)
	t.Logf("   - Amount to send: %s", amountToSend.String())
	t.Logf("   - Source: Chain A, User A (%s)", userA.FormattedAddress())
	t.Logf("   - Destination: Chain C, User C (%s)", userC.FormattedAddress())
	t.Logf("   - Route: Chain A → Chain B → Chain C")

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
	t.Logf("   - PFM memo: %s", string(memo))
	t.Log("✅ Transfer prepared")

	// Execute PFM transfer
	t.Log("")
	t.Log("🚀 Step 7: Executing PFM transfer")
	t.Logf("   - Sending IBC transfer on channel %s", channelAB.ChannelID)

	tx, err := chainA.SendIBCTransfer(ctx, channelAB.ChannelID, userA.KeyName(), transfer, ibc.TransferOptions{
		Memo: string(memo),
	})
	if err != nil {
		t.Logf("❌ Failed to send IBC transfer: %v", err)
	}
	require.NoError(t, err)
	t.Logf("   - Transaction hash: %s", tx.TxHash)

	err = tx.Validate()
	if err != nil {
		t.Logf("❌ Transaction validation failed: %v", err)
	}
	require.NoError(t, err)
	t.Logf("   - Packet: Sequence %d, Source Port: %s, Source Channel: %s",
		tx.Packet.Sequence, tx.Packet.SourcePort, tx.Packet.SourceChannel)
	t.Log("✅ Transfer transaction submitted and validated")

	// Wait for packet processing
	t.Log("")
	t.Log("⏳ Step 8: Waiting for packet acknowledgment")
	t.Logf("   - Polling for ACK from height %d to %d", chainAHeight, chainAHeight+50)

	ack, err := testutil.PollForAck(ctx, chainA, chainAHeight, chainAHeight+50, tx.Packet)
	if err != nil {
		t.Logf("❌ Failed to receive ACK: %v", err)
	} else {
		t.Log("✅ Packet ACK received")
		if len(ack.Acknowledgement) > 0 {
			t.Logf("   - ACK: %s", string(ack.Acknowledgement))
		}
	}
	require.NoError(t, err)

	// NOTE: for now we're waiting a bunch of blocks here to account for the PFM forwarding from chain B to C.
	//
	// TODO: we should rather check for ack on chain B that the PFM transfer was completed.
	t.Log("")
	t.Log("⏳ Step 9: Waiting for PFM forwarding to complete (B→C)")
	t.Log("   - Waiting 30 blocks for multi-hop forwarding")

	err = testutil.WaitForBlocks(ctx, 30, chainA, chainB, chainC)
	if err != nil {
		t.Logf("❌ Failed to wait for blocks: %v", err)
	}
	require.NoError(t, err)
	t.Log("✅ Blocks produced, forwarding should be complete")

	// Verify final balances
	t.Log("")
	t.Log("🔍 Step 10: Verifying final balances")

	// Chain A: Initial balance - sent amount
	expectedBalA := balA.Sub(amountToSend)
	newBalA, err := chainA.GetBalance(ctx, userA.FormattedAddress(), chainA.Config().Denom)
	require.NoError(t, err)
	if newBalA.String() == expectedBalA.String() {
		t.Logf("✅ Chain A balance correct: %s (expected: %s, sent: %s)",
			newBalA.String(), expectedBalA.String(), amountToSend.String())
	} else {
		t.Logf("❌ Chain A balance incorrect: got %s, expected %s (sent %s)",
			newBalA.String(), expectedBalA.String(), amountToSend.String())
	}
	require.Equal(t, balA.Sub(amountToSend).String(), newBalA.String(), "expected different balance for chain A after PFM transfer")

	// Chain C: Should have received the forwarded amount with appropriate IBC denom
	firstHopDenom := transfertypes.GetPrefixedDenom(channelAB.Counterparty.PortID, channelAB.Counterparty.ChannelID, chainA.Config().Denom)
	secondHopDenom := transfertypes.GetPrefixedDenom(channelBC.Counterparty.PortID, channelBC.Counterparty.ChannelID, firstHopDenom)
	firstHopIBCDenom := transfertypes.ParseDenomTrace(firstHopDenom).IBCDenom()
	secondHopIBCDenom := transfertypes.ParseDenomTrace(secondHopDenom).IBCDenom()

	t.Logf("   - First hop IBC denom: %s", firstHopIBCDenom)
	t.Logf("   - Second hop IBC denom: %s", secondHopIBCDenom)

	// Chain B: Should have zero balance (PFM should forward immediately)
	newBalB, err := chainB.GetBalance(ctx, userB.FormattedAddress(), firstHopIBCDenom)
	require.NoError(t, err)
	if newBalB.String() == "0" {
		t.Logf("✅ Chain B balance correct: %s (expected: 0, PFM forwarded immediately)",
			newBalB.String())
	} else {
		t.Logf("❌ Chain B balance incorrect: got %s, expected 0 (PFM should have forwarded)",
			newBalB.String())
	}
	require.Equal(t, "0", newBalB.String(), "expected zero balance for chain B after PFM transfer")

	// Chain C: Final recipient should have the amount
	newBalC, err := chainC.GetBalance(ctx, userC.FormattedAddress(), secondHopIBCDenom)
	require.NoError(t, err)
	if newBalC.String() == amountToSend.String() {
		t.Logf("✅ Chain C balance correct: %s (expected: %s, received via PFM)",
			newBalC.String(), amountToSend.String())
	} else {
		t.Logf("❌ Chain C balance incorrect: got %s, expected %s",
			newBalC.String(), amountToSend.String())
	}
	require.Equal(t, amountToSend.String(), newBalC.String(), "expected different balance for chain C after PFM transfer")

	t.Log("")
	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	t.Log("✅ PFM Transfer Test PASSED")
	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	t.Logf("   Summary:")
	t.Logf("   - Sent %s from Chain A to Chain C via Chain B", amountToSend.String())
	t.Logf("   - Chain A balance: %s → %s", balA.String(), newBalA.String())
	t.Logf("   - Chain B balance: 0 (forwarded immediately)")
	t.Logf("   - Chain C balance: 0 → %s", newBalC.String())
}
