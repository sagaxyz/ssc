package types

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:           DefaultParams(),
		ChainletAccounts: []ChainletAccount{},
		Pools:            []DenomPool{},
		Funders:          []GenesisFunder{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	// Validate chainlet accounts have unique chain IDs
	chainletIDs := make(map[string]bool)
	for _, acc := range gs.ChainletAccounts {
		if chainletIDs[acc.ChainId] {
			return ErrChainletAccountNotFound // reuse existing error type
		}
		chainletIDs[acc.ChainId] = true
	}

	// Validate pools have unique {chainId, denom} pairs
	poolKeys := make(map[string]bool)
	for _, pool := range gs.Pools {
		key := pool.ChainId + "/" + pool.Denom
		if poolKeys[key] {
			return ErrChainletAccountNotFound
		}
		poolKeys[key] = true
	}

	// Validate funders have unique {chainId, denom, address} tuples
	funderKeys := make(map[string]bool)
	for _, f := range gs.Funders {
		key := f.ChainId + "/" + f.Denom + "/" + f.Address
		if funderKeys[key] {
			return ErrFunderNotFound
		}
		funderKeys[key] = true
	}

	return gs.Params.Validate()
}
