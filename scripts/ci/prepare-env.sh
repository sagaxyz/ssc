#!/bin/bash

sscd init test --chain-id testchain
cp ./scripts/ci/config/client.toml ~/.ssc/config/
sscd keys add alice
sscd keys add bob
sscd add-genesis-account "$(sscd keys show alice -a)" 100000000000000000000000000utsaga,100000000stake
sscd add-genesis-account "$(sscd keys show bob -a)" 100000000000000000000000000utsaga,100000000stake
jq '.app_state["chainlet"]["params"]["chainletStackProtections"]=true' ~/.ssc/config/genesis.json > ~/.ssc/config/tmp_genesis.json && mv ~/.ssc/config/tmp_genesis.json ~/.ssc/config/genesis.json
jq '.app_state["chainlet"]["params"]["nEpochDeposit"]="30"' ~/.ssc/config/genesis.json > ~/.ssc/config/tmp_genesis.json && mv ~/.ssc/config/tmp_genesis.json ~/.ssc/config/genesis.json
jq '.app_state["acl"]["params"]["enable"]=true' > ~/.ssc/config/tmp_genesis.json ~/.ssc/config/genesis.json && mv ~/.ssc/config/tmp_genesis.json ~/.ssc/config/genesis.json
jq ".app_state[\"acl\"][\"allowed\"]=[{\"format\":1,\"value\":\"$(sscd keys show bob -a)\"}]" ~/.ssc/config/genesis.json > ~/.ssc/config/tmp_genesis.json && mv ~/.ssc/config/tmp_genesis.json ~/.ssc/config/genesis.json
sscd gentx alice 100000000stake --chain-id testchain
sscd collect-gentxs
sscd start &
sleep 10
