package utils

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
)

// ExecTxJSON signs & broadcasts a tx using an Interchaintest wallet (no CLI keyring).
// Forces "-o json" and "-y"; injects "--fees <fees>" when fees != "".
// Accepts both JSON stdout (with "txhash") and raw 64-hex txhash stdout.
// Waits a bit, then queries the tx for the final code.
// Returns: txhash, code, raw stdout (string) from initial ExecTx, error.
func ExecTxJSON(ctx context.Context, chain ibc.Chain, signer ibc.Wallet, fees string, args ...string) (string, uint32, string, error) {
	node, err := getNode(chain)
	if err != nil {
		return "", 0, "", fmt.Errorf("get node: %w", err)
	}

	// ensure required flags
	if !slices.Contains(args, "-o") {
		args = append(args, "-o", "json")
	}
	if !slices.Contains(args, "-y") {
		args = append(args, "-y")
	}
	if fees != "" && !slices.Contains(args, "--fees") {
		args = append(args, "--fees", fees)
	}

	// Exec (your version returns (string, error))
	outStr, runErr := node.ExecTx(ctx, signer.KeyName(), args...)
	if runErr != nil {
		return "", 0, outStr, fmt.Errorf("exec tx: %w; stdout=%s", runErr, outStr)
	}

	// Allow tx to index a little
	_ = testutil.WaitForBlocks(ctx, 2, chain)

	// Parse txhash from stdout: JSON or raw hex
	txhash := parseTxHash(outStr)
	if txhash == "" {
		return "", 0, outStr, fmt.Errorf("missing txhash in output (stdout=%s)", outStr)
	}

	// Query tx with retries (indexing can still lag beyond 2 blocks)
	code, qErr := queryTxCodeWithRetry(ctx, chain, txhash, 8, 500*time.Millisecond)
	if qErr != nil {
		return txhash, 0, outStr, fmt.Errorf("query tx failed for %s: %w", txhash, qErr)
	}

	return txhash, code, outStr, nil
}

// QueryJSON runs a read-only query via the node (returns stdout, stderr, err).
func QueryJSON(ctx context.Context, chain ibc.Chain, args ...string) ([]byte, []byte, error) {
	node, err := getNode(chain)
	if err != nil {
		return nil, nil, fmt.Errorf("get node: %w", err)
	}
	return node.ExecQuery(ctx, args...)
}

// PollUntil runs fn up to `attempts` times, sleeping `sleep` between tries.
func PollUntil(ctx context.Context, attempts int, sleep time.Duration, fn func() error) error {
	for i := 0; i < attempts; i++ {
		if err := fn(); err == nil {
			return nil
		}
		select {
		case <-time.After(sleep):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return fmt.Errorf("condition not met in %d attempts", attempts)
}

// ---- internal helpers ----

// parseTxHash tries JSON first, then plain hex, then naive extract.
func parseTxHash(stdout string) string {
	stdout = strings.TrimSpace(stdout)

	// 1) JSON with "txhash"
	var pre struct {
		TxHash string `json:"txhash"`
	}
	if json.Unmarshal([]byte(stdout), &pre) == nil && pre.TxHash != "" {
		return pre.TxHash
	}

	// 2) raw 64-hex (uppercase/lowercase)
	if isHex64(stdout) {
		return stdout
	}

	// 3) naive extract from JSON-ish text
	return ExtractField(stdout, "txhash")
}

// isHex64 returns true if s is exactly 64 hex chars (no 0x prefix).
func isHex64(s string) bool {
	if len(s) != 64 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

func queryTxCodeWithRetry(ctx context.Context, chain ibc.Chain, txhash string, attempts int, sleep time.Duration) (uint32, error) {
	node, err := getNode(chain)
	if err != nil {
		return 0, err
	}
	var lastErr error
	for i := 0; i < attempts; i++ {
		qOut, qErrStd, qErr := node.ExecQuery(ctx, "tx", txhash, "-o", "json")
		if qErr == nil {
			var qr struct {
				Code uint32 `json:"code"`
			}
			if err := json.Unmarshal(qOut, &qr); err != nil {
				return 0, fmt.Errorf("unmarshal query tx: %w; raw=%s", err, string(qOut))
			}
			return qr.Code, nil
		}
		lastErr = fmt.Errorf("query attempt %d/%d: %v (stderr=%s)", i+1, attempts, qErr, string(qErrStd))
		select {
		case <-time.After(sleep):
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}
	return 0, lastErr
}

// ExtractField fetches a naive `"key":"value"` from a small JSON string.
func ExtractField(s, key string) string {
	key = `"` + key + `":"`
	i := strings.Index(s, key)
	if i < 0 {
		return ""
	}
	i += len(key)
	j := strings.Index(s[i:], `"`)
	if j < 0 {
		return ""
	}
	return s[i : i+j]
}
