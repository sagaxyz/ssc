#!/bin/bash

KEY=${KEY:-"bob"}
DENOM=${DENOM:-"utsaga"}

rm -rf tmp out
mkdir -p tmp
mkdir -p out

chainlets=$(sscd q chainlet list-chainlets --output json | jq -r '.Chainlets | .[] .chainId ')
COUNT=1
while read -r line; do
    echo $line
    
    sscd tx escrow deposit 100000$DENOM $line --from $KEY --gas 10000000 --fees 1000000$DENOM --generate-only > tmp/$COUNT.json
    COUNT=$((COUNT+1))
done <<< "$chainlets"


jq -s 'reduce .[] as $item ([]; . + $item.body.messages)' tmp/*.json > out/msgs.json
jq '.body.messages |= inputs' tmp/1.json out/msgs.json > out/tx.json

sscd tx sign out/tx.json --from $KEY > out/signed_tx.json
sscd tx broadcast out/signed_tx.json