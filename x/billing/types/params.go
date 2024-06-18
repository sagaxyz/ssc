package types

import (
	fmt "fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{ValidatorPayoutEpoch: SAGA_EPOCH_IDENTIFIER}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {

	psp := paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte("ValidatorPayoutEpoch"), &p.ValidatorPayoutEpoch, validatePayoutEpochParam),
	}

	return psp
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validatePayoutEpochParam(v interface{}) error {
	_, ok := v.(string)
	if !ok {
		return fmt.Errorf("could not unmarshal validator-payout-epoch parm for validation")
	}
	return nil
}
