#!/usr/bin/env bash
set -euo pipefail

# ==============================
# Configuration (match happypath.sh style)
# ==============================
tx_block_time=.3
key="bob"
keyring_backend="test"
denom="utsaga"              # Denom A (priority)
denom_b="utagas"            # Denom B (secondary)
chainlet_denom="asaga"
gas_limit=500000
fees="5000${denom}"

# Stack + versions
stack_name="sagaevm"
stack_img_v1="sagaxyz/sagaevm:0.7.0"
stack_ver_v1="0.7.0"
stack_checksum_v1="abc123"

stack_img_v2="sagaxyz/sagaevm:0.8.0"
stack_ver_v2="0.8.0"
stack_checksum_v2="def456"

# Fees (keep small for local billing)
epoch_fee_a="1000${denom}"
setup_fee_a="1000${denom}"     # can omit (defaults to epoch), kept explicit for clarity
epoch_fee_b="1000${denom_b}"
setup_fee_b="1000${denom_b}"

# Chain IDs used across tests
cid_a_1="mychaina_100001-1"  # Test 1
cid_b_1="mychainb_1337-1"    # Test 2
cid_b_2="mychainb_4242-1"    # Test 3
cid_a_2="mychaina_5151-1"    # Test 4

# ==============================
# Utilities (same style as happypath.sh)
# ==============================
cleanup_sscd() {
  echo "🧹 Cleaning up sscd..."
  if pid=$(pgrep -f "sscd start --home"); then
    echo "🛑 Stopping sscd (pid $pid)..."
    kill "$pid"
  fi
}

wait_tx() {
  local hash=$1
  echo "⏳ Waiting ${tx_block_time}s for tx $hash..."
  sleep "${tx_block_time}"
}

run_test() {
  local cmd=$1
  local want=$2
  local pass_msg=$3
  local fail_msg=$4

  echo "▶️  ${pass_msg}"
  txhash=$(eval "${cmd}" | jq -r .txhash)
  wait_tx "$txhash"

  code=$(sscd q tx "$txhash" -o json | jq -r .code)

  local ok
  if [[ "$want" == "nonzero" ]]; then
    (( code != 0 )) && ok=0 || ok=1
  else
    (( code == want )) && ok=0 || ok=1
  fi

  if (( ok == 0 )); then
    echo "✅ ${pass_msg}"
  else
    echo "❌ ${fail_msg} (got code=$code, want ${want})"
    cleanup_sscd
    exit 1
  fi
  echo
}

ensure_status() {
  local chain_id=$1
  local expect=$2
  if sscd q chainlet get-chainlet "$chain_id" -o json | jq -r '.Chainlet.status' | grep -q "$expect"; then
    echo "✅ $chain_id is $expect"
  else
    echo "❌ $chain_id not $expect"
    cleanup_sscd
    exit 1
  fi
}

sleep_billing() {
  local secs=${1:-70}
  echo "🛌 Sleeping ${secs}s for billing/payout..."
  sleep "$secs"
}

# # ==============================
# # Setup: stack + fees (A priority, B secondary)
# # ==============================
echo "=== create/update stack & fees ==="

run_test "sscd tx chainlet create-chainlet-stack ${stack_name} \"Your personal EVM\" \
  ${stack_img_v1} ${stack_ver_v1} ${stack_checksum_v1} ${epoch_fee_a} ${setup_fee_a} false \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "created ${stack_name} stack (A fees)" "failed to create ${stack_name}"

run_test "sscd tx chainlet update-chainlet-stack ${stack_name} \
  ${stack_img_v2} ${stack_ver_v2} ${stack_checksum_v2} false \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "updated ${stack_name} to ${stack_ver_v2}" "failed to update stack version"

# NEW: Use update-stack-fees with CSV tokens in order (A first => priority, then B).
# token format per CLI: "epoch[:setup]" ; if setup omitted, it copies epoch.
# Example (explicit): "1000utsaga:1000utsaga,1000utagas:1000utagas"
fees_csv="${epoch_fee_a}:${setup_fee_a},${epoch_fee_b}:${setup_fee_b}"
run_test "sscd tx chainlet update-stack-fees ${stack_name} \"${fees_csv}\" \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "updated fees (A then B; A priority)" "failed to update fees"

