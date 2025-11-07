package types

import (
	fmt "fmt"

	cosmossdkerrors "cosmossdk.io/errors"
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
		SupportedDenoms: []string{"utsaga", "utagas"},
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	psp := paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte("SupportedDenoms"), &p.SupportedDenoms, validateDenoms),
	}

	return psp
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateDenoms(p.SupportedDenoms); err != nil {
		return cosmossdkerrors.Wrapf(ErrInvalidParams, err.Error())
	}
	return nil
}

func validateDenoms(v interface{}) error {
	denoms, ok := v.([]string)
	if !ok {
		return fmt.Errorf("not a string slice")
	}
	for _, denom := range denoms {
		if err := validateDenom(denom); err != nil {
			return err
		}
	}
	return nil
}

func validateDenom(v interface{}) error {
	_, ok := v.(string)
	if !ok {
		return fmt.Errorf("not a string")
	}
	return nil
}
