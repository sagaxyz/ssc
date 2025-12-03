package types

import (
	fmt "fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		MaxData: 1024,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	psp := paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte("MaxData"), &p.MaxData, validateUint32),
	}

	return psp
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

func validateUint32(v interface{}) error {
	_, ok := v.(uint32)
	if !ok {
		return fmt.Errorf("param not uint32")
	}
	return nil
}