# ==============================
# Test 1 — Launch with A, bill, payout, withdraw, offline, restart with B, query pool A
# ==============================
echo "=== Test 1: A-priority launch, withdraw, restart with B ==="

run_test "sscd tx chainlet launch-chainlet \"\$(sscd keys show -a ${key} --keyring-backend ${keyring_backend})\" \
  ${stack_name} ${stack_ver_v1} mychaina ${chainlet_denom} '{}' \
  --evm-chain-id 100001 --network-version 1 --gas ${gas_limit} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "launched ${cid_a_1} with denom A" "failed to launch ${cid_a_1}"

ensure_status "${cid_a_1}" "STATUS_ONLINE"
sleep_billing 70

if sscd q billing get-billing-history "${cid_a_1}" > /dev/null 2>&1; then
  echo "✅ billing history fetched"
else
  echo "❌ billing history missing"; cleanup_sscd; exit 1
fi

if sscd q billing get-validator-payout-history "${cid_a_1}" > /dev/null 2>&1; then
  echo "✅ validator payout fetched"
else
  echo "❌ validator payout missing"; cleanup_sscd; exit 1
fi

# NEW: withdraw funds from escrow (owner)
run_test "sscd tx escrow withdraw ${cid_a_1} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "withdrew from ${cid_a_1} escrow" "failed to withdraw from ${cid_a_1}"

sleep_billing 60
ensure_status "${cid_a_1}" "STATUS_OFFLINE"

run_test "sscd tx escrow deposit 20000${denom_b} ${cid_a_1} \
  --from alice --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "deposited B into ${cid_a_1}" "failed to deposit B"

sleep_billing 65
ensure_status "${cid_a_1}" "STATUS_ONLINE"

echo "🔎 querying pool for A on ${cid_a_1}"
sscd q escrow balance "${cid_a_1}" "${denom}" -o json | jq .

# ==============================
# Test 2 — Launch with B, bill/payout, withdraw, offline, deposit A, restart, query pool B
# ==============================
echo "=== Test 2: B launch, withdraw, restart with A ==="

run_test "sscd tx chainlet launch-chainlet \"\$(sscd keys show -a ${key} --keyring-backend ${keyring_backend})\" \
  ${stack_name} ${stack_ver_v1} mychainb ${chainlet_denom} '{}' \
  --evm-chain-id 1337 --network-version 1 --gas ${gas_limit} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "launched ${cid_b_1} with denom B" "failed to launch ${cid_b_1}"

ensure_status "${cid_b_1}" "STATUS_ONLINE"
sleep_billing 70

if sscd q billing get-billing-history "${cid_b_1}" > /dev/null 2>&1; then
  echo "✅ billing history fetched"
else
  echo "❌ billing history missing"; cleanup_sscd; exit 1
fi

if sscd q billing get-validator-payout-history "${cid_b_1}" > /dev/null 2>&1; then
  echo "✅ validator payout fetched"
else
  echo "❌ validator payout missing"; cleanup_sscd; exit 1
fi

# NEW: withdraw funds from escrow (owner)
run_test "sscd tx escrow withdraw ${cid_b_1} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "withdrew from ${cid_b_1} escrow" "failed to withdraw from ${cid_b_1}"

sleep_billing 60
ensure_status "${cid_b_1}" "STATUS_OFFLINE"

run_test "sscd tx escrow deposit 20000${denom} ${cid_b_1} \
  --from alice --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "deposited A into ${cid_b_1}" "failed to deposit A"

sleep_billing 65
ensure_status "${cid_b_1}" "STATUS_ONLINE"

echo "🔎 querying pool for B on ${cid_b_1}"
sscd q escrow balance "${cid_b_1}" "${denom_b}" -o json | jq .

# ==============================
# Test 3 — Launch with B, deposit A before billing, ensure bill in A, payouts, pools,
#           withdraw -> offline -> deposit -> restart
# ==============================
echo "=== Test 3: B launch, pre-billing A deposit, bill in A, cycle withdraw/restart ==="

