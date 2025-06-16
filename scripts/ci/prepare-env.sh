#!/bin/bash
DIR=~/.ssc
SSC_HOME="--home $DIR"
KEYRING="--keyring-backend=test"

function log() {
  local msg=$1

  echo "$(date) $msg" >&2
}
function fail() {
  if [ $# -gt 0 ]; then
    local msg=$1
    log "$msg"
  fi

  exit 1
}


sscd init test --chain-id testchain $SSC_HOME || fail "failed to init configuration"
cp ./scripts/ci/config/client.toml $DIR/config/ || fail "failed to copy client.toml"
sscd keys add alice $SSC_HOME $KEYRING || fail "failed to add key alice"
sscd keys add bob $SSC_HOME $KEYRING || fail "failed to add key bob"
sscd add-genesis-account "$(sscd keys show alice -a $SSC_HOME $KEYRING)" 100000000000000000000000000utsaga,100000000stake $SSC_HOME || fail "failed to add genesis account for alice"
sscd add-genesis-account "$(sscd keys show bob -a $SSC_HOME $KEYRING)" 100000000000000000000000000utsaga,100000000stake $SSC_HOME || fail "failed to add genesis account for bob"
jq '.app_state["chainlet"]["params"]["chainletStackProtections"]=true' $DIR/config/genesis.json > $DIR/config/tmp_genesis.json && mv $DIR/config/tmp_genesis.json $DIR/config/genesis.json || fail "failed to set chainletStackProtections"
jq '.app_state["chainlet"]["params"]["nEpochDeposit"]="30"' $DIR/config/genesis.json > $DIR/config/tmp_genesis.json && mv $DIR/config/tmp_genesis.json $DIR/config/genesis.json || fail "failed to set nEpochDeposit"
jq '.app_state["acl"]["params"]["enable"]=true' > $DIR/config/tmp_genesis.json $DIR/config/genesis.json && mv $DIR/config/tmp_genesis.json $DIR/config/genesis.json || fail "failed to enable ACL"
jq ".app_state[\"acl\"][\"allowed\"]=[\"$(sscd keys show bob -a $SSC_HOME $KEYRING)\"]" $DIR/config/genesis.json > $DIR/config/tmp_genesis.json && mv $DIR/config/tmp_genesis.json $DIR/config/genesis.json || fail "failed to set allowed ACL address"
jq ".app_state[\"acl\"][\"admins\"]=[\"$(sscd keys show bob -a $SSC_HOME $KEYRING)\"]" $DIR/config/genesis.json > $DIR/config/tmp_genesis.json && mv $DIR/config/tmp_genesis.json $DIR/config/genesis.json || fail "failed to set ACL admin address"
jq '.app_state["gov"]["params"]["max_deposit_period"]="600s"' > $DIR/config/tmp_genesis.json $DIR/config/genesis.json && mv $DIR/config/tmp_genesis.json $DIR/config/genesis.json || fail "failed to set max_deposit_period"
jq '.app_state["gov"]["params"]["voting_period"]="600s"' > $DIR/config/tmp_genesis.json $DIR/config/genesis.json && mv $DIR/config/tmp_genesis.json $DIR/config/genesis.json || fail "failed to set voting_period"
jq '.app_state["gov"]["params"]["expedited_voting_period"]="60s"' > $DIR/config/tmp_genesis.json $DIR/config/genesis.json && mv $DIR/config/tmp_genesis.json $DIR/config/genesis.json || fail "failed to set expedited_voting_period"
sscd genesis gentx alice 100000000stake --chain-id testchain $SSC_HOME $KEYRING || fail "failed to create genesis transaction for alice"
sscd genesis collect-gentxs $SSC_HOME || fail "failed to collect genesis transactions"
sscd start $SSC_HOME &
sleep 10
