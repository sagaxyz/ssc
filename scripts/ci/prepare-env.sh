#!/bin/bash
set -euo pipefail

DIR=~/.ssc
SSC_HOME="--home $DIR"
KEYRING="--keyring-backend=test"
GENESIS="$DIR/config/genesis.json"
TMP="$DIR/config/tmp_genesis.json"

log() { echo "$(date) $1" >&2; }
fail() { log "$1"; exit 1; }

update() {
  jq "$1" "$GENESIS" > "$TMP" && mv "$TMP" "$GENESIS" || fail "failed: $1"
}

sscd init test --chain-id testchain $SSC_HOME || fail "init failed"
cp ./scripts/ci/config/client.toml "$DIR/config/" || fail "copy client.toml"

sscd keys add alice $SSC_HOME $KEYRING || fail "add alice"
sscd keys add bob $SSC_HOME $KEYRING || fail "add bob"

sscd add-genesis-account "$(sscd keys show alice -a $SSC_HOME $KEYRING)" 100000000000000000000000000utsaga,100000000stake $SSC_HOME || fail "add alice acct"
sscd add-genesis-account "$(sscd keys show bob -a $SSC_HOME $KEYRING)" 100000000000000000000000000utsaga,100000000stake $SSC_HOME || fail "add bob acct"

BOB_ADDR=$(sscd keys show bob -a $SSC_HOME $KEYRING)

update '.app_state["chainlet"]["params"]["chainletStackProtections"]=true'
update '.app_state["chainlet"]["params"]["nEpochDeposit"]="30"'
update '.app_state["acl"]["params"]["enable"]=true'
update ".app_state[\"acl\"][\"allowed\"]=[\"$BOB_ADDR\"]"
update ".app_state[\"acl\"][\"admins\"]=[\"$BOB_ADDR\"]"
update '.app_state["gov"]["params"]["max_deposit_period"]="600s"'
update '.app_state["gov"]["params"]["voting_period"]="600s"'
update '.app_state["gov"]["params"]["expedited_voting_period"]="60s"'

sscd genesis gentx alice 100000000stake --chain-id testchain $SSC_HOME $KEYRING || fail "gentx failed"
sscd genesis collect-gentxs $SSC_HOME || fail "collect-gentxs failed"

sscd start $SSC_HOME &
sleep 10
