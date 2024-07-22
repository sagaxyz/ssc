package types

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// DefaultIndex is the default capability global index.
const DefaultIndex uint64 = 1

func NewGenesisState(epochs []EpochInfo) *GenesisState {
	return &GenesisState{Epochs: epochs}
}

// DefaultGenesis returns the default Capability genesis state.
func DefaultGenesis() *GenesisState {
	epochs := []EpochInfo{
		NewGenesisEpochInfo(DayEpochID, time.Hour*24), // alphabetical order
		NewGenesisEpochInfo(HourEpochID, time.Hour),
		NewGenesisEpochInfo(MinuteEpochID, time.Minute),
		NewGenesisEpochInfo(WeekEpochID, time.Hour*24*7),
	}
	return NewGenesisState(epochs)
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	epochIdentifiers := map[string]bool{}
	for _, epoch := range gs.Epochs {
		if err := epoch.Validate(); err != nil {
			return err
		}
		if epochIdentifiers[epoch.Identifier] {
			return errors.New("epoch identifier should be unique")
		}
		epochIdentifiers[epoch.Identifier] = true
	}
	return nil
}

// Validate also validates epoch info.
func (ei EpochInfo) Validate() error {
	if strings.TrimSpace(ei.Identifier) == "" {
		return errors.New("epoch identifier cannot be blank")
	}
	if ei.Duration == 0 {
		return errors.New("epoch duration cannot be 0")
	}
	if ei.CurrentEpoch < 0 {
		return fmt.Errorf("current epoch cannot be negative: %d", ei.CurrentEpochStartHeight)
	}
	if ei.CurrentEpochStartHeight < 0 {
		return fmt.Errorf("current epoch start height cannot be negative: %d", ei.CurrentEpochStartHeight)
	}
	return nil
}

func NewGenesisEpochInfo(identifier string, duration time.Duration) EpochInfo {
	return EpochInfo{
		Identifier:              identifier,
		StartTime:               time.Time{},
		Duration:                duration,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		CurrentEpochStartTime:   time.Time{},
		EpochCountingStarted:    false,
	}
}
