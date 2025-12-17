package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/cosmos/cosmos-sdk/x/authx"
	"github.com/cosmos/cosmos-sdk/x/authx/keeper"
)

// NewDecodeStore returns a decoder function closure that umarshals the KVPair's
// Value to the corresponding authx type.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], keeper.GrantKey):
			var grantA, grantB authx.Grant
			cdc.MustUnmarshal(kvA.Value, &grantA)
			cdc.MustUnmarshal(kvB.Value, &grantB)
			return fmt.Sprintf("%v\n%v", grantA, grantB)
		case bytes.Equal(kvA.Key[:1], keeper.GrantQueuePrefix):
			var grantA, grantB authx.GrantQueueItem
			cdc.MustUnmarshal(kvA.Value, &grantA)
			cdc.MustUnmarshal(kvB.Value, &grantB)
			return fmt.Sprintf("%v\n%v", grantA, grantB)
		default:
			panic(fmt.Sprintf("invalid authx key %X", kvA.Key))
		}
	}
}
