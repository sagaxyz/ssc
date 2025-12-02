package types

import (
	chainlettypes "github.com/sagaxyz/saga-sdk/x/chainlet/types"
	host "github.com/cosmos/ibc-go/v10/modules/core/24-host"
)

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	df := DefaultParams()
	return &GenesisState{
		Params: df,
		PortId: chainlettypes.PortID,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}

	return gs.Params.Validate()
}
