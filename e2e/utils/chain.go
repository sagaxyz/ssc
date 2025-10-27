package utils

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	interchaintest "github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

const (
	// Genesis params (aligned with prepare-env.sh, with staking/bond denom = utsaga)
	VotingPeriod     = "10s"
	MaxDepositPeriod = "10s"
	ExpeditedVoting  = "5s"

	MinDeposit          = "250000000"
	Quorum              = "0.4"
	Threshold           = "0.5"
	VetoThreshold       = "0.4"
	DelegationAmount    = 8000000
	Denom               = "utsaga" // bank/fee denom
	StakeDenom          = "utsaga" // staking (bond) denom
	MaxValidators       = 21
	SignedBlocksWindow  = "17280"
	Inflation           = "0.07"
	InflationMin        = "0.02"
	InflationMax        = "0.08"
	InflationRateChange = "1.0"
)

func GetDockerImage() (repo, version string) {
	repo = "sagaxyz/ssc"
	version = "e2e"
	return repo, version
}

func (icn *InterchainNetwork) createAndAddChains(t *testing.T) ([]ibc.Chain, error) {
	chainSpecs := make([]*interchaintest.ChainSpec, icn.config.nChains)
	for i := range icn.config.nChains {
		chainSpecs[i] = createChainSpec(
			icn.config,
			fmt.Sprintf("chain%d", i+1),
			fmt.Sprintf("ssc_%d-1", i+1),
			Denom,
			StakeDenom,
		)
	}

	chains, err := createIbcChains(t, chainSpecs)
	if err != nil {
		return nil, fmt.Errorf("failed to create chains: %w", err)
	}

	icn.chains = chains
	for _, chain := range chains {
		icn.interchain.AddChain(chain)
	}

	return chains, nil
}

func createChainSpec(config networkConfig, name, chainID, denom, stakeDenom string) *interchaintest.ChainSpec {
	repo, version := GetDockerImage()
	nVals := int(config.nValsPerChain)
	nFullNodes := int(config.nFullNodes)

	defaultGenesisKV := []cosmos.GenesisKV{
		// --- Governance ---
		{Key: "app_state.gov.params.voting_period", Value: VotingPeriod},
		{Key: "app_state.gov.params.max_deposit_period", Value: MaxDepositPeriod},
		{Key: "app_state.gov.params.expedited_voting_period", Value: ExpeditedVoting},
		{Key: "app_state.gov.params.quorum", Value: Quorum},
		{Key: "app_state.gov.params.threshold", Value: Threshold},
		{Key: "app_state.gov.params.veto_threshold", Value: VetoThreshold},
		{Key: "app_state.gov.params.min_deposit", Value: []map[string]string{{"denom": denom, "amount": MinDeposit}}},

		// --- Chainlet params ---
		{Key: "app_state.chainlet.params.chainletStackProtections", Value: false},
		{Key: "app_state.chainlet.params.nEpochDeposit", Value: "30"},

		// --- ACL ---
		{Key: "app_state.acl.params.enable", Value: true},

		// --- Billing epochs ---
		{Key: "app_state.billing.params.billing_epoch", Value: "minute"},
		{Key: "app_state.billing.params.validator_payout_epoch", Value: "minute"},

		// --- Distribution / Staking / Slashing / Mint ---
		{Key: "app_state.distribution.params.community_tax", Value: "0.042"},
		{Key: "app_state.distribution.params.base_proposer_reward", Value: "0.01"},
		{Key: "app_state.distribution.params.bonus_proposer_reward", Value: "0.04"},

		// IMPORTANT: make bond denom = utsaga
		{Key: "app_state.staking.params.bond_denom", Value: stakeDenom},
		{Key: "app_state.staking.params.max_validators", Value: MaxValidators},

		{Key: "app_state.slashing.params.signed_blocks_window", Value: SignedBlocksWindow},

		// Mint module: align mint denom with utsaga
		{Key: "app_state.mint.params.mint_denom", Value: denom},
		{Key: "app_state.mint.minter.inflation", Value: Inflation},
		{Key: "app_state.mint.params.inflation_min", Value: InflationMin},
		{Key: "app_state.mint.params.inflation_max", Value: InflationMax},
		{Key: "app_state.mint.params.inflation_rate_change", Value: InflationRateChange},
	}

	return &interchaintest.ChainSpec{
		ChainConfig: ibc.ChainConfig{
			Type:    "cosmos",
			Name:    name,
			ChainID: chainID,
			Images: []ibc.DockerImage{
				{Repository: repo, Version: version, UIDGID: "1025:1025"},
			},
			Bech32Prefix:   "saga",
			Bin:            "sscd",
			Denom:          denom,       // bank/fee denom: utsaga
			GasPrices:      "0" + denom, // "0utsaga"
			GasAdjustment:  1.3,
			TrustingPeriod: "508h",
			NoHostMount:    false,
			ModifyGenesis:  cosmos.ModifyGenesis(defaultGenesisKV),

			// ensure validators delegate in the bond denom (utsaga)
			ModifyGenesisAmounts: func(i int) (sdk.Coin, sdk.Coin) {
				delegation := sdk.NewCoin(stakeDenom, math.NewInt(DelegationAmount)) // utsaga
				return delegation, delegation
			},
		},
		NumValidators: &nVals,
		NumFullNodes:  &nFullNodes,
	}
}

func createIbcChains(t *testing.T, chainSpecs []*interchaintest.ChainSpec) ([]ibc.Chain, error) {
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t, zaptest.Level(zap.ErrorLevel)), chainSpecs)
	return cf.Chains(t.Name())
}

// getNode casts the given ibc.Chain to a Cosmos chain and returns
// the first validator node.
func getNode(chain ibc.Chain) (*cosmos.ChainNode, error) {
	cc, ok := chain.(*cosmos.CosmosChain)
	if !ok {
		return nil, fmt.Errorf("expected CosmosChain; got %T", chain)
	}

	return cc.GetNode(), nil
}
