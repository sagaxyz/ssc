#!/bin/bash

# run `ignite chain serve` before running this script

TXBLOCKTIME=5
KEY=bob
KEYRING_BACKEND=test
DENOM=utsaga
GAS_LIMIT=500000
FEES=5000$DENOM

WaitTx() {
  HASH=$1
  echo "waiting $TXBLOCKTIME seconds for the tx $HASH to be broadcasted..."
  sleep $TXBLOCKTIME
}

echo "testing create-chainlet-stack"

TX_HASH=$(sscd tx chainlet create-chainlet-stack sagaevm "Your personal EVM" sagaxyz/sagaevm:gamesjam-pi14 1.0.0 abc123 1000$DENOM minute 1000$DENOM --from $KEY --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if [ "$(echo $TX_RES | jq .code)" -eq 0 ]; then
    echo "pass: created a 'sagaevm' chainlet stack"
else
    echo "fail: failed to create a 'sagaevm' chainlet stack"
    exit 1
fi

TX_HASH=$(sscd tx chainlet create-chainlet-stack sagaevm "Your personal EVM" sagaxyz/sagaevm:2.0.0 2.0.0 def456 1000$DENOM minute 1000$DENOM --from $KEY --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if [ "$(echo $TX_RES | jq .code)" -eq 1 ]; then
    echo "pass: rejected duplicate 'sagaevm' chainlet stack name"
else
    echo "fail: accepted duplicate 'sagaevm' chainlet stack name"
	exit 1
fi

TX_HASH=$(sscd tx chainlet create-chainlet-stack sagavm "Your personal EVM" sagaxyz/sagavm:1.0.0 1.0.0 123123 1000$DENOM minute 1000$DENOM --from alice --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if [ "$(echo $TX_RES | jq .code)" -eq 1 ]; then
    echo "pass: rejected non-admin creation of a chainlet stack"
else
    echo "fail: accepted a non-admin chainlet stack creation"
	exit 1
fi

echo "testing update-chainlet-stack"

TX_HASH=$(sscd tx chainlet update-chainlet-stack sagaevm sagaxyz/sagaevm:5ed0edf 1.2.3 sha256:b4cfab4354a11805b0b60cc52f43bd5d6b41f8a291c724fc8aabfa1d5a836aed --from $KEY --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if [ "$(echo $TX_RES | jq .code)" -eq 0 ]; then
    echo "pass: successfully updated new chainlet stack version"
else
    echo "fail: failed to update chainlet stack version"
	exit 1
fi

TX_HASH=$(sscd tx chainlet update-chainlet-stack sagaevm2 sagaxyz/sagaevm:3.0.0 3.0.0 def456 --from $KEY --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if [ "$(echo $TX_RES | jq .code)" -eq 1 ]; then
    echo "pass: failed to update non-existent chainlet stack"
else
    echo "fail: accepted non-existent chainlet stack update"
	exit 1
fi

echo "testing launch-chainlet"

TX_HASH=$(sscd tx chainlet launch-chainlet "$(sscd keys show -a $KEY)" sagaevm 1.0.0 mychain asaga '{}' --evm-chain-id 100001 --network-version 1 --gas $GAS_LIMIT --from $KEY --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if [ "$(echo $TX_RES | jq .code)" -eq 0 ]; then
    echo "pass: launched a chainlet from a valid chainlet stack"
else
    echo "fail: failed to launch a chainlet from a valid chainlet stack"
	exit 1
fi



TX_HASH=$(sscd tx chainlet launch-chainlet "$(sscd keys show -a $KEY)" sagaevm 2.0.0 mychainabc asaga '{}' --evm-chain-id 100001 --network-version 1 --gas $GAS_LIMIT --from $KEY --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if [ "$(echo $TX_RES | jq .code)" -eq 0 ]; then
    echo "fail: launched a chainlet with an invalid chainlet stack"
    exit 1
fi
echo "pass: did not launch a chainlet with an invalid chainlet stack"

TX_HASH=$(sscd tx chainlet launch-chainlet "$(sscd keys show -a $KEY)" sagaevm 1.0.0 mychain asaga '{}' --evm-chain-id 13371337 --network-version 1 --gas $GAS_LIMIT --from $KEY --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if [ "$(echo $TX_RES | jq .code)" -eq 0 ]; then
    echo "pass: launched another chainlet from a valid chainlet stack"
else
    echo "fail: failed to launch a chainlet from a valid chainlet stack"
	exit 1
fi

TX_HASH=$(sscd tx chainlet launch-chainlet "$(sscd keys show -a $KEY)" sagaevm 1.0.0 kukkoo asaga '{"gasLimit":10000000,"genAcctBalances":"saga1mk92pa54q8ehgcdqh0qp4pj6ddjwgt25aknqxn=1000,saga18xqr6cnyezq4pudqnf53klj3ppq3mvm4eea6dp=100000"}' --gas $GAS_LIMIT --from $KEY --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if [ "$(echo $TX_RES | jq .code)" -eq 0 ]; then
    echo "pass: launched a chainlet with chainlet parametes, including genesis account balances"
else
    echo "fail: failed to launch a chainlet with chainlet parametes, including genesis account balances"
	exit 1
fi

TX_HASH=$(sscd tx chainlet launch-chainlet "$(sscd keys show -a alice --keyring-backend $KEYRING_BACKEND)" sagaevm 1.0.0 mychain '{}' --evm-chain-id 515151 --network-version 1 --gas $GAS_LIMIT --service-chainlet --from alice --keyring-backend $KEYRING_BACKEND --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if echo $TX_RES | jq .code | grep -q 1; then
    echo "pass: rejected non-admin launch of a service chainlet"
else
    echo "fail: accepted a non-admin launch of a service chainlet"
	exit 1
fi

TX_HASH=$(sscd tx chainlet launch-chainlet "$(sscd keys show -a $KEY --keyring-backend $KEYRING_BACKEND)" sagaevm 1.0.0 mychain '{}' --evm-chain-id 424242 --network-version 1 --gas $GAS_LIMIT --service-chainlet --from $KEY --keyring-backend $KEYRING_BACKEND --fees $FEES -o json -y | jq -r .txhash)
WaitTx $TX_HASH
TX_RES=$(sscd q tx $TX_HASH -o json)

if echo $TX_RES | jq .code | grep -q 0; then
    echo "pass: launched a service chainlet from admin account"
else
    echo "fail: failed to a service chainlet from admin account"
	exit 1
fi

printf "\n\ntesting epoch queries"
sscd q epochs epoch-infos

echo "testing chainlet stack queries"
if [ "$(sscd q chainlet list-chainlet-stack -o json | jq '.ChainletStacks | length')" -eq 1 ]; then
    echo "pass: found one chainlet stack"
else
    echo "fail: did not find one chainlet stack"
	exit 1
fi

if [ "$(sscd q chainlet get-chainlet-stack sagaevm -o json | jq '.ChainletStack.versions | length')" -eq 2 ]; then
    echo "pass: found two chainlet stack versions for 'sagaevm'"
else
    echo "fail: did not find two chainlet stack versions for 'sagaevm'"
	exit 1
fi

echo "testing chainlet queries"
if [ "$(sscd q chainlet list-chainlets -o json | jq '.Chainlets | length')" -eq 3 ]; then
    echo "pass: found three chainlets"
else
    echo "fail: did not find three chainlets"
	exit 1
fi

sscd q chainlet list-chainlets


echo
echo
echo "Sleeping for 1 minute to enable chainlets to be billed and validators to be paid out"
echo
sleep 60

echo "testing billing history retrieval"
sscd q billing get-billing-history mychain_100001-1 1> /dev/null 2>&1
RETCODE=$?
if [ $RETCODE -eq 0 ]; then
    echo "pass: fetched chainlet billing history"
else
    echo "fail: could not fetch chainlet billing history"
    exit 1
fi

echo "testing validator payout history retrieval"
sscd q billing get-validator-payout-history "$(sscd keys show -a $KEY)" 1> /dev/null 2>&1
if [ $RETCODE -eq 0 ]; then
    echo "pass: fetched validator payout history"
else
    echo "fail: could not fetch validator payout history"
    exit 1
fi

echo
echo


echo "lgtm"

