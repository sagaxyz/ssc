#!/usr/bin/env bash
set -euo pipefail

# ==============================
# Configuration
# ==============================
tx_block_time=.5
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
# Prechecks (single-denom modes)
cid_a0="mychaina_7777-1"
cid_b0="mychainb_8888-1"
# Prioritized-denom suite
cid_a_1="mychaina_100001-1"  # Test 1
cid_b_1="mychainb_1337-1"    # Test 2
cid_b_2="mychainb_4242-1"    # Test 3
cid_a_2="mychaina_5151-1"    # Test 4

# ==============================
# Harness
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

# Resolve addresses for per-funder prints
key_addr="$(sscd keys show -a "${key}" --keyring-backend "${keyring_backend}")"
alice_addr="$(sscd keys show -a alice --keyring-backend "${keyring_backend}")"

# ==============================
# DRY print helpers (still using sscd directly)
# ==============================

# Print the single pool (if present) from a cached pools JSON
_print_pool_from_json() {
  local pools_json="$1" cid="$2" d="$3"
  echo "${pools_json}" | jq --arg d "$d" '
    (.pools // .Pools // [])
    | map(select((.denom // .Denom) == $d)) | .[0] // {}
    | if . == {} then
        "no pool yet"
      else
        {chainId: (.chainId // .ChainId), denom: (.denom // .Denom), balance: (.balance // .Balance), shares: (.shares // .Shares)}
      end
  '
}

# Print a funder position if it exists
_print_funder_if_any() {
  local cid="$1" d="$2" addr="$3"
  local json shares

  if ! json="$(sscd q escrow funder "$cid" "$d" "$addr" -o json 2>/dev/null)"; then
    echo "no position"
    return 0
  fi

  # Extract decimal string from common shapes
  shares="$(jq -r '.shares.shares' <<<"$json")"

  if [[ -z "$shares" || "$shares" == "null" ]]; then
    echo "no position"
    return 0
  fi

  # Print a consistent object with context from parameters
  jq -n --arg cid "$cid" --arg d "$d" --arg addr "$addr" --arg shares "$shares" \
    '{chainId:$cid, denom:$d, address:$addr, shares:$shares}'
}

# One-call state dump for both pools & both addresses
show_escrow_state() {
  local cid="$1"
  echo "—— Escrow State for ${cid} ——"
  local pools_json
  pools_json="$(sscd q escrow pools "${cid}" -o json)"

  echo "🔎 pool(${denom}) on ${cid}"
  _print_pool_from_json "${pools_json}" "${cid}" "${denom}"
  echo "🔎 pool(${denom_b}) on ${cid}"
  _print_pool_from_json "${pools_json}" "${cid}" "${denom_b}"

  echo "👤 funder(${key_addr}) in ${cid}/${denom}"
  _print_funder_if_any "${cid}" "${denom}" "${key_addr}"
  echo "👤 funder(${key_addr}) in ${cid}/${denom_b}"
  _print_funder_if_any "${cid}" "${denom_b}" "${key_addr}"

  echo "👤 funder(${alice_addr}) in ${cid}/${denom}"
  _print_funder_if_any "${cid}" "${denom}" "${alice_addr}"
  echo "👤 funder(${alice_addr}) in ${cid}/${denom_b}"
  _print_funder_if_any "${cid}" "${denom_b}" "${alice_addr}"
  echo
}

# ==============================
# Setup: create stack and bump version
# ==============================
echo "=== create/update stack ==="

run_test "sscd tx chainlet create-chainlet-stack ${stack_name} \"Your personal EVM\" \
  ${stack_img_v1} ${stack_ver_v1} ${stack_checksum_v1} ${epoch_fee_a} ${setup_fee_a} false \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "created ${stack_name} stack (initial A fees)" "failed to create ${stack_name}"

run_test "sscd tx chainlet update-chainlet-stack ${stack_name} \
  ${stack_img_v2} ${stack_ver_v2} ${stack_checksum_v2} false \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "updated ${stack_name} to ${stack_ver_v2}" "failed to update stack version"

# ==============================
# Precheck #1: Single-denom A mode (fees = A only)
# ==============================
echo "=== Precheck #1: single-denom A mode ==="

# Restrict fees to A only
fees_csv="${epoch_fee_a}:${setup_fee_a}"
run_test "sscd tx chainlet update-stack-fees ${stack_name} \"${fees_csv}\" \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "updated fees to A only" "failed to set A-only fees"

# Launch A0, deposit A, bill, withdraw, offline, re-deposit A, online
run_test "sscd tx chainlet launch-chainlet \"\$(sscd keys show -a ${key} --keyring-backend ${keyring_backend})\" \
  ${stack_name} ${stack_ver_v1} mychaina ${chainlet_denom} '{}' \
  --evm-chain-id 7777 --network-version 1 --gas ${gas_limit} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "launched ${cid_a0} (A-only mode)" "failed to launch ${cid_a0}"

ensure_status "${cid_a0}" "STATUS_ONLINE"

run_test "sscd tx escrow deposit 15000${denom} ${cid_a0} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "deposited A into ${cid_a0}" "failed A deposit for ${cid_a0}"
show_escrow_state "${cid_a0}"

sleep_billing 70
sscd q billing get-billing-history "${cid_a0}" > /dev/null 2>&1 && echo "✅ billing history fetched" || { echo "❌ billing history missing"; cleanup_sscd; exit 1; }
sscd q billing get-validator-payout-history "${cid_a0}" > /dev/null 2>&1 && echo "✅ validator payout fetched" || { echo "❌ validator payout missing"; cleanup_sscd; exit 1; }

run_test "sscd tx escrow withdraw ${cid_a0} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "withdrew all from ${cid_a0}" "failed withdraw in A-only mode"
show_escrow_state "${cid_a0}"

sleep_billing 60
ensure_status "${cid_a0}" "STATUS_OFFLINE"

run_test "sscd tx escrow deposit 5000${denom} ${cid_a0} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "re-deposited A to restart ${cid_a0}" "failed re-deposit A"
show_escrow_state "${cid_a0}"

sleep_billing 65
ensure_status "${cid_a0}" "STATUS_ONLINE"

# Try a B operation and expect failure in A-only mode.
# NOTE: Launch doesn’t depend on fee denoms, so we assert by attempting a B **deposit** (should be rejected).
run_test "sscd tx escrow deposit 12345${denom_b} ${cid_a0} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  nonzero "rejected B deposit while A-only" "B deposit unexpectedly accepted in A-only mode"

# ==============================
# Precheck #2: Single-denom B mode (fees = B only)
# ==============================
echo "=== Precheck #2: single-denom B mode ==="

fees_csv="${epoch_fee_b}:${setup_fee_b}"
run_test "sscd tx chainlet update-stack-fees ${stack_name} \"${fees_csv}\" \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "updated fees to B only" "failed to set B-only fees"

run_test "sscd tx chainlet launch-chainlet \"\$(sscd keys show -a ${key} --keyring-backend ${keyring_backend})\" \
  ${stack_name} ${stack_ver_v1} mychainb ${chainlet_denom} '{}' \
  --evm-chain-id 8888 --network-version 1 --gas ${gas_limit} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "launched ${cid_b0} (B-only mode)" "failed to launch ${cid_b0}"

ensure_status "${cid_b0}" "STATUS_ONLINE"

run_test "sscd tx escrow deposit 15000${denom_b} ${cid_b0} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "deposited B into ${cid_b0}" "failed B deposit for ${cid_b0}"
show_escrow_state "${cid_b0}"

sleep_billing 70
sscd q billing get-billing-history "${cid_b0}" > /dev/null 2>&1 && echo "✅ billing history fetched" || { echo "❌ billing history missing"; cleanup_sscd; exit 1; }
sscd q billing get-validator-payout-history "${cid_b0}" > /dev/null 2>&1 && echo "✅ validator payout fetched" || { echo "❌ validator payout missing"; cleanup_sscd; exit 1; }

run_test "sscd tx escrow withdraw ${cid_b0} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "withdrew all from ${cid_b0}" "failed withdraw in B-only mode"
show_escrow_state "${cid_b0}"

sleep_billing 60
ensure_status "${cid_b0}" "STATUS_OFFLINE"

run_test "sscd tx escrow deposit 5000${denom_b} ${cid_b0} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "re-deposited B to restart ${cid_b0}" "failed re-deposit B"
show_escrow_state "${cid_b0}"

sleep_billing 65
ensure_status "${cid_b0}" "STATUS_ONLINE"

# ==============================
# Enable prioritized multi-denom fees (A then B; A priority)
# ==============================
echo "=== enable prioritized (A,B) fees ==="
fees_csv="${epoch_fee_a}:${setup_fee_a},${epoch_fee_b}:${setup_fee_b}"
run_test "sscd tx chainlet update-stack-fees ${stack_name} \"${fees_csv}\" \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "updated fees (A then B; A priority)" "failed to set prioritized fees"

# ==============================
# Test 1 — Launch with A, bill, payout, withdraw, offline, deposit B, show state
# ==============================
echo "=== Test 1: A-priority launch, withdraw, restart with B ==="

run_test "sscd tx chainlet launch-chainlet \"\$(sscd keys show -a ${key} --keyring-backend ${keyring_backend})\" \
  ${stack_name} ${stack_ver_v1} mychaina ${chainlet_denom} '{}' \
  --evm-chain-id 100001 --network-version 1 --gas ${gas_limit} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "launched ${cid_a_1} with denom A" "failed to launch ${cid_a_1}"

ensure_status "${cid_a_1}" "STATUS_ONLINE"
sleep_billing 70
sscd q billing get-billing-history "${cid_a_1}" > /dev/null 2>&1 && echo "✅ billing history fetched" || { echo "❌ billing history missing"; cleanup_sscd; exit 1; }
sscd q billing get-validator-payout-history "${cid_a_1}" > /dev/null 2>&1 && echo "✅ validator payout fetched" || { echo "❌ validator payout missing"; cleanup_sscd; exit 1; }

run_test "sscd tx escrow withdraw ${cid_a_1} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "withdrew from ${cid_a_1} escrow" "failed to withdraw from ${cid_a_1}"
show_escrow_state "${cid_a_1}"

sleep_billing 60
ensure_status "${cid_a_1}" "STATUS_OFFLINE"

run_test "sscd tx escrow deposit 20000${denom_b} ${cid_a_1} \
  --from alice --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "deposited B into ${cid_a_1}" "failed to deposit B"
show_escrow_state "${cid_a_1}"

sleep_billing 65
ensure_status "${cid_a_1}" "STATUS_ONLINE"

# ==============================
# Test 2 — Launch with B, bill/payout, withdraw, offline, deposit A, show state
# ==============================
echo "=== Test 2: B launch, withdraw, restart with A ==="

run_test "sscd tx chainlet launch-chainlet \"\$(sscd keys show -a ${key} --keyring-backend ${keyring_backend})\" \
  ${stack_name} ${stack_ver_v1} mychainb ${chainlet_denom} '{}' \
  --evm-chain-id 1337 --network-version 1 --gas ${gas_limit} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "launched ${cid_b_1} with denom B" "failed to launch ${cid_b_1}"

ensure_status "${cid_b_1}" "STATUS_ONLINE"
sleep_billing 70
sscd q billing get-billing-history "${cid_b_1}" > /dev/null 2>&1 && echo "✅ billing history fetched" || { echo "❌ billing history missing"; cleanup_sscd; exit 1; }
sscd q billing get-validator-payout-history "${cid_b_1}" > /dev/null 2>&1 && echo "✅ validator payout fetched" || { echo "❌ validator payout missing"; cleanup_sscd; exit 1; }

run_test "sscd tx escrow withdraw ${cid_b_1} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "withdrew from ${cid_b_1} escrow" "failed to withdraw from ${cid_b_1}"
show_escrow_state "${cid_b_1}"

sleep_billing 60
ensure_status "${cid_b_1}" "STATUS_OFFLINE"

run_test "sscd tx escrow deposit 20000${denom} ${cid_b_1} \
  --from alice --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "deposited A into ${cid_b_1}" "failed to deposit A"
show_escrow_state "${cid_b_1}"

sleep_billing 65
ensure_status "${cid_b_1}" "STATUS_ONLINE"

# ==============================
# Test 3 — Launch with B, deposit A pre-billing, bill in A, withdraw, deposit, show state
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
show_escrow_state "${cid_b_2}"

sleep_billing 70
sscd q billing get-billing-history "${cid_b_2}" > /dev/null 2>&1 && echo "✅ billing history fetched" || { echo "❌ billing history missing"; cleanup_sscd; exit 1; }
sscd q billing get-validator-payout-history "${cid_b_2}" > /dev/null 2>&1 && echo "✅ validator payout fetched" || { echo "❌ validator payout missing"; cleanup_sscd; exit 1; }

# Show state after billing
show_escrow_state "${cid_b_2}"

run_test "sscd tx escrow withdraw ${cid_b_2} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "withdrew from ${cid_b_2} (owner)" "failed withdraw ${cid_b_2}"
show_escrow_state "${cid_b_2}"

sleep_billing 60
ensure_status "${cid_b_2}" "STATUS_OFFLINE"

run_test "sscd tx escrow deposit 5000${denom} ${cid_b_2} \
  --from alice --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "re-deposited to restart ${cid_b_2}" "failed re-deposit ${cid_b_2}"
show_escrow_state "${cid_b_2}"

sleep_billing 65
ensure_status "${cid_b_2}" "STATUS_ONLINE"

# ==============================
# Test 4 — Launch with A (acct1), acct2 deposits, bill/payout, withdraw acct1, show state
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
show_escrow_state "${cid_a_2}"

run_test "sscd tx escrow deposit 8000${denom} ${cid_a_2} \
  --from alice --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "acct2 deposited A into ${cid_a_2}" "acct2 deposit failed"
show_escrow_state "${cid_a_2}"

sleep_billing 70
sscd q billing get-billing-history "${cid_a_2}" > /dev/null 2>&1 && echo "✅ billing history fetched" || { echo "❌ billing history missing"; cleanup_sscd; exit 1; }
sscd q billing get-validator-payout-history "${cid_a_2}" > /dev/null 2>&1 && echo "✅ validator payout fetched" || { echo "❌ validator payout missing"; cleanup_sscd; exit 1; }

echo "🔎 pool(A) pre-withdraw on ${cid_a_2}"
show_escrow_state "${cid_a_2}"

run_test "sscd tx escrow withdraw ${cid_a_2} \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "acct1 withdrew from ${cid_a_2}" "acct1 withdraw failed"
show_escrow_state "${cid_a_2}"

# ==============================
# Done
# ==============================
echo -e "\n👍 ALL TESTS PASSED (single-denom prechecks + prioritized-denom suite)\n"
