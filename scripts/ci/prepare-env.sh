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

sed -i.bak 's/^timeout_propose *=.*/timeout_propose = "500ms"/' "$DIR/config/config.toml" && rm -f "$DIR/config/config.toml.bak" || fail "failed to set timeout_propose"
sed -i.bak 's/^timeout_propose_delta *=.*/timeout_propose_delta = "500ms"/' "$DIR/config/config.toml" && rm -f "$DIR/config/config.toml.bak" || fail "failed to set timeout_propose_delta"
sed -i.bak 's/^timeout_prevote *=.*/timeout_prevote = "500ms"/' "$DIR/config/config.toml" && rm -f "$DIR/config/config.toml.bak" || fail "failed to set timeout_prevote"
sed -i.bak 's/^timeout_prevote_delta *=.*/timeout_prevote_delta = "500ms"/' "$DIR/config/config.toml" && rm -f "$DIR/config/config.toml.bak" || fail "failed to set timeout_prevote_delta"
sed -i.bak 's/^timeout_precommit *=.*/timeout_precommit = "500ms"/' "$DIR/config/config.toml" && rm -f "$DIR/config/config.toml.bak" || fail "failed to set timeout_precommit"
sed -i.bak 's/^timeout_precommit_delta *=.*/timeout_precommit_delta = "500ms"/' "$DIR/config/config.toml" && rm -f "$DIR/config/config.toml.bak" || fail "failed to set timeout_precommit_delta"
sed -i.bak 's/^timeout_commit *=.*/timeout_commit = "200ms"/' "$DIR/config/config.toml" && rm -f "$DIR/config/config.toml.bak" || fail "failed to set timeout_commit"

sscd genesis gentx alice 100000000stake --chain-id testchain $SSC_HOME $KEYRING || fail "failed to create genesis transaction for alice"
sscd genesis collect-gentxs $SSC_HOME || fail "failed to collect genesis transactions"

sscd start $SSC_HOME &
sleep 10
