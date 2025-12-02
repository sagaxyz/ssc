package types

import "strconv"

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:                 DefaultParams(),
		BillingHistory:         []SaveBillingHistory{},
		ValidatorPayoutHistory: []ValidatorPayoutHistory{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	// Validate billing history records have unique {chainletId, epochIdentifier, epochNumber} tuples
	billingKeys := make(map[string]bool)
	for _, bh := range gs.BillingHistory {
		key := bh.ChainletId + "/" + bh.EpochIdentifier + "/" + strconv.FormatInt(bh.EpochNumber, 10)
		if billingKeys[key] {
			return ErrDuplicateRecord
		}
		billingKeys[key] = true
	}

	// Validate validator payout history records have unique {validatorAddress, epochIdentifier, epochNumber} tuples
	payoutKeys := make(map[string]bool)
	for _, vph := range gs.ValidatorPayoutHistory {
		key := vph.ValidatorAddress + "/" + vph.EpochIdentifier + "/" + strconv.FormatInt(vph.EpochNumber, 10)
		if payoutKeys[key] {
			return ErrDuplicateRecord
		}
		payoutKeys[key] = true
	}

	return gs.Params.Validate()
}
