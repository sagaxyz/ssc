package utils

import (
	"context"
	"fmt"

	"github.com/strangelove-ventures/interchaintest/v8/ibc"
)

type CreateStackParams struct {
	Name        string // "sagaevm"
	Description string // "Your personal EVM"
	Image       string // "sagaxyz/sagaevm:0.7.0"
	Version     string // "0.7.0"
	Hash        string // "abc123"
	MinDeposit  string // e.g. "1000utsaga"
	MinTopup    string // e.g. "1000utsaga"
	CcvConsumer bool   // false
}

type UpdateStackParams struct {
	Name        string // "sagaevm"
	Image       string // "sagaxyz/sagaevm:0.8.0"
	Version     string // "0.8.0"
	Hash        string // "abc234"
	CcvConsumer bool   // false
}

type LaunchChainletParams struct {
	OwnerAddr     string // bob.FormattedAddress()
	StackName     string // "sagaevm"
	StackVersion  string // "0.7.0"
	ChainletID    string // "mychain"
	ChainletDenom string // "asaga"
	CustomJSON    string // "{}" or {"gasLimit":...}
	// flags
	EVMChainID     string // "100001"
	NetworkVersion string // "1"
	Gas            string // "500000"
	Service        bool   // optional flag --service-chainlet
}

func ChainletCreateStack(ctx context.Context, chain ibc.Chain, signer ibc.Wallet, fees string, p CreateStackParams) (txhash string, code uint32, raw string, err error) {
	args := []string{
		"chainlet", "create-chainlet-stack",
		p.Name, p.Description, p.Image, p.Version, p.Hash,
		p.MinDeposit, p.MinTopup, fmt.Sprintf("%v", p.CcvConsumer),
	}
	return ExecTxJSON(ctx, chain, signer, fees, args...)
}

func ChainletUpdateStack(ctx context.Context, chain ibc.Chain, signer ibc.Wallet, fees string, p UpdateStackParams) (txhash string, code uint32, raw string, err error) {
	args := []string{
		"chainlet", "update-chainlet-stack",
		p.Name, p.Image, p.Version, p.Hash, fmt.Sprintf("%v", p.CcvConsumer),
	}
	return ExecTxJSON(ctx, chain, signer, fees, args...)
}

func ChainletLaunch(ctx context.Context, chain ibc.Chain, signer ibc.Wallet, fees string, p LaunchChainletParams) (txhash string, code uint32, raw string, err error) {
	args := []string{
		"chainlet", "launch-chainlet",
		p.OwnerAddr, p.StackName, p.StackVersion, p.ChainletID, p.ChainletDenom, p.CustomJSON,
	}
	if p.Service {
		args = append(args, "--service-chainlet")
	}
	if p.EVMChainID != "" {
		args = append(args, "--evm-chain-id", p.EVMChainID)
	}
	if p.NetworkVersion != "" {
		args = append(args, "--network-version", p.NetworkVersion)
	}
	if p.Gas != "" {
		args = append(args, "--gas", p.Gas)
	}
	return ExecTxJSON(ctx, chain, signer, fees, args...)
}
