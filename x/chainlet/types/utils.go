package types

import "regexp"

var (
	chainIdRegexp = regexp.MustCompile(`^[a-z]+_[1-9]\d*-[1-9]\d*$`)
	denomRegexp   = regexp.MustCompile(`^[a-z]{3,6}$`)
)

func validateChainId(chainId string) bool {
	if len(chainId) > 40 {
		return false
	}
	return chainIdRegexp.MatchString(chainId)
}

func validateDenom(denom string) bool {
	return denomRegexp.MatchString(denom)
}
