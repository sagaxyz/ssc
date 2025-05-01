#!/bin/bash

# run `ignite chain serve` before running this script

TXBLOCKTIME=5
KEY=bob
ANOTHER_KEY=alice
GAS_LIMIT=500000
DENOM=utsaga
FEES=5000$DENOM
CHAINLET_DENOM=asaga

WaitTx() {
  HASH=$1
  echo "waiting $TXBLOCKTIME seconds for the tx $HASH to be broadcasted..."
  sleep $TXBLOCKTIME
}

echo "testing create-chainlet-stack"

TX_HASH=$(sscd tx chainlet create-chainlet-stack sagaevm "Your personal EVM" sagaxyz/sagaevm:8.2.8 1.0.0 abc123 1000utsaga minute 1000utsaga --from $KEY --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if [ "$(echo $TX_RES | jq .code)" -eq 0 ]; then
    echo "pass: created a 'sagaevm' chainlet stack"
else
    echo "fail: failed to create a 'sagaevm' chainlet stack"
    exit 1
fi

echo "testing launch-chainlet"

TX_HASH=$(sscd tx chainlet launch-chainlet "$(sscd keys show -a $KEY)" sagaevm 1.0.0 mychain $CHAINLET_DENOM '{}' --evm-chain-id 100001 --network-version 1 --gas $GAS_LIMIT --from $KEY --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if [ "$(echo $TX_RES | jq .code)" -eq 0 ]; then
    echo "pass: launched a chainlet from a valid chainlet stack"
else
    echo "fail: failed to launch a chainlet from a valid chainlet stack"
	exit 1
fi

if sscd q chainlet get-chainlet mychain_100001-1 -o json | jq '.Chainlet.status' | grep -q 'STATUS_ONLINE'; then
    echo "pass: chainlet status is STATUS_ONLINE"
else
    echo "fail: chainlet has incorrect status"
    exit 1
fi

echo "testing escrow"
sleep 5

sscd q bank balances "$(sscd keys show -a $ANOTHER_KEY)"

TX_HASH=$(sscd tx escrow deposit 30000utsaga mychain_100001-1 --from $ANOTHER_KEY --fees $FEES -y)
WaitTx $TX_HASH

TX_HASH=$(sscd tx escrow deposit 30000utsaga mychain_100001-1 --from $ANOTHER_KEY --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if [ "$(echo $TX_RES | jq .code)" -eq 0 ]; then
    echo "pass: deposited funds into escrow"
else
    echo "fail: failed to deposit funds into escrow"
	exit 1
fi

echo "testing withdraw"

TX_HASH=$(sscd tx escrow withdraw mychain_100001-1 --from $KEY --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if [ "$(echo $TX_RES | jq .code)" -eq 0 ]; then
    echo "pass: withdraw successfull"
else
    echo "fail: unable to withdraw"
	exit 1
fi

sleep 60

if sscd q chainlet get-chainlet mychain_100001-1 -o json | jq '.Chainlet.status' | grep -q 'STATUS_ONLINE'; then
    echo "pass: chainlet status is STATUS_ONLINE"
else
    echo "fail: chainlet has incorrect status"
    exit 1
fi

echo "second account withdraws as well"

TX_HASH=$(sscd tx escrow withdraw mychain_100001-1 --from $ANOTHER_KEY --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if [ "$(echo $TX_RES | jq .code)" -eq 0 ]; then
    echo "pass: withdraw successfull"
else
    echo "fail: unable to withdraw"
	exit 1
fi

sleep 60

if sscd q chainlet get-chainlet mychain_100001-1 -o json | jq '.Chainlet.status' | grep -q 'STATUS_OFFLINE'; then
    echo "pass: chainlet status is STATUS_OFFLINE"
else
    echo "fail: chainlet has incorrect status"
    exit 1
fi

echo "testing chainlet restarts"

TX_HASH=$(sscd tx escrow deposit 29999utsaga mychain_100001-1 --from $ANOTHER_KEY --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if [ "$(echo $TX_RES | jq .code)" -eq 0 ]; then
echo "pass: deposited funds into escrow"
else
    echo "fail: failed to deposit funds into escrow"
	exit 1
fi

sleep 75
sscd q chainlet get-chainlet mychain_100001-1
if sscd q chainlet get-chainlet mychain_100001-1 -o json | jq '.Chainlet.status' | grep -q 'STATUS_OFFLINE'; then
    echo "pass: chainlet status is STATUS_OFFLINE"
else
    echo "fail: chainlet has incorrect status"
    exit 1
fi

TX_HASH=$(sscd tx escrow deposit 1utsaga mychain_100001-1 --from $KEY --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if [ "$(echo $TX_RES | jq .code)" -eq 0 ]; then
echo "pass: deposited funds into escrow"
else
    echo "fail: failed to deposit funds into escrow"
	exit 1
fi

sleep 60

if sscd q chainlet get-chainlet mychain_100001-1 -o json | jq '.Chainlet.status' | grep -q 'STATUS_ONLINE'; then
    echo "pass: chainlet status is STATUS_ONLINE"
else
    echo "fail: chainlet has incorrect status"
    exit 1
fi
