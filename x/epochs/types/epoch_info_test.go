package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidateEpochInfo(t *testing.T) {
	testCases := []struct {
		name       string
		ei         EpochInfo
		expectPass bool
	}{
		{
			"invalid - blank identifier",
			EpochInfo{
				"  ",
				time.Now(),
				time.Hour * 24,
				1,
				time.Now(),
				true,
				1,
			},
			false,
		},
		{
			"invalid - epoch duration zero",
			EpochInfo{
				WeekEpochID,
				time.Now(),
				time.Hour * 0,
				1,
				time.Now(),
				true,
				1,
			},
			false,
		},
		{
			"invalid - negative current epoch",
			EpochInfo{
				WeekEpochID,
				time.Now(),
				time.Hour * 24,
				-1,
				time.Now(),
				true,
				1,
			},
			false,
		},
		{
			"invalid - negative epoch start height",
			EpochInfo{
				WeekEpochID,
				time.Now(),
				time.Hour * 24,
				1,
				time.Now(),
				true,
				-1,
			},
			false,
		},
		{
			"pass",
			EpochInfo{
				WeekEpochID,
				time.Now(),
				time.Hour * 24,
				1,
				time.Now(),
				true,
				1,
			},
			true,
		},
	}

	for _, tc := range testCases {
		err := tc.ei.Validate()
		if !tc.expectPass {
			require.NotNil(t, err)
			return
		}
		require.NoError(t, err)
	}
}
