package types

import "fmt"

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	df := DefaultParams()
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:        df,
		PeerData:      []GenesisPeerData{},
		ChainCounters: []GenesisChainCounter{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	// Validate peer data entries have unique {chainId, validatorAddress} pairs
	peerDataKeys := make(map[string]bool)
	for _, pd := range gs.PeerData {
		key := pd.ChainId + "/" + pd.ValidatorAddress
		if peerDataKeys[key] {
			return fmt.Errorf("duplicate peer data for chain %s and validator %s", pd.ChainId, pd.ValidatorAddress)
		}
		peerDataKeys[key] = true
	}

	// Validate chain counters have unique chain IDs
	chainCounterKeys := make(map[string]bool)
	for _, cc := range gs.ChainCounters {
		if chainCounterKeys[cc.ChainId] {
			return fmt.Errorf("duplicate chain counter for chain %s", cc.ChainId)
		}
		chainCounterKeys[cc.ChainId] = true
	}

	return gs.Params.Validate()
}
