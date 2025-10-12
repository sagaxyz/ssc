package types

import (
	"errors"
	"fmt"
	"strconv"
	"time"

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
		ChainletStackProtections:         false,
		NEpochDeposit:                    "30",
		AutomaticChainletUpgrades:        true,
		AutomaticChainletUpgradeInterval: 100,
		LaunchDelay:                      3 * time.Minute,
		MaxChainlets:                     500,
		EnableCCV:                        true,
		UpgradeMinimumHeightDelta:        100,
		UpgradeTimeoutHeight:             500,
		UpgradeTimeoutTime:               12 * time.Hour,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	psp := paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte("ChainletStackProtections"), &p.ChainletStackProtections, validateBool),
		paramtypes.NewParamSetPair([]byte("NEpochDeposit"), &p.NEpochDeposit, validateED),
		paramtypes.NewParamSetPair([]byte("MaxChainlets"), &p.MaxChainlets, validateUint64),
		paramtypes.NewParamSetPair([]byte("AutomaticChainletUpgrades"), &p.AutomaticChainletUpgrades, validateBool),
		paramtypes.NewParamSetPair([]byte("AutomaticChainletUpgradeInterval"), &p.AutomaticChainletUpgradeInterval, validateInt64),
		paramtypes.NewParamSetPair([]byte("LaunchDelay"), &p.LaunchDelay, validateDuration),
		paramtypes.NewParamSetPair([]byte("EnableCCV"), &p.EnableCCV, validateBool),
		paramtypes.NewParamSetPair([]byte("UpgradeMinimumHeightDelta"), &p.UpgradeMinimumHeightDelta, validateUint64),
		paramtypes.NewParamSetPair([]byte("UpgradeTimeoutHeight"), &p.UpgradeTimeoutHeight, validateUint64),
		paramtypes.NewParamSetPair([]byte("UpgradeTimeoutTime"), &p.UpgradeTimeoutTime, validateUint64),
	}

	return psp
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateBool(p.ChainletStackProtections); err != nil {
		return fmt.Errorf("param ChainletStackProtections validation failed: %v", err)
	}
	if err := validateED(p.NEpochDeposit); err != nil {
		return fmt.Errorf("param NEpochDeposit validation failed: %v", err)
	}
	if err := validateBool(p.AutomaticChainletUpgrades); err != nil {
		return fmt.Errorf("param AutomaticChainletUpgrades validation failed: %v", err)
	}
	if err := validateInt64(p.AutomaticChainletUpgradeInterval); err != nil {
		return fmt.Errorf("param AutomaticChainletUpgradeInterval validation failed: %v", err)
	}
	if err := validateBool(p.EnableCCV); err != nil {
		return fmt.Errorf("param EnableCCV validation failed: %v", err)
	}
	if err := validateUint64(p.UpgradeMinimumHeightDelta); err != nil {
		return fmt.Errorf("param UpgradeMinimumHeightDelta validation failed: %v", err)
	}
	if err := validateUint64(p.UpgradeTimeoutHeight); err != nil {
		return fmt.Errorf("param UpgradeTimeoutHeight validation failed: %v", err)
	}
	if err := validateTime(p.UpgradeTimeoutTime); err != nil {
		return fmt.Errorf("param UpgradeTimeoutTime validation failed: %v", err)
	}
	return nil
}

func validateBool(v interface{}) error {
	_, ok := v.(bool)
	if !ok {
		return fmt.Errorf("param not bool")
	}
	return nil
}

func validateInt64(v interface{}) error {
	_, ok := v.(int64)
	if !ok {
		return fmt.Errorf("param not int64")
	}
	return nil
}

func validateED(v interface{}) error {
	vv, ok := v.(string)
	if !ok {
		return fmt.Errorf("could not unmarshal EpochDeposit param for validation")
	}
	_, err := strconv.Atoi(vv)
	if err != nil {
		return err
	}
	return nil
}

func validateDuration(v interface{}) error {
	vv, ok := v.(time.Duration)
	if !ok {
		return errors.New("param not duration")
	}
	if vv < 0 {
		return errors.New("duration negative")
	}
	return nil
}

func validateTime(v interface{}) error {
	vv, ok := v.(time.Time)
	if !ok {
		return errors.New("param not time")
	}
	if vv.IsZero() {
		return errors.New("zero time")
	}
	return nil
}

func validateUint64(v interface{}) error {
	_, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("param not uint64")
	}
	return nil
}
