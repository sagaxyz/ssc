#!/usr/bin/env bash
set -euo pipefail

# Configuration
tx_block_time=.2
key="bob"
keyring_backend="test"
denom="utsaga"
chainlet_denom="asaga"
gas_limit=500000
fees="5000${denom}"

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

# === Testing create-chainlet-stack ===
echo "=== create-chainlet-stack ==="

run_test "sscd tx chainlet create-chainlet-stack sagaevm \"Your personal EVM\" \
  sagaxyz/sagaevm:0.7.0 0.7.0 abc123 1000${denom} minute 1000${denom} \
  --from ${key} --fees ${fees} -o json -y" \
  0 "created sagaevm chainlet stack" "failed to create sagaevm stack"

run_test "sscd tx chainlet create-chainlet-stack sagaevm \"Your personal EVM\" \
  sagaxyz/sagaevm:2.0.0 2.0.0 def456 1000${denom} minute 1000${denom} \
  --from ${key} --fees ${fees} -o json -y" \
  nonzero "rejected duplicate sagaevm stack name" "accepted duplicate sagaevm name"

run_test "sscd tx chainlet create-chainlet-stack sagavm \"Your personal EVM\" \
  sagaxyz/sagavm:1.0.0 1.0.0 123123 1000${denom} minute 1000${denom} \
  --from alice --fees ${fees} -o json -y" \
  nonzero "rejected non-admin creation" "accepted non-admin creation"

# === Testing update-chainlet-stack ===
echo "=== update-chainlet-stack ==="

run_test "sscd tx chainlet update-chainlet-stack sagaevm \
  sagaxyz/sagaevm:0.8.0 0.8.0 abc234 \
  --from ${key} --fees ${fees} -o json -y" \
  0 "updated chainlet stack version" "failed to update stack version"

run_test "sscd tx chainlet update-chainlet-stack sagaevm2 \
  sagaxyz/sagaevm:3.0.0 3.0.0 def456 \
  --from ${key} --fees ${fees} -o json -y" \
  nonzero "rejected update of non-existent stack" "accepted update of non-existent stack"

# === Testing launch-chainlet ===
echo "=== launch-chainlet ==="

run_test "sscd tx chainlet launch-chainlet \"\$(sscd keys show -a ${key})\" sagaevm 0.7.0 mychain ${chainlet_denom} '{}' \
  --evm-chain-id 100001 --network-version 1 --gas ${gas_limit} \
  --from ${key} --fees ${fees} -o json -y" \
  0 "launched chainlet from valid stack" "failed to launch chainlet from valid stack"

run_test "sscd tx chainlet launch-chainlet \"\$(sscd keys show -a ${key})\" sagavm 2.0.0 mychainabc ${chainlet_denom} '{}' \
  --evm-chain-id 100001 --network-version 1 --gas ${gas_limit} \
  --from ${key} --fees ${fees} -o json -y" \
  nonzero "rejected launch with invalid stack" "accepted invalid stack launch"

run_test "sscd tx chainlet launch-chainlet \"\$(sscd keys show -a ${key})\" sagaevm 0.7.0 mychain ${chainlet_denom} '{}' \
  --evm-chain-id 13371337 --network-version 1 --gas ${gas_limit} \
  --from ${key} --fees ${fees} -o json -y" \
  0 "launched another chainlet from valid stack" "failed second valid launch"

run_test "sscd tx chainlet launch-chainlet \"\$(sscd keys show -a ${key})\" sagaevm 0.8.0 kukkoo ${chainlet_denom} '{\"gasLimit\":10000000,\"genAcctBalances\":\"saga1mk92pa54q8ehgcdqh0qp4pj6ddjwgt25aknqxn=1000,saga18xqr6cnyezq4pudqnf53klj3ppq3mvm4eea6dp=100000\"}' \
  --gas ${gas_limit} --from ${key} --fees ${fees} -o json -y" \
  0 "launched with custom params" "failed launch with custom params"

run_test "sscd tx chainlet launch-chainlet \"\$(sscd keys show -a alice --keyring-backend ${keyring_backend})\" sagaevm 0.7.0 mychain ${chainlet_denom} '{}' \
  --evm-chain-id 515151 --network-version 1 --gas ${gas_limit} --service-chainlet \
  --from alice --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  nonzero "rejected non-admin service launch" "accepted non-admin service launch"

run_test "sscd tx chainlet launch-chainlet \"\$(sscd keys show -a ${key} --keyring-backend ${keyring_backend})\" sagaevm 0.7.0 mychain ${chainlet_denom} '{}' \
  --evm-chain-id 424242 --network-version 1 --gas ${gas_limit} --service-chainlet \
  --from ${key} --keyring-backend ${keyring_backend} --fees ${fees} -o json -y" \
  0 "launched service chainlet from admin" "failed admin service launch"

# === Queries & billing ===
echo "=== queries & billing ==="
sscd q epochs epoch-infos

# chainlet-stack list & get
if sscd q chainlet list-chainlet-stacks -o json | jq '.ChainletStacks | length' | grep -q '^1$'; then
  echo "✅ 1 stack"
else
  cleanup_sscd
  exit 1
fi

if sscd q chainlet get-chainlet-stack sagaevm -o json | jq '.ChainletStack.versions | length' | grep -q '^2$'; then
  echo "✅ 2 versions"
else
  cleanup_sscd
  exit 1
fi

# list chainlets
if sscd q chainlet list-chainlets -o json | jq '.Chainlets | length' | grep -q '^[34]$'; then
  echo "✅ chainlet count OK"
else
  cleanup_sscd
  exit 1
fi

sscd q chainlet list-chainlets

echo "🛌 Sleeping 1m for billing..."
sleep 60

# billing & payouts
if sscd q billing get-billing-history mychain_100001-1 > /dev/null 2>&1; then
  echo "✅ billing history  fetched"
else
  cleanup_sscd
  exit 1
fi

if sscd q billing get-validator-payout-history mychain_100001-1 > /dev/null 2>&1; then
  echo "✅ validator payout fetched"
else
  cleanup_sscd
  exit 1
fi

# Final cleanup
echo -e "
👍 ALL TESTS PASSED
"
cleanup_sscd
