package versions_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

func AddBatch(sv *versions.Versions, versions []string) error {
	for _, v := range versions {
		err := sv.Add(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestVersionsInsertRemove(t *testing.T) {
	tests := []struct {
		insert   []string
		remove   []string
		expected []string
	}{
		{
			[]string{"0.1.2", "0.1.3"},
			nil,
			[]string{"0.1.2", "0.1.3"},
		},
		{
			[]string{"0.1.2"},
			[]string{"0.1.2"},
			[]string{},
		},
		{
			[]string{"0.1.2"},
			[]string{"0.1.3"},
			[]string{"0.1.2"},
		},
		{
			[]string{"0.1.2", "0.1.3", "0.1.4", "1.0.0", "2.0.0", "2.0.1"},
			[]string{"0.1.2", "0.1.3", "0.1.4", "1.0.0", "2.0.0"},
			[]string{"2.0.1"},
		},
		{
			[]string{"0.1.2", "0.1.3", "0.1.4", "1.0.0", "2.0.0", "2.0.1"},
			[]string{"0.1.2", "0.1.3", "0.1.4", "1.0.0", "2.0.0", "2.0.1"},
			[]string{},
		},
		{
			[]string{"2.0.0", "0.1.3", "0.1.2", "1.0.0", "2.0.1", "0.1.4"},
			[]string{"1.0.0", "0.1.3", "2.0.1", "2.0.0", "0.1.4", "0.1.2"},
			[]string{},
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d: insert %v, remove %v", i, tt.insert, tt.remove), func(t *testing.T) {
			versions := versions.New()

			err := AddBatch(versions, tt.insert)
			require.NoError(t, err)

			for _, v := range tt.remove {
				err := versions.Remove(v)
				require.NoError(t, err)

				e := versions.Export()
				t.Logf("removed %s, now: %+v", v, e)
			}

			exported := versions.Export()
			require.Equal(t, tt.expected, exported)
		})
	}
}

func TestVersionsLatestCompatible(t *testing.T) {
	tests := []struct {
		versions       []string
		current        string
		expectedLatest string
		err            bool
	}{
		// No other versions
		{[]string{"0.1.2"}, "0.1.2", "0.1.2", false},
		{[]string{"1.2.3"}, "1.2.3", "1.2.3", false},

		// Only older versions
		{[]string{"0.1.2", "0.1.3"}, "0.1.3", "0.1.3", false},
		{[]string{"1.2.2", "1.2.3"}, "1.2.3", "1.2.3", false},

		// Do not jump major version
		{[]string{"0.1.2", "0.1.3", "0.2.0"}, "0.1.2", "0.1.3", false},
		{[]string{"0.1.2", "0.1.3", "1.0.0", "0.2.0"}, "0.1.2", "0.1.3", false},
		{[]string{"0.1.0", "1.0.0", "2.0.0"}, "1.0.0", "1.0.0", false},
		{[]string{"0.1.0", "1.0.0", "1.1.0", "2.0.0"}, "1.0.0", "1.1.0", false},
		{[]string{"0.1.0", "1.0.0", "1.1.0", "1.1.1", "2.0.0"}, "1.0.0", "1.1.1", false},
		{[]string{"0.1.0", "1.0.0", "1.1.0", "1.1.1", "2.0.0"}, "1.0.0-alpha", "1.1.1", false},

		// Current version is newer than any existing one
		{[]string{"0.1.2"}, "0.1.3", "0.1.3", false},
		{[]string{"1.2.3"}, "1.2.4", "1.2.4", false},
		{[]string{"1.2.3"}, "0.3.0", "0.3.0", false},
		{[]string{"1.2.3"}, "2.0.0", "2.0.0", false},
		{[]string{"1.2.2", "1.2.3"}, "1.2.4-alpha", "1.2.4-alpha", false},
		{[]string{"1.2.2", "1.2.3"}, "1.2.4-alpha", "1.2.4-alpha", false},
		{[]string{"1.2.2", "1.2.3"}, "1.3.0-alpha", "1.3.0-alpha", false},
		{[]string{"1.2.2", "1.2.3"}, "2.0.0-alpha", "2.0.0-alpha", false},

		// Invalid input
		{[]string{}, "v0.1.2", "", true},
		{[]string{}, "", "", true},
		{[]string{}, "test", "", true},
		{[]string{}, "1.0", "", true},
	}

	for i, tt := range tests {
		tt := tt
		if tt.err {
			tt.expectedLatest = "error"
		}
		t.Run(fmt.Sprintf("%d: %s -> %s", i, tt.current, tt.expectedLatest), func(t *testing.T) {
			versions := versions.New()

			err := AddBatch(versions, tt.versions)
			require.NoError(t, err)

			latest, err := versions.LatestCompatible(tt.current)
			if tt.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedLatest, latest)
		})
	}
}
