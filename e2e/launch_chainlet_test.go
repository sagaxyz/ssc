package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	interchaintest "github.com/cosmos/interchaintest/v10"
	"github.com/cosmos/interchaintest/v10/chain/cosmos"
	e2eutils "github.com/sagaxyz/ssc/e2e/utils"
	"github.com/stretchr/testify/require"
)

func TestChainletLaunch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	t.Log("🧪 Starting Chainlet Launch Test")
	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	t.Parallel()

	ctx := context.Background()

	t.Log("📡 Step 1: Creating and starting single-chain network")
	icn, err := e2eutils.CreateAndStartFullyConnectedNetwork(t, ctx, e2eutils.WithNChains(1))
	if err != nil {
		t.Logf("❌ Failed to create network: %v", err)
	}
	require.NoError(t, err)
	t.Log("✅ Network created successfully")

	chain, err := icn.GetChain(0)
	require.NoError(t, err)

	_, ok := chain.(*cosmos.CosmosChain)
	require.True(t, ok)
	t.Logf("   - Chain: %s (denom: %s)", chain.Config().Name, chain.Config().Denom)

	t.Log("")
	t.Log("💰 Step 2: Funding test users")
	fundAmt := math.NewInt(10_000_000)
	t.Logf("   - Funding amount: %s", fundAmt.String())

	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", fundAmt, chain, chain)
	require.Len(t, users, 2)
	alice := users[0]
	bob := users[1]

	t.Logf("   - Alice: %s", alice.FormattedAddress())
	t.Logf("   - Bob: %s", bob.FormattedAddress())
	t.Log("✅ Users funded")

	denom := chain.Config().Denom
	fees := "5000" + denom
	t.Logf("   - Transaction fees: %s", fees)

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

	t.Log("")
	t.Log("📦 Step 3: Creating chainlet stack")
	t.Log("   - Stack name: sagaevm")
	t.Log("   - Description: Your personal EVM")
	t.Log("   - Image: sagaxyz/sagaevm:0.7.0")
	t.Log("   - Version: 0.7.0")
	t.Log("   - Min Deposit: 1000" + denom)
	t.Log("   - Min Topup: 1000" + denom)

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
	if err != nil {
		t.Logf("❌ Failed to create stack: %v", err)
	}
	require.NoError(t, err)
	assertTxCode("create sagaevm stack", "0", code, txh)

	t.Log("")
	t.Log("🔄 Step 4: Updating chainlet stack")
	t.Log("   - Stack name: sagaevm")
	t.Log("   - New Image: sagaxyz/sagaevm:0.8.0")
	t.Log("   - New Version: 0.8.0")

	txh, code, _, err = e2eutils.ChainletUpdateStack(ctx, chain, bob, fees, e2eutils.UpdateStackParams{
		Name:        "sagaevm",
		Image:       "sagaxyz/sagaevm:0.8.0",
		Version:     "0.8.0",
		Hash:        "abc234",
		CcvConsumer: false,
	})
	if err != nil {
		t.Logf("❌ Failed to update stack: %v", err)
	}
	require.NoError(t, err)
	assertTxCode("update sagaevm to 0.8.0", "0", code, txh)

	t.Log("")
	t.Log("🚀 Step 5: Launching chainlets")
	t.Log("   - Launch variant 1: EVM Chain ID 100001 (version 0.7.0)")

	// valid 100001
	txh, code, _, err = e2eutils.ChainletLaunch(ctx, chain, bob, fees, e2eutils.LaunchChainletParams{
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
	if err != nil {
		t.Logf("❌ Failed to launch chainlet (100001): %v", err)
	}
	require.NoError(t, err)
	assertTxCode("launch 0.7.0 (100001)", "0", code, txh)
	t.Logf("   - Chainlet ID: mychain_100001-1")

	t.Log("   - Launch variant 2: EVM Chain ID 13371337 (version 0.7.0)")

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
	if err != nil {
		t.Logf("❌ Failed to launch chainlet (13371337): %v", err)
	}
	require.NoError(t, err)
	assertTxCode("launch 0.7.0 (13371337)", "0", code, txh)
	t.Logf("   - Chainlet ID: mychain_13371337-1")

	t.Log("   - Launch variant 3: Custom params (version 0.8.0)")

	// custom params on 0.8.0
	custom := fmt.Sprintf(`{"gasLimit":10000000,"genAcctBalances":"%s=1000,%s=100000"}`, alice.FormattedAddress(), bob.FormattedAddress())
	t.Logf("   - Custom JSON: %s", custom)
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
	if err != nil {
		t.Logf("❌ Failed to launch chainlet (custom): %v", err)
	}
	require.NoError(t, err)
	assertTxCode("launch 0.8.0 (custom)", "0", code, txh)
	t.Logf("   - Chainlet ID: kukkoo_<id>-1")

	t.Log("")
	t.Log("🔍 Step 6: Querying chainlet information")
	t.Log("   - Querying epoch infos")

	_, _, err = e2eutils.QueryJSON(ctx, chain, "epochs", "epoch-infos")
	if err != nil {
		t.Logf("❌ Failed to query epoch infos: %v", err)
	}
	require.NoError(t, err)
	t.Log("✅ Epoch infos queried")

	t.Log("   - Querying chainlet stacks list")
	var stacks struct {
		Stacks []any `json:"ChainletStacks"`
	}
	err = e2eutils.QueryInto(ctx, chain, &stacks, "chainlet", "list-chainlet-stack", "-o", "json")
	if err != nil {
		t.Logf("❌ Failed to list chainlet stacks: %v", err)
	}
	require.NoError(t, err)
	if len(stacks.Stacks) == 1 {
		t.Logf("✅ Found %d chainlet stack (expected: 1)", len(stacks.Stacks))
	} else {
		t.Logf("❌ Found %d chainlet stacks (expected: 1)", len(stacks.Stacks))
	}
	require.Len(t, stacks.Stacks, 1)

	t.Log("   - Querying chainlet stack details")
	var get struct {
		Stack struct {
			Versions []any `json:"versions"`
		} `json:"ChainletStack"`
	}
	err = e2eutils.QueryInto(ctx, chain, &get, "chainlet", "get-chainlet-stack", "sagaevm", "-o", "json")
	if err != nil {
		t.Logf("❌ Failed to get chainlet stack: %v", err)
	}
	require.NoError(t, err)
	if len(get.Stack.Versions) == 2 {
		t.Logf("✅ Found %d versions in stack (expected: 2)", len(get.Stack.Versions))
	} else {
		t.Logf("❌ Found %d versions in stack (expected: 2)", len(get.Stack.Versions))
	}
	require.Len(t, get.Stack.Versions, 2)

	t.Log("   - Querying chainlets list")
	var cl struct {
		Chainlets []any `json:"Chainlets"`
	}
	err = e2eutils.QueryInto(ctx, chain, &cl, "chainlet", "list-chainlets", "-o", "json")
	if err != nil {
		t.Logf("❌ Failed to list chainlets: %v", err)
	}
	require.NoError(t, err)
	if len(cl.Chainlets) >= 3 {
		t.Logf("✅ Found %d chainlets (expected: >= 3)", len(cl.Chainlets))
	} else {
		t.Logf("❌ Found %d chainlets (expected: >= 3)", len(cl.Chainlets))
	}
	require.GreaterOrEqual(t, len(cl.Chainlets), 3)

	t.Log("")
	t.Log("⏳ Step 7: Waiting for billing history to be available")
	target := "mychain_100001-1"
	t.Logf("   - Target chainlet: %s", target)
	t.Log("   - Polling for billing history and validator payout history")

	err = e2eutils.PollUntil(ctx, 24, 500*time.Millisecond, func() error {
		if _, _, err := e2eutils.QueryJSON(ctx, chain, "billing", "get-billing-history", target); err != nil {
			return err
		}
		if _, _, err := e2eutils.QueryJSON(ctx, chain, "billing", "get-validator-payout-history", target); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Logf("❌ Failed to get billing history: %v", err)
	}
	require.NoError(t, err)
	t.Log("✅ Billing history available")

	t.Log("")
	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	t.Log("✅ Chainlet Launch Test PASSED")
	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	t.Logf("   Summary:")
	t.Logf("   - Created chainlet stack: sagaevm")
	t.Logf("   - Updated stack to version 0.8.0")
	t.Logf("   - Launched 3 chainlets (2 with EVM Chain IDs, 1 with custom params)")
	t.Logf("   - Verified stack and chainlet queries")
	t.Logf("   - Confirmed billing history is available")
}

/* ---------------- local helper ---------------- */

func mustQueryJSON(ctx context.Context, chain *cosmos.CosmosChain, args ...string) ([]byte, error) {
	stdout, stderr, err := e2eutils.QueryJSON(ctx, chain, args...)
	if err != nil {
		return stdout, fmt.Errorf("query failed: %v; stderr=%s", err, string(stderr))
	}
	return stdout, nil
}
