#!/bin/bash

# create a chainlet stack and launch a chainlet
sscd tx chainlet create-chainlet-stack sagaevm "Your personal EVM" sagaxyz/sagaevm:gamesjam-pi14 1.0.0 abc123 1utsaga minute 1utsaga --from bob --fees 5000utsaga -y
sscd tx chainlet launch-chainlet "$(sscd keys show -a bob)" sagaevm 1.0.0 mychain asaga '{}' --evm-chain-id 100001 --fees 5000utsaga --from bob -y

# check the chainlet escrow account
sscd q escrow get-chainlet-account mychain_100001-1
sscd q chainlet get-chainlet mychain_100001-1

# deposit into the escrow account from a different user
sscd tx escrow deposit 2utsaga mychain_100001-1 --from alice --fees 5000utsaga -y
sscd q escrow get-chainlet-account mychain_100001-1
sscd q chainlet get-chainlet mychain_100001-1

# withdraw from the escrow account
sscd tx escrow withdraw mychain_100001-1 --from bob --fees 5000utsaga -y
sscd q escrow get-chainlet-account mychain_100001-1
sscd q escrow get-chainlet-account mychain_100001-1

sscd q escrow get-chainlet-account mychain_100001-1
