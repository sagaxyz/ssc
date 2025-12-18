package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("saga", "sagapub")
	config.SetBech32PrefixForValidator("sagavaloper", "sagavaloperpub")
	config.SetBech32PrefixForConsensusNode("sagavalcons", "sagavalconspub")
}

func TestValidatePlatformValidatorsParam(t *testing.T) {
	addr1 := "saga14grgksm5pe5u4cf8pvcchfsl8mfg8mvj3c95l0"
	addr2 := "saga1nkm3et2qcqgya0ad8wt6e20l6206xdtccw28c9"
	addr3 := "saga1zggrvdnjzsfpc7sr7jnw2jl5v9u8vjn263x0s0"

	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid empty list",
			input:   []string{},
			wantErr: false,
		},
		{
			name:    "valid single address",
			input:   []string{addr1},
			wantErr: false,
		},
		{
			name:    "valid multiple addresses",
			input:   []string{addr1, addr2},
			wantErr: false,
		},
		{
			name:    "invalid address format",
			input:   []string{"invalid-address"},
			wantErr: true,
			errMsg:  "invalid platform validator address",
		},
		{
			name:    "duplicate addresses",
			input:   []string{addr1, addr1},
			wantErr: true,
			errMsg:  "duplicate platform validator address",
		},
		{
			name:    "duplicate addresses with valid in between",
			input:   []string{addr1, addr2, addr1},
			wantErr: true,
			errMsg:  "duplicate platform validator address",
		},
		{
			name:    "three unique addresses",
			input:   []string{addr1, addr2, addr3},
			wantErr: false,
		},
		{
			name:    "wrong type",
			input:   "not a slice",
			wantErr: true,
			errMsg:  "could not unmarshal platform-validators",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePlatformValidatorsParam(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
