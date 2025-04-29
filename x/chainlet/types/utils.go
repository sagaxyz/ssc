package types

import (
	fmt "fmt"
	"regexp"
)

var (
	chainIdRegexp         = regexp.MustCompile(`^[a-z]+_[1-9]\d*-[1-9]\d*$`)
	denomRegexp           = regexp.MustCompile(`^[a-z]{3,6}$`)
	nonAdminChainIdRegexp = regexp.MustCompile(`^[a-z]+_[1-9]\\d*-1$`)
)

func validateChainId(chainId string) bool {
	if len(chainId) > 40 {
		return false
	}
	return chainIdRegexp.MatchString(chainId)
}

func ValidateNonAdminChainId(chainId string) (bool, error) {
	if len(chainId) > 40 {
		return false, nil
	}
	// return nonAdminChainIdRegexp.MatchString(chainId), nil
	return regexp.MatchString("^[a-z]+_[1-9]\\d*-1$", chainId)
}

func validateDenom(denom string) bool {
	return denomRegexp.MatchString(denom)
}

func GenerateChainId(name string, evm, version int64) string {
	return fmt.Sprintf("%s_%d-%d", name, evm, version)
}