run_test "sscd tx chainlet launch-chainlet \"\$(sscd keys show -a ${key} --keyring-backend ${keyring_backend})\" \
  ${stack_name} ${stack_ver_v2} mychainb ${chainlet_denom} '{}' \
  --evm-chain-id 4242 --network-version 1 --gas ${gas_limit} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "launched ${cid_b_2} with denom B" "failed to launch ${cid_b_2}"

ensure_status "${cid_b_2}" "STATUS_ONLINE"

run_test "sscd tx escrow deposit 15000${denom} ${cid_b_2} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "deposited A into ${cid_b_2} before billing" "failed pre-billing A deposit"

sleep_billing 70

if sscd q billing get-billing-history "${cid_b_2}" > /dev/null 2>&1; then
  echo "✅ billing history fetched"
else
  echo "❌ billing history missing"; cleanup_sscd; exit 1
fi

if sscd q billing get-validator-payout-history "${cid_b_2}" > /dev/null 2>&1; then
  echo "✅ validator payout fetched"
else
  echo "❌ validator payout missing"; cleanup_sscd; exit 1
fi

echo "🔎 querying pools A & B on ${cid_b_2}"
sscd q escrow balance "${cid_b_2}" "${denom}" -o json | jq .
sscd q escrow balance "${cid_b_2}" "${denom_b}" -o json | jq .

run_test "sscd tx escrow withdraw ${cid_b_2} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "withdrew from ${cid_b_2} (owner)" "failed withdraw ${cid_b_2}"

sleep_billing 60
ensure_status "${cid_b_2}" "STATUS_OFFLINE"

run_test "sscd tx escrow deposit 5000${denom} ${cid_b_2} \
  --from alice --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "re-deposited to restart ${cid_b_2}" "failed re-deposit ${cid_b_2}"

sleep_billing 65
ensure_status "${cid_b_2}" "STATUS_ONLINE"

# ==============================
# Test 4 — Launch with A (acct1), acct2 deposits different amount, bill/payout,
#           query pool A, acct1 withdraw, query pool A
# ==============================
echo "=== Test 4: Two depositors in A, bill, withdraw, re-check pool ==="

run_test "sscd tx chainlet launch-chainlet \"\$(sscd keys show -a ${key} --keyring-backend ${keyring_backend})\" \
  ${stack_name} ${stack_ver_v1} mychaina ${chainlet_denom} '{}' \
  --evm-chain-id 5151 --network-version 1 --gas ${gas_limit} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "launched ${cid_a_2} with denom A (acct1)" "failed to launch ${cid_a_2}"

run_test "sscd tx escrow deposit 12000${denom} ${cid_a_2} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "acct1 deposited A into ${cid_a_2}" "acct1 deposit failed"

run_test "sscd tx escrow deposit 8000${denom} ${cid_a_2} \
  --from alice --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "acct2 deposited A into ${cid_a_2}" "acct2 deposit failed"

sleep_billing 70

if sscd q billing get-billing-history "${cid_a_2}" > /dev/null 2>&1; then
  echo "✅ billing history fetched"
else
  echo "❌ billing history missing"; cleanup_sscd; exit 1
fi

if sscd q billing get-validator-payout-history "${cid_a_2}" > /dev/null 2>&1; then
  echo "✅ validator payout fetched"
else
  echo "❌ validator payout missing"; cleanup_sscd; exit 1
fi

echo "🔎 pool(A) pre-withdraw on ${cid_a_2}"
sscd q escrow balance "${cid_a_2}" "${denom}" -o json | jq .

run_test "sscd tx escrow withdraw ${cid_a_2} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "acct1 withdrew from ${cid_a_2}" "acct1 withdraw failed"

echo "🔎 pool(A) post-withdraw on ${cid_a_2}"
sscd q escrow balance "${cid_a_2}" "${denom}" -o json | jq .

# ==============================
# Done
# ==============================
echo -e "\n👍 ALL PRIORITIZED-DENOM TESTS PASSED\n"
