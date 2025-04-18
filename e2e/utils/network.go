package utils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/client"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
)

// InterchainNetwork is the central struct to manage the interchaintest setup
// consisting of the running chains, any relayers as well as
// a pointer to the actual interchain object.
type InterchainNetwork struct {
	chains      []ibc.Chain
	config      networkConfig
	client      *client.Client
	eRep        *testreporter.RelayerExecReporter
	interchain  *interchaintest.Interchain
	links       RelayerLinks
	networkName string
}

func (icn *InterchainNetwork) GetInterchain() *interchaintest.Interchain { return icn.interchain }
func (icn *InterchainNetwork) GetChains() []ibc.Chain                    { return icn.chains }
func (icn *InterchainNetwork) GetChain(idx uint8) (ibc.Chain, error) {
	if int(idx) >= len(icn.chains) {
		return nil, fmt.Errorf("chain index out of bounds: %d", idx)
	}

	return icn.chains[idx], nil
}

// GetChannelInfo returns the channel information for the first stored channel between the chains with the given indices.
// This information includes e.g. channel and port ids of the connection and the counterparty information as well.
func (icn *InterchainNetwork) GetChannelInfo(ctx context.Context, relayerPath RelayerPath) (ibc.ChannelOutput, error) {
	r := icn.GetLink(relayerPath[0], relayerPath[1]).Relayer

	chain, err := icn.GetChain(relayerPath[0])
	if err != nil {
		return ibc.ChannelOutput{}, err
	}

	// do not use ibc.GetTransferChannel to get the channel, it may return incorrect channel id
	channels, err := r.GetChannels(ctx, icn.eRep, chain.Config().ChainID)
	if err != nil {
		return ibc.ChannelOutput{}, err
	}

	return channels[0], nil
}

// CreateAndStartFullyConnectedNetwork will start the interchaintest setup and create a fully functional
// and connected network.
//
// NOTE 1: This is the main entrypoint to be used in E2E tests.
//
// NOTE 2: There are sensible defaults set to enable running the setup with no options.
// In general, it will be useful to specify the desired setup state using the ConfigOption
// functions that can be passed as an argument.
func CreateAndStartFullyConnectedNetwork(t *testing.T, ctx context.Context, opts ...ConfigOption) (*InterchainNetwork, error) {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	icn, err := newNetwork(t, *config)
	if err != nil {
		return nil, fmt.Errorf("failed to create network: %w", err)
	}

	chains, err := icn.createAndAddChains(t)
	if err != nil {
		return nil, fmt.Errorf("failed to create chains: %w", err)
	}

	if err = icn.createAndAddRelayers(t, chains); err != nil {
		return nil, fmt.Errorf("failed to create relayer: %w", err)
	}

	if err = icn.Build(t, ctx); err != nil {
		return nil, fmt.Errorf("failed to build interchain: %w", err)
	}

	if err = icn.startRelayers(ctx); err != nil {
		return nil, fmt.Errorf("failed to start relayer: %w", err)
	}

	icn.registerRelayerCleanup(t, ctx)

	return icn, nil
}

func newNetwork(t *testing.T, config networkConfig) (*InterchainNetwork, error) {
	dockerClient, network := interchaintest.DockerSetup(t)

	ic := interchaintest.NewInterchain()
	f, err := interchaintest.CreateLogFile(fmt.Sprintf("%d.json", time.Now().Unix()))
	if err != nil {
		return nil, err
	}

	rep := testreporter.NewReporter(f)
	eRep := rep.RelayerExecReporter(t)

	return &InterchainNetwork{
		client:      dockerClient,
		config:      config,
		eRep:        eRep,
		interchain:  ic,
		networkName: network,
	}, nil
}

func (icn *InterchainNetwork) Build(t *testing.T, ctx context.Context) error {
	return icn.interchain.Build(ctx, icn.eRep, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            icn.client,
		NetworkID:         icn.networkName,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		SkipPathCreation:  false,
	})
}
