#!/bin/bash

KEY=${KEY:-"bob"}

NUM_TXS=${NUM_TXS:-50}
CHAINLET_STACK_NAME=${CHAINLET_STACK_NAME:-"sagaevm"}
CHAINLET_STACK_VERSION=${CHAINLET_STACK_VERSION:-"1.0"}
CHAINLET_STACK_REPO=${CHAINLET_STACK_REPO:-"sagaxyz/sagaevm:gamesjam-pi14"}
CHAINLET_NAME=${CHAINLET_NAME:-"chainlet"}
CHAINLET_DENOM=${CHAINLET_DENOM:-"asaga"}
CHAINLET_PARAMS=${CHAINLET_PARAMS:-'{"genAcctBalances":"saga1mk92pa54q8ehgcdqh0qp4pj6ddjwgt25aknqxn=1000,saga18xqr6cnyezq4pudqnf53klj3ppq3mvm4eea6dp=2000"}'}
MAINTAINERS=$(sscd keys show -a $KEY)
DENOM=${DENOM:-"utsaga"}

rm -rf tmp out
mkdir -p tmp
mkdir -p out

sscd tx chainlet create-chainlet-stack $CHAINLET_STACK_NAME "Your personal EVM" $CHAINLET_STACK_REPO $CHAINLET_STACK_VERSION abc123 1000$DENOM minute 1000$DENOM --from $KEY -o json -y 

sleep 8

for (( i=1; i<=NUM_TXS; i++ ))
do
    sscd tx chainlet launch-chainlet $MAINTAINERS $CHAINLET_STACK_NAME $CHAINLET_STACK_VERSION $CHAINLET_NAME $CHAINLET_DENOM $CHAINLET_PARAMS --from $KEY --gas 10000000 --fees 1000000$DENOM --generate-only > tmp/$i.json
done

jq -s 'reduce .[] as $item ([]; . + $item.body.messages)' tmp/*.json > out/msgs.json
jq '.body.messages |= inputs' tmp/1.json out/msgs.json > out/tx.json

sscd tx sign out/tx.json --from $KEY > out/signed_tx.json
sscd tx broadcast out/signed_tx.json
