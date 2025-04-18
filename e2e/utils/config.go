package utils

import "fmt"

// ConfigOption is used to pass adjustments to the default network configuration.
type ConfigOption func(config *networkConfig)

func WithNChains(nChains uint8) ConfigOption {
	return func(config *networkConfig) {
		config.nChains = nChains
	}
}

func WithNFullNodes(nFullNodes uint8) ConfigOption {
	return func(config *networkConfig) {
		config.nFullNodes = nFullNodes
	}
}

func WithNValsPerChain(nValsPerChain uint8) ConfigOption {
	return func(config *networkConfig) {
		config.nValsPerChain = nValsPerChain
	}
}

func WithRelayerPaths(relayerPaths ...RelayerPath) ConfigOption {
	return func(config *networkConfig) {
		config.relayerPaths = relayerPaths
	}
}

// RelayerPath is a type to specify two chain indices in the existing interchain testing network.
type RelayerPath [2]uint8

type networkConfig struct {
	// nChains specifies the number of chains being created.
	nChains uint8
	// nFullNodes specifies the number of full nodes, that are NOT validators, running per chain.
	nFullNodes uint8
	// relayerPaths specifies the relayers to be created in form of indices of two chains that should be connected.
	relayerPaths []RelayerPath
	// nValsPerChain specifies the number of validator nodes, that will be created per chain.
	nValsPerChain uint8
}

// defaultConfig returns a minimal functional network testing configuration.
func defaultConfig() *networkConfig {
	return &networkConfig{
		nChains:       1,
		nFullNodes:    0,
		relayerPaths:  nil,
		nValsPerChain: 1,
	}
}

func (nc networkConfig) validate() error {
	if nc.nChains < 1 {
		return fmt.Errorf("invalid number of chains: %d", nc.nChains)
	}

	if nc.nValsPerChain < 1 {
		return fmt.Errorf("invalid number of validators per chain: %d", nc.nValsPerChain)
	}

	maxRelayers := factorial(nc.nChains - 1)
	if len(nc.relayerPaths) > int(maxRelayers) {
		return fmt.Errorf("incorrect number of relayer paths; expected max. of %d, got %d", maxRelayers, len(nc.relayerPaths))
	}

	for _, relayerPath := range nc.relayerPaths {
		for _, idx := range relayerPath {
			if idx >= nc.nChains {
				return fmt.Errorf("relayer path contains invalid index: %d; max %d allowed", idx, nc.nChains-1)
			}
		}
	}

	return nil
}

// factorial is used to determine the max number of relayer connections between the chains.
func factorial(n uint8) uint8 {
	var result uint8
	if n > 0 {
		result = n * factorial(n-1)
		return result
	}
	return 1
}
