package keeper

import (
	"github.com/sagaxyz/ssc/x/peers/types"
)

var _ types.QueryServer = Keeper{}
