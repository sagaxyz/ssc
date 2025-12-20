package e2e

import (
	"context"
	"testing"

	"cosmossdk.io/math"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	interchaintest "github.com/cosmos/interchaintest/v10"
	"github.com/cosmos/interchaintest/v10/ibc"
	"github.com/cosmos/interchaintest/v10/testutil"
	e2eutils "github.com/sagaxyz/ssc/e2e/utils"
	"github.com/stretchr/testify/require"
)

// TestBasicIBCTransfer is a basic test to start 2 SSC chains and send an IBC transfer between them.
func TestBasicIBCTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	t.Log("🧪 Starting Basic IBC Transfer Test")
	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	t.Parallel()
	ctx := context.Background()

	t.Log("📡 Step 1: Creating and starting 2-chain network with relayer")
	pathAB := e2eutils.RelayerPath{0, 1}
	t.Logf("   - Relayer path: %v", pathAB)

	icn, err := e2eutils.CreateAndStartFullyConnectedNetwork(t, ctx,
		e2eutils.WithNChains(2),
		e2eutils.WithRelayerPaths(pathAB),
	)
	if err != nil {
		t.Logf("❌ Failed to create network: %v", err)
	}
	require.NoError(t, err)
	require.NotNil(t, icn)
	t.Log("✅ Network created successfully")

	chainFrom, err := icn.GetChain(0)
	require.NoError(t, err)
	chainTo, err := icn.GetChain(1)
	require.NoError(t, err)

	t.Logf("   - Chain From: %s (denom: %s)", chainFrom.Config().Name, chainFrom.Config().Denom)
	t.Logf("   - Chain To: %s (denom: %s)", chainTo.Config().Name, chainTo.Config().Denom)

	t.Log("")
	t.Log("💰 Step 2: Funding test users")
	fundAmount := math.NewInt(10_000_000)
	t.Logf("   - Funding amount: %s", fundAmount.String())

	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", fundAmount, chainFrom, chainTo)
	userA := users[0]
	userB := users[1]

	t.Logf("   - User A (Chain From): %s", userA.FormattedAddress())
	t.Logf("   - User B (Chain To): %s", userB.FormattedAddress())
	t.Log("✅ Users funded")

	t.Log("")
	t.Log("⏳ Step 3: Waiting for blocks to be produced")
	err = testutil.WaitForBlocks(ctx, 3, chainFrom, chainTo)
	if err != nil {
		t.Logf("❌ Failed to wait for blocks: %v", err)
	}
	require.NoError(t, err)
	t.Log("✅ Blocks produced")

	t.Log("")
	t.Log("🔍 Step 4: Verifying initial balances")
	balA, err := chainFrom.GetBalance(ctx, userA.FormattedAddress(), chainFrom.Config().Denom)
	require.NoError(t, err)
	if balA.Equal(fundAmount) {
		t.Logf("✅ Chain From balance verified: %s (expected: %s)", balA.String(), fundAmount.String())
	} else {
		t.Logf("❌ Chain From balance mismatch: got %s, expected %s", balA.String(), fundAmount.String())
	}
	require.True(t, balA.Equal(fundAmount))

	balB, err := chainTo.GetBalance(ctx, userB.FormattedAddress(), chainTo.Config().Denom)
	require.NoError(t, err)
	if balB.Equal(fundAmount) {
		t.Logf("✅ Chain To balance verified: %s (expected: %s)", balB.String(), fundAmount.String())
	} else {
		t.Logf("❌ Chain To balance mismatch: got %s, expected %s", balB.String(), fundAmount.String())
	}
	require.True(t, balB.Equal(fundAmount))

	t.Log("")
	t.Log("🔗 Step 5: Retrieving IBC channel information")
	channel, err := icn.GetChannelInfo(ctx, pathAB)
	if err != nil {
		t.Logf("❌ Failed to get channel: %v", err)
	}
	require.NoError(t, err)
	t.Logf("   - Channel ID: %s", channel.ChannelID)
	t.Logf("   - Port ID: %s", channel.PortID)
	t.Log("✅ Channel retrieved")

	chainFromHeight, err := chainFrom.Height(ctx)
	require.NoError(t, err)
	t.Logf("   - Chain From height before transfer: %d", chainFromHeight)

	t.Log("")
	t.Log("📦 Step 6: Preparing IBC transfer")
	amountToSend := math.NewInt(1_000_000)
	t.Logf("   - Amount to send: %s", amountToSend.String())
	t.Logf("   - Source: Chain From, User A (%s)", userA.FormattedAddress())
	t.Logf("   - Destination: Chain To, User B (%s)", userB.FormattedAddress())

	transfer := ibc.WalletAmount{
		Address: userB.FormattedAddress(),
		Denom:   chainFrom.Config().Denom,
		Amount:  amountToSend,
	}
	t.Log("✅ Transfer prepared")

	t.Log("")
	t.Log("🚀 Step 7: Executing IBC transfer")
	t.Logf("   - Sending IBC transfer on channel %s", channel.ChannelID)

	tx, err := chainFrom.SendIBCTransfer(ctx, channel.ChannelID, userA.KeyName(), transfer, ibc.TransferOptions{})
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

	t.Log("")
	t.Log("⏳ Step 8: Waiting for packet acknowledgment")
	t.Logf("   - Polling for ACK from height %d to %d", chainFromHeight, chainFromHeight+50)

	ack, err := testutil.PollForAck(ctx, chainFrom, chainFromHeight, chainFromHeight+50, tx.Packet)
	if err != nil {
		t.Logf("❌ Failed to receive ACK: %v", err)
	} else {
		t.Log("✅ Packet ACK received")
		if len(ack.Acknowledgement) > 0 {
			t.Logf("   - ACK: %s", string(ack.Acknowledgement))
		}
	}
	require.NoError(t, err)

	t.Log("")
	t.Log("⏳ Step 9: Waiting for blocks to finalize transfer")
	err = testutil.WaitForBlocks(ctx, 10, chainFrom)
	if err != nil {
		t.Logf("❌ Failed to wait for blocks: %v", err)
	}
	require.NoError(t, err)
	t.Log("✅ Blocks produced, transfer should be complete")

	t.Log("")
	t.Log("🔍 Step 10: Verifying final balances")

	chainFromDenom := transfertypes.GetPrefixedDenom(channel.Counterparty.PortID, channel.Counterparty.ChannelID, chainFrom.Config().Denom)
	chainFromIBCDenom := transfertypes.ParseDenomTrace(chainFromDenom).IBCDenom()
	t.Logf("   - IBC denom on Chain To: %s", chainFromIBCDenom)

	expectedBalA := balA.Sub(amountToSend)
	newBalA, err := chainFrom.GetBalance(ctx, userA.FormattedAddress(), chainFrom.Config().Denom)
	require.NoError(t, err)
	if newBalA.Equal(expectedBalA) {
		t.Logf("✅ Chain From balance correct: %s (expected: %s, sent: %s)",
			newBalA.String(), expectedBalA.String(), amountToSend.String())
	} else {
		t.Logf("❌ Chain From balance incorrect: got %s, expected %s (sent %s)",
			newBalA.String(), expectedBalA.String(), amountToSend.String())
	}
	require.True(t, newBalA.Equal(balA.Sub(amountToSend)))

	newBalB, err := chainTo.GetBalance(ctx, userB.FormattedAddress(), chainFromIBCDenom)
	require.NoError(t, err)
	if newBalB.Equal(amountToSend) {
		t.Logf("✅ Chain To balance correct: %s (expected: %s, received via IBC)",
			newBalB.String(), amountToSend.String())
	} else {
		t.Logf("❌ Chain To balance incorrect: got %s, expected %s",
			newBalB.String(), amountToSend.String())
	}
	require.True(t, newBalB.Equal(amountToSend))

	t.Log("")
	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	t.Log("✅ Basic IBC Transfer Test PASSED")
	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	t.Logf("   Summary:")
	t.Logf("   - Sent %s from Chain From to Chain To", amountToSend.String())
	t.Logf("   - Chain From balance: %s → %s", balA.String(), newBalA.String())
	t.Logf("   - Chain To balance: %s → %s (IBC denom)", balB.String(), newBalB.String())
}
