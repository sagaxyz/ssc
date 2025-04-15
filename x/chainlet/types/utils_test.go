package types

import (
	"testing"
)

func TestValidateChainId(t *testing.T) {
	tests := []struct {
		ChainId string
		Valid   bool
	}{
		{"12345", false},
		{"12345_1", false},
		{"12345_a-1", false},
		{"12345_1-a", false},
		{"12345_1--1", false},
		{"12345_-1-1", false},
		{"12345_1-0", false},
		{"12345_0-1", false},
		{"!2345_1-1", false},
		{"12345_1-1a", false},

		{"12345_1-1", false},
		{"12345_11-11", false},
		{"12345_1-9000", false},
		{"12345_123-9000", false},
		{"abcde_123-9000", true},
		{"321AbCdE123_123-9000", false},
		{"_saf-9000", false},
		{"BajkhzXGHASDLK_123-9000", false},
		{"abc_123-9000", true},
	}

	for i, test := range tests {
		valid := validateChainId(test.ChainId)
		if valid != test.Valid {
			t.Fatalf("test %d failed: got %v, expected %v", i+1, valid, test.Valid)
		}
	}
}

func validateTestResult(res bool, err error) bool {
	return err == nil && res
}
