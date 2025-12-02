package types

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	df := DefaultParams()
	return &GenesisState{
		Params:         df,
		Chainlets:      []Chainlet{},
		ChainletStacks: []ChainletStack{},
		ChainletCount:  0,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	// Validate chainlets have unique chain IDs
	chainletIDs := make(map[string]bool)
	for _, chainlet := range gs.Chainlets {
		if chainletIDs[chainlet.ChainId] {
			return ErrChainletExists
		}
		chainletIDs[chainlet.ChainId] = true
	}

	// Validate chainlet stacks have unique display names
	stackNames := make(map[string]bool)
	for _, stack := range gs.ChainletStacks {
		if stackNames[stack.DisplayName] {
			return ErrInvalidChainletStack
		}
		stackNames[stack.DisplayName] = true
	}

	return gs.Params.Validate()
}
