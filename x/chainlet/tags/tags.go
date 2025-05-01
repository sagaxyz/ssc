package tags

import "golang.org/x/exp/slices"

var Tags = []string{"hub", "dex", "oracle"}

func VerifyAndTruncate(input []string) ([]string, bool) {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range input {
		if !slices.Contains(Tags, item) {
			return nil, false
		}
		if seen[item] {
			continue
		}
		seen[item] = true
		result = append(result, item)
	}

	return result, true
}
