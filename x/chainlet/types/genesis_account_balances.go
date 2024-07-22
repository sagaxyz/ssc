package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

// This is used by the controller in order for the chainlet params
// values be something that can be passed as an env variable.
func (gab GenesisAccountBalances) MarshalJSON() ([]byte, error) {
	var buf strings.Builder
	for i, acc := range gab.List {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(acc.Address)
		buf.WriteString("=")
		buf.WriteString(acc.Balance)
	}
	return json.Marshal(buf.String())
}

func (gab *GenesisAccountBalances) UnmarshalJSON(data []byte) error {
	input := string(data)
	input = strings.TrimPrefix(input, "\"")
	input = strings.TrimSuffix(input, "\"")
	abs := strings.Split(input, ",")
	gab.List = make([]*AccountBalance, 0, len(abs))
	for _, ab := range abs {
		split := strings.Split(ab, "=")
		if len(split) != 2 {
			return fmt.Errorf("bad entry in genesis balances: %s", ab)
		}

		gab.List = append(gab.List, &AccountBalance{
			Address: split[0],
			Balance: split[1],
		})
	}

	return nil
}
