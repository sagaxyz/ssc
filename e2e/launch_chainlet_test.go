package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	e2eutils "github.com/sagaxyz/ssc/e2e/utils"
	interchaintest "github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/stretchr/testify/require"
)

type txResp struct {
	TxHash string `json:"txhash"`
	Code   uint32 `json:"code"`
}

func TestChainletLaunch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	t.Parallel()

	ctx := context.Background()

	// Spin up a single SSC chain
	icn, err := e2eutils.CreateAndStartFullyConnectedNetwork(t, ctx, e2eutils.WithNChains(1))
	require.NoError(t, err)
	chain, err := icn.GetChain(0)
	require.NoError(t, err)

	// Assert type so address formatting, etc., matches Cosmos expectations
	_, ok := chain.(*cosmos.CosmosChain)
	require.True(t, ok)

	// Two programmatic users on SAME chain
	fundAmt := math.NewInt(10_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", fundAmt, chain, chain)
	require.Len(t, users, 2)
	alice := users[0]
	bob := users[1]

	denom := chain.Config().Denom
	fees := "5000" + denom
	gasLimit := "500000"
	chainletDenom := "asaga"

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
		// Fetch the full tx JSON for rich debug
		ch, ok := chain.(*cosmos.CosmosChain)
		if !ok {
			t.Fatalf("❌ %s (got code=%d, want %s)\n  txhash: %s\n  cannot fetch full tx details, chain is not cosmos", label, got, want, txhash)
		}

		qout, qerr := mustQueryJSON(ctx, ch, "tx", txhash, "-o", "json")
		t.Fatalf("❌ %s (got code=%d, want %s)\n  txhash: %s\n  q tx: %s\n  query_err: %v",
			label, got, want, txhash, string(qout), qerr)
	}

	// === create stack ===
	{
		txh, code, _, err := e2eutils.ExecTxJSON(ctx, chain, bob, fees,
			"chainlet", "create-chainlet-stack",
			"sagaevm", "Your personal EVM",
			"sagaxyz/sagaevm:0.7.0", "0.7.0", "abc123",
			"1000"+denom, "1000"+denom, "false",
		)
		require.NoError(t, err, "create sagaevm stack exec error")
		assertTxCode("create sagaevm stack", "0", code, txh)
	}

	// === update stack ===
	{
		txh, code, _, err := e2eutils.ExecTxJSON(ctx, chain, bob, fees,
			"chainlet", "update-chainlet-stack",
			"sagaevm",
			"sagaxyz/sagaevm:0.8.0", "0.8.0", "abc234", "false",
		)
		require.NoError(t, err)
		assertTxCode("update existing stack to 0.8.0", "0", code, txh)
	}

	// === launch chainlet ===
	{
		txh, code, _, err := e2eutils.ExecTxJSON(ctx, chain, bob, fees,
			"chainlet", "launch-chainlet",
			bob.FormattedAddress(), "sagaevm", "0.7.0", "mychain", chainletDenom, "{}",
			"--evm-chain-id", "100001", "--network-version", "1", "--gas", gasLimit,
		)
		require.NoError(t, err)
		assertTxCode("valid launch (0.7.0, 100001)", "0", code, txh)

		txh, code, _, err = e2eutils.ExecTxJSON(ctx, chain, bob, fees,
			"chainlet", "launch-chainlet",
			bob.FormattedAddress(), "sagaevm", "0.7.0", "mychain", chainletDenom, "{}",
			"--evm-chain-id", "13371337", "--network-version", "1", "--gas", gasLimit,
		)
		require.NoError(t, err)
		assertTxCode("valid launch (0.7.0, 13371337)", "0", code, txh)

		custom := fmt.Sprintf(`{"gasLimit":10000000,"genAcctBalances":"%s=1000,%s=100000"}`, alice.FormattedAddress(), bob.FormattedAddress())
		txh, code, _, err = e2eutils.ExecTxJSON(ctx, chain, bob, fees,
			"chainlet", "launch-chainlet",
			bob.FormattedAddress(), "sagaevm", "0.8.0", "kukkoo", chainletDenom, custom,
			"--gas", gasLimit,
		)
		require.NoError(t, err)
		assertTxCode("custom params launch (0.8.0)", "0", code, txh)
	}

	// === queries & billing ===
	{
		_, _, err := e2eutils.QueryJSON(ctx, chain, "epochs", "epoch-infos")
		require.NoError(t, err, "epochs query failed")

		raw, _, err := e2eutils.QueryJSON(ctx, chain, "chainlet", "list-chainlet-stack", "-o", "json")
		require.NoError(t, err)
		var stacks struct {
			Stacks []any `json:"ChainletStacks"`
		}
		require.NoError(t, json.Unmarshal(raw, &stacks), "parse list-chainlet-stack")
		require.Len(t, stacks.Stacks, 1, "expected exactly 1 stack")

		raw, _, err = e2eutils.QueryJSON(ctx, chain, "chainlet", "get-chainlet-stack", "sagaevm", "-o", "json")
		require.NoError(t, err)
		var get struct {
			Stack struct {
				Versions []any `json:"versions"`
			} `json:"ChainletStack"`
		}
		require.NoError(t, json.Unmarshal(raw, &get), "parse get-chainlet-stack")
		require.Len(t, get.Stack.Versions, 2, "expected 2 versions in sagaevm stack")

		raw, _, err = e2eutils.QueryJSON(ctx, chain, "chainlet", "list-chainlets", "-o", "json")
		require.NoError(t, err)
		var cl struct {
			Chainlets []any `json:"Chainlets"`
		}
		require.NoError(t, json.Unmarshal(raw, &cl), "parse list-chainlets")
		require.GreaterOrEqual(t, len(cl.Chainlets), 3, "expected at least 3 chainlets")
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

/* ---------------- local helpers for this test ---------------- */

func mustQueryJSON(ctx context.Context, chain *cosmos.CosmosChain, args ...string) ([]byte, error) {
	stdout, stderr, err := e2eutils.QueryJSON(ctx, chain, args...)
	if err != nil {
		return stdout, fmt.Errorf("query failed: %v; stderr=%s", err, string(stderr))
	}
	return stdout, nil
}
