package e2e

import (
	"context"
	"encoding/json"
	"testing"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	interchaintest "github.com/cosmos/interchaintest/v10"
	"github.com/cosmos/interchaintest/v10/ibc"
	"github.com/cosmos/interchaintest/v10/testutil"
	e2eutils "github.com/sagaxyz/ssc/e2e/utils"
	"github.com/stretchr/testify/require"
)

// GMPMessage represents an Axelar GMP message in the ICS-20 memo field
type GMPMessage struct {
	SourceChain   string `json:"source_chain"`
	SourceAddress string `json:"source_address"`
	Payload       []byte `json:"payload"`
	Type          int64  `json:"type"`
}

const (
	TypeGeneralMessageWithToken = 2
)

// TestGMPWithPFMTransfer tests that GMP middleware correctly extracts PFM forward instructions
// from an ABI-encoded payload and forwards via PFM across 3 chains.
//
// This simulates an Axelar GMP message arriving on Saga with type=2 (TypeGeneralMessageWithToken)
// where the payload contains ABI-encoded PFM forward instructions.
func TestGMPWithPFMTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Log("Starting GMP + PFM Transfer Test")
	t.Log("This test verifies that GMP middleware correctly extracts PFM forward instructions")

	t.Parallel()
	ctx := context.Background()

	// Set up 3 chains: A -> B -> C
	// Chain A simulates Axelar sending a GMP message
	// Chain B is Saga (with GMP + PFM middleware)
	// Chain C is the final destination
	pathAB := e2eutils.RelayerPath{0, 1}
	pathBC := e2eutils.RelayerPath{1, 2}

	t.Log("Step 1: Creating 3-chain network")

	icn, err := e2eutils.CreateAndStartFullyConnectedNetwork(t, ctx,
		e2eutils.WithNChains(3),
		e2eutils.WithRelayerPaths(pathAB, pathBC),
	)
	require.NoError(t, err)
	t.Log("Network created")

	chainA, err := icn.GetChain(0)
	require.NoError(t, err)
	chainB, err := icn.GetChain(1)
	require.NoError(t, err)
	chainC, err := icn.GetChain(2)
	require.NoError(t, err)

	t.Logf("Chain A: %s", chainA.Config().Name)
	t.Logf("Chain B: %s", chainB.Config().Name)
	t.Logf("Chain C: %s", chainC.Config().Name)

	// Fund users
	t.Log("Step 2: Funding test users")
	fundAmount := math.NewInt(10_000_000)

	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", fundAmount, chainA, chainB, chainC)
	userA := users[0]
	userB := users[1]
	userC := users[2]

	t.Logf("User A: %s", userA.FormattedAddress())
	t.Logf("User B: %s", userB.FormattedAddress())
	t.Logf("User C: %s", userC.FormattedAddress())

	err = testutil.WaitForBlocks(ctx, 3, chainA, chainB, chainC)
	require.NoError(t, err)

	// Get channel info
	t.Log("Step 3: Getting channel information")
	channelAB, err := icn.GetChannelInfo(ctx, pathAB)
	require.NoError(t, err)
	t.Logf("Channel A->B: %s", channelAB.ChannelID)

	channelBC, err := icn.GetChannelInfo(ctx, pathBC)
	require.NoError(t, err)
	t.Logf("Channel B->C: %s", channelBC.ChannelID)

	// Record initial balances
	t.Log("Step 4: Recording initial balances")
	balA, err := chainA.GetBalance(ctx, userA.FormattedAddress(), chainA.Config().Denom)
	require.NoError(t, err)
	t.Logf("User A balance: %s", balA.String())

	chainAHeight, err := chainA.Height(ctx)
	require.NoError(t, err)

	// Build GMP message with ABI-encoded PFM forward instructions
	t.Log("Step 5: Building GMP message with ABI-encoded PFM payload")

	// Create PFM forward instructions as JSON
	pfmPayload := map[string]interface{}{
		"forward": map[string]interface{}{
			"receiver": userC.FormattedAddress(),
			"port":     channelBC.PortID,
			"channel":  channelBC.ChannelID,
		},
	}
	pfmPayloadJSON, err := json.Marshal(pfmPayload)
	require.NoError(t, err)
	t.Logf("PFM payload: %s", string(pfmPayloadJSON))

	// ABI-encode the PFM payload as a string
	payloadType, err := abi.NewType("string", "", nil)
	require.NoError(t, err)
	args := abi.Arguments{{Type: payloadType}}
	abiEncodedPayload, err := args.Pack(string(pfmPayloadJSON))
	require.NoError(t, err)
	t.Logf("ABI-encoded payload length: %d bytes", len(abiEncodedPayload))

	// Create GMP message
	gmpMsg := GMPMessage{
		SourceChain:   "Ethereum",
		SourceAddress: "0xce16F69375520ab01377ce7B88f5BA8C48F8D666",
		Payload:       abiEncodedPayload,
		Type:          TypeGeneralMessageWithToken,
	}
	gmpMemo, err := json.Marshal(gmpMsg)
	require.NoError(t, err)
	t.Logf("GMP memo: %s", string(gmpMemo))

	// Send IBC transfer with GMP memo
	t.Log("Step 6: Sending IBC transfer with GMP memo")
	amountToSend := math.NewInt(1_000_000)

	transfer := ibc.WalletAmount{
		Address: userB.FormattedAddress(), // Initial receiver on chain B
		Denom:   chainA.Config().Denom,
		Amount:  amountToSend,
	}

	tx, err := chainA.SendIBCTransfer(ctx, channelAB.ChannelID, userA.KeyName(), transfer, ibc.TransferOptions{
		Memo: string(gmpMemo),
	})
	require.NoError(t, err)
	t.Logf("Transaction hash: %s", tx.TxHash)

	err = tx.Validate()
	require.NoError(t, err)
	t.Log("Transfer submitted")

	// Wait for ACK on chain A
	t.Log("Step 7: Waiting for packet acknowledgment")
	_, err = testutil.PollForAck(ctx, chainA, chainAHeight, chainAHeight+50, tx.Packet)
	require.NoError(t, err)
	t.Log("ACK received")

	// Wait for PFM forwarding to complete
	t.Log("Step 8: Waiting for PFM forwarding (B->C)")
	err = testutil.WaitForBlocks(ctx, 30, chainA, chainB, chainC)
	require.NoError(t, err)

	// Verify balances
	t.Log("Step 9: Verifying final balances")

	// Chain A: Should have sent amount
	newBalA, err := chainA.GetBalance(ctx, userA.FormattedAddress(), chainA.Config().Denom)
	require.NoError(t, err)
	expectedBalA := balA.Sub(amountToSend)
	require.Equal(t, expectedBalA.String(), newBalA.String(), "Chain A balance incorrect")
	t.Logf("Chain A balance: %s -> %s (sent %s)", balA.String(), newBalA.String(), amountToSend.String())

	// Calculate IBC denoms
	firstHopDenom := transfertypes.GetPrefixedDenom(channelAB.Counterparty.PortID, channelAB.Counterparty.ChannelID, chainA.Config().Denom)
	secondHopDenom := transfertypes.GetPrefixedDenom(channelBC.Counterparty.PortID, channelBC.Counterparty.ChannelID, firstHopDenom)
	firstHopIBCDenom := transfertypes.ParseDenomTrace(firstHopDenom).IBCDenom()
	secondHopIBCDenom := transfertypes.ParseDenomTrace(secondHopDenom).IBCDenom()

	// Chain B: Should have 0 (PFM forwarded immediately)
	newBalB, err := chainB.GetBalance(ctx, userB.FormattedAddress(), firstHopIBCDenom)
	require.NoError(t, err)
	require.Equal(t, "0", newBalB.String(), "Chain B should have 0 balance (PFM should forward)")
	t.Logf("Chain B balance: %s (expected 0, PFM forwarded)", newBalB.String())

	// Chain C: Should have received the forwarded amount
	newBalC, err := chainC.GetBalance(ctx, userC.FormattedAddress(), secondHopIBCDenom)
	require.NoError(t, err)
	require.Equal(t, amountToSend.String(), newBalC.String(), "Chain C should have received the forwarded amount")
	t.Logf("Chain C balance: %s (expected %s)", newBalC.String(), amountToSend.String())

	t.Log("GMP + PFM Transfer Test PASSED")
	t.Log("GMP correctly extracted PFM forward instructions and PFM forwarded the packet")
}
