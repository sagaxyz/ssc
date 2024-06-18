package versions_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

func TestCheckVersion(t *testing.T) {
	tests := []struct {
		version string
		valid   bool
	}{
		{"0.1.2", true},
		{"0.1.2-alpha.1", true},
		{"0.1.2_alpha", false},
		{"v0.1.2", false},
		{"", false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.version, func(t *testing.T) {
			valid := versions.Check(tt.version)
			require.Equal(t, tt.valid, valid)
		})
	}
}

func TestCheckVersionUpgrade(t *testing.T) {
	tests := []struct {
		versionOld string
		versionNew string
		valid      bool
		major      bool
	}{
		{"0.1.2", "0.1.3", true, false},
		{"0.1.2", "0.2.0", true, true},
		{"1.2.3", "1.3.0", true, false},
		{"1.2.3", "2.0.0", true, true},
		{"2.0.0", "2.0.0", true, false},
		{"1.2.3", "v2.0.0", false, false},
		{"1.2.3", "", false, false},
		{"", "1.2.3", false, false},
		{"v2.0.0", "3.0.0", false, false},
		{"1.0.0", "0.1.2", false, false},
		{"2.0.0", "1.0.0", false, false},
		{"1.0.0", "3.0.0", false, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%s->%s", tt.versionOld, tt.versionNew), func(t *testing.T) {
			major, err := versions.CheckUpgrade(tt.versionOld, tt.versionNew)
			if tt.valid {
				require.NoError(t, err)
				require.Equal(t, tt.major, major)
			} else {
				require.Error(t, err)
			}
		})
	}
}
