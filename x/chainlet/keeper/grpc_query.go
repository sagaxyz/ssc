package keeper

import (
	"github.com/sagaxyz/ssc/x/chainlet/types"
)

var _ types.QueryServer = &Keeper{}
