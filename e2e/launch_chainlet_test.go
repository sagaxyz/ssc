package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	e2eutils "github.com/sagaxyz/ssc/e2e/utils"
	interchaintest "github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/stretchr/testify/require"
)

func TestChainletLaunch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	t.Parallel()

	ctx := context.Background()

	icn, err := e2eutils.CreateAndStartFullyConnectedNetwork(t, ctx, e2eutils.WithNChains(1))
	require.NoError(t, err)
	chain, err := icn.GetChain(0)
	require.NoError(t, err)

	_, ok := chain.(*cosmos.CosmosChain)
	require.True(t, ok)

	// two wallets on same chain
	fundAmt := math.NewInt(10_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", fundAmt, chain, chain)
	require.Len(t, users, 2)
	alice := users[0]
	bob := users[1]

	denom := chain.Config().Denom
	fees := "5000" + denom

	expect := func(want string, got uint32) bool {
		if want == "nonzero" {
			return got != 0
		}
		return want == "0" && got == 0
	}
	assertTxCode := func(label, want string, got uint32, txhash string) {
		if expect(want, got) {
			t.Logf("✅ %s (code=%d, tx=%s)", label, got, txhash)
			return
		}
		ch := chain.(*cosmos.CosmosChain)
		qout, qerr := mustQueryJSON(ctx, ch, "tx", txhash, "-o", "json")
		t.Fatalf("❌ %s (got code=%d, want %s)\n  txhash: %s\n  q tx: %s\n  err: %v",
			label, got, want, txhash, string(qout), qerr)
	}

	// --- create ---
	{
		txh, code, _, err := e2eutils.ChainletCreateStack(ctx, chain, bob, fees, e2eutils.CreateStackParams{
			Name:        "sagaevm",
			Description: "Your personal EVM",
			Image:       "sagaxyz/sagaevm:0.7.0",
			Version:     "0.7.0",
			Hash:        "abc123",
			MinDeposit:  "1000" + denom,
			MinTopup:    "1000" + denom,
			CcvConsumer: false,
		})
		require.NoError(t, err)
		assertTxCode("create sagaevm stack", "0", code, txh)
	}

	// --- update ---
	{
		txh, code, _, err := e2eutils.ChainletUpdateStack(ctx, chain, bob, fees, e2eutils.UpdateStackParams{
			Name:        "sagaevm",
			Image:       "sagaxyz/sagaevm:0.8.0",
			Version:     "0.8.0",
			Hash:        "abc234",
			CcvConsumer: false,
		})
		require.NoError(t, err)
		assertTxCode("update sagaevm to 0.8.0", "0", code, txh)
	}

	// --- launch variants ---
	{
		// valid 100001
		txh, code, _, err := e2eutils.ChainletLaunch(ctx, chain, bob, fees, e2eutils.LaunchChainletParams{
			OwnerAddr:      bob.FormattedAddress(),
			StackName:      "sagaevm",
			StackVersion:   "0.7.0",
			ChainletID:     "mychain",
			ChainletDenom:  "asaga",
			CustomJSON:     "{}",
			EVMChainID:     "100001",
			NetworkVersion: "1",
			Gas:            "500000",
		})
		require.NoError(t, err)
		assertTxCode("launch 0.7.0 (100001)", "0", code, txh)

		// valid 13371337
		txh, code, _, err = e2eutils.ChainletLaunch(ctx, chain, bob, fees, e2eutils.LaunchChainletParams{
			OwnerAddr:      bob.FormattedAddress(),
			StackName:      "sagaevm",
			StackVersion:   "0.7.0",
			ChainletID:     "mychain",
			ChainletDenom:  "asaga",
			CustomJSON:     "{}",
			EVMChainID:     "13371337",
			NetworkVersion: "1",
			Gas:            "500000",
		})
		require.NoError(t, err)
		assertTxCode("launch 0.7.0 (13371337)", "0", code, txh)

		// custom params on 0.8.0
		custom := fmt.Sprintf(`{"gasLimit":10000000,"genAcctBalances":"%s=1000,%s=100000"}`, alice.FormattedAddress(), bob.FormattedAddress())
		txh, code, _, err = e2eutils.ChainletLaunch(ctx, chain, bob, fees, e2eutils.LaunchChainletParams{
			OwnerAddr:      bob.FormattedAddress(),
			StackName:      "sagaevm",
			StackVersion:   "0.8.0",
			ChainletID:     "kukkoo",
			ChainletDenom:  "asaga",
			CustomJSON:     custom,
			EVMChainID:     "",
			NetworkVersion: "",
			Gas:            "500000",
		})
		require.NoError(t, err)
		assertTxCode("launch 0.8.0 (custom)", "0", code, txh)
	}

	// --- queries & billing ---
	{
		_, _, err := e2eutils.QueryJSON(ctx, chain, "epochs", "epoch-infos")
		require.NoError(t, err)

		var stacks struct {
			Stacks []any `json:"ChainletStacks"`
		}
		require.NoError(t, e2eutils.QueryInto(ctx, chain, &stacks, "chainlet", "list-chainlet-stack", "-o", "json"))
		require.Len(t, stacks.Stacks, 1)

		var get struct {
			Stack struct {
				Versions []any `json:"versions"`
			} `json:"ChainletStack"`
		}
		require.NoError(t, e2eutils.QueryInto(ctx, chain, &get, "chainlet", "get-chainlet-stack", "sagaevm", "-o", "json"))
		require.Len(t, get.Stack.Versions, 2)

		var cl struct {
			Chainlets []any `json:"Chainlets"`
		}
		require.NoError(t, e2eutils.QueryInto(ctx, chain, &cl, "chainlet", "list-chainlets", "-o", "json"))
		require.GreaterOrEqual(t, len(cl.Chainlets), 3)
	}

	target := "mychain_100001-1"
	require.NoError(t, e2eutils.PollUntil(ctx, 24, 500*time.Millisecond, func() error {
		if _, _, err := e2eutils.QueryJSON(ctx, chain, "billing", "get-billing-history", target); err != nil {
			return err
		}
		if _, _, err := e2eutils.QueryJSON(ctx, chain, "billing", "get-validator-payout-history", target); err != nil {
			return err
		}
		return nil
	}))
}

/* ---------------- local helper ---------------- */

func mustQueryJSON(ctx context.Context, chain *cosmos.CosmosChain, args ...string) ([]byte, error) {
	stdout, stderr, err := e2eutils.QueryJSON(ctx, chain, args...)
	if err != nil {
		return stdout, fmt.Errorf("query failed: %v; stderr=%s", err, string(stderr))
	}
	return stdout, nil
}
