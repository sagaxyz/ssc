#!/bin/bash
set -u

# Expand globs to nothing
shopt -s nullglob

# Validate dependencies are installed
command -v jq &> /dev/null || fail "jq not installed"
command -v curl &> /dev/null || fail "curl not installed"

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

function check_env_vars() {
  for name in "$@"; do
    local value="${!name:-}"
    if [ -z "$value" ]; then
      echo "Variable $name is empty"
      return 1
    fi
  done
}

# Print all env variables for debugging
env | sed -e 's/^.*MNEMONIC.*$/<removed>/' -e 's/^.*KEY.*$/<removed>/'

# Optional env vars
LOGLEVEL=${LOGLEVEL:-"info"}
OPTS=${OPTS:-""}
EXTERNAL_ADDRESS=${EXTERNAL_ADDRESS:-""}
PEERS=${PEERS:-""}
VALIDATOR_KEY=${VALIDATOR_KEY:-""}
NODE_KEY=${NODE_KEY:-""}
MNEMONIC=${MNEMONIC:-""}

# Check that mandatory env variables are not empty
check_env_vars CHAIN_ID MONIKER DENOM GENESIS || fail

# Constant vars
CONFIG_DIR="$HOME/.ssc/config"

# Note that we are always recreating the config directory in order to avoid handling migrations
log "initializing config directory"
#TODO Remove the entire directory when launched from the controller.
#     It is not used now only because without the controller we cannot
#     recover the node key.
#rm -rf "$CONFIG_DIR"
rm -rf "$CONFIG_DIR/*toml"
#TODO Remove --recover variant with the controller being used.
if [ -n "$MNEMONIC" ]; then
  echo "$MNEMONIC" | sscd init "$MONIKER" --chain-id "$CHAIN_ID" --default-denom "$DENOM" --recover --overwrite || fail "failed to init configuration"
else
  sscd init "$MONIKER" --chain-id "$CHAIN_ID" --default-denom "$DENOM" --overwrite || fail "failed to init configuration"
fi

# Overwrite the randomly generated validator private key
if [ -n "$VALIDATOR_KEY" ]; then
  val_key_json=$(echo "$VALIDATOR_KEY" | base64 -d) || fail "failed to decode validator key"
  echo "$val_key_json" > $CONFIG_DIR/priv_validator_key.json || fail
fi
# Overwrite the randomly generated node private key
if [ -n "$NODE_KEY" ]; then
  node_key_json=$(echo "$NODE_KEY" | base64 -d) || fail "failed to decode node key"
  echo "$node_key_json" > $CONFIG_DIR/node_key.json || fail
fi

# Overwrite the blank genesis file created by the init command
if [[ ! $GENESIS =~ ^http ]]; then
  genesis_json=$(echo "$GENESIS" | base64 -d) || fail
  echo "$genesis_json" > $CONFIG_DIR/genesis.json || fail
else
  curl "$GENESIS" --output $CONFIG_DIR/genesis.json -f || fail
fi
sscd genesis validate || fail

peers=$PEERS
if [ -z "$peers" ]; then
  peers=$(jq -r '.app_state.genutil.gen_txs[].body.memo' $CONFIG_DIR/genesis.json | grep -v "$EXTERNAL_ADDRESS" | paste -sd, -)
  log "extracted peers from the genesis file: $peers"
fi

# Node-specific configuration
# XXX: Do not add things here that are not related to the functionality of this script
# or are a requirement to start nodes. Those should be passed as a simple command line
# argument in the $OPTS env variable.
log "configuring node"
sed -i "s/^minimum-gas-prices =/ s/= .*/= \"0.01$DENOM,0.01stake\"/g" $CONFIG_DIR/app.toml
sed -i "s/^log_level =.*/log_level = \"$LOGLEVEL\"/g" $CONFIG_DIR/config.toml
sed -i 's/^create_empty_blocks = true/create_empty_blocks = false/g' $CONFIG_DIR/config.toml
sed -i "s/^external_address =.*/external_address = \"$EXTERNAL_ADDRESS\"/g" $CONFIG_DIR/config.toml
sed -i "s/^persistent_peers =.*/persistent_peers = \"$peers\"/g" $CONFIG_DIR/config.toml

#TODO replace constant changes here with default values in the Go code
sed -i 's/^address = .*:9090"/address = "0.0.0.0:9090"/g' $CONFIG_DIR/app.toml #grpc.address
sed -i 's/^laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' $CONFIG_DIR/config.toml
sed -i 's/^allow_duplicate_ip = false/allow_duplicate_ip = true/g' $CONFIG_DIR/config.toml
sed -i 's/^send_rate = 5120000/send_rate = 20000000/g' $CONFIG_DIR/config.toml
sed -i 's/^recv_rate = 5120000/recv_rate = 20000000/g' $CONFIG_DIR/config.toml
sed -i 's/^max_packet_msg_payload_size =.*/max_packet_msg_payload_size = 10240/g' $CONFIG_DIR/config.toml
sed -i 's/^flush_throttle_timeout = \"100ms\"/flush_throttle_timeout = \"10ms\"/g' $CONFIG_DIR/config.toml
sed -i 's/^iavl-cache-size = .*/iavl-cache-size = 100000/g' $CONFIG_DIR/app.toml
sed -i 's/^cors_allowed_origins = .*/cors_allowed_origins = ["*"]/g' $CONFIG_DIR/config.toml
sed -i 's/^enabled-unsafe-cors =.*/enabled-unsafe-cors = true/g' $CONFIG_DIR/app.toml
sed -i 's/^enable-unsafe-cors =.*/enable-unsafe-cors = true/g' $CONFIG_DIR/app.toml
sed -i 's/^addr_book_strict = true/addr_book_strict = false/g' $CONFIG_DIR/config.toml
sed -i 's/prometheus = false/prometheus = true/g' $CONFIG_DIR/config.toml
echo 'json-log-file = "/var/log/saga/ssc.log"' >> $CONFIG_DIR/config.toml

exec sscd start $OPTS
