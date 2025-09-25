package e2e

import (
	"context"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	e2eutils "github.com/sagaxyz/ssc/e2e/utils"
	interchaintest "github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/stretchr/testify/require"
)

// TestBasicIBCTransfer is a basic test to start 2 SSC chains and send an IBC transfer between them.
func TestBasicIBCTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()
	ctx := context.Background()

	pathAB := e2eutils.RelayerPath{0, 1}
	icn, err := e2eutils.CreateAndStartFullyConnectedNetwork(t, ctx,
		e2eutils.WithNChains(2),
		e2eutils.WithRelayerPaths(pathAB),
	)
	require.NoError(t, err)
	require.NotNil(t, icn)

	chainFrom, err := icn.GetChain(0)
	require.NoError(t, err)
	chainTo, err := icn.GetChain(1)
	require.NoError(t, err)

	fundAmount := math.NewInt(10_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", fundAmount, chainFrom, chainTo)
	userA := users[0]
	userB := users[1]

	err = testutil.WaitForBlocks(ctx, 3, chainFrom, chainTo)
	require.NoError(t, err)

	balA, err := chainFrom.GetBalance(ctx, userA.FormattedAddress(), chainFrom.Config().Denom)
	require.NoError(t, err)
	require.True(t, balA.Equal(fundAmount))
	balB, err := chainTo.GetBalance(ctx, userB.FormattedAddress(), chainTo.Config().Denom)
	require.NoError(t, err)
	require.True(t, balB.Equal(fundAmount))

	fmt.Printf("userA=%v, bal=%v\n", userA.FormattedAddress(), balA)
	fmt.Printf("userB=%v, bal=%v\n", userB.FormattedAddress(), balB)

	amountToSend := math.NewInt(1_000_000)
	transfer := ibc.WalletAmount{
		Address: userB.FormattedAddress(),
		Denom:   chainFrom.Config().Denom,
		Amount:  amountToSend,
	}

	channel, err := icn.GetChannelInfo(ctx, pathAB)
	require.NoError(t, err)

	chainFromHeight, err := chainFrom.Height(ctx)
	require.NoError(t, err)

	tx, err := chainFrom.SendIBCTransfer(ctx, channel.ChannelID, userA.KeyName(), transfer, ibc.TransferOptions{})
	require.NoError(t, err)
	require.NoError(t, tx.Validate())

	_, err = testutil.PollForAck(ctx, chainFrom, chainFromHeight, chainFromHeight+50, tx.Packet)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 10, chainFrom)
	require.NoError(t, err)

	chainFromDenom := transfertypes.GetPrefixedDenom(channel.Counterparty.PortID, channel.Counterparty.ChannelID, chainFrom.Config().Denom)
	chainFromIBCDenom := transfertypes.ParseDenomTrace(chainFromDenom).IBCDenom()

	newBalA, err := chainFrom.GetBalance(ctx, userA.FormattedAddress(), chainFrom.Config().Denom)
	require.NoError(t, err)
	require.True(t, newBalA.Equal(balA.Sub(amountToSend)))

	newBalB, err := chainTo.GetBalance(ctx, userB.FormattedAddress(), chainFromIBCDenom)
	require.NoError(t, err)
	require.True(t, newBalB.Equal(amountToSend))
}
