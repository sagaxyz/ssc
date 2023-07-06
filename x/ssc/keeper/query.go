package keeper

import (
	"github.com/sagaxyz/ssc/x/ssc/types"
)

var _ types.QueryServer = Keeper{}
