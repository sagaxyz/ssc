package keeper

import (
	"github.com/sagaxyz/ssc/x/escrow/types"
)

var _ types.QueryServer = Keeper{}
