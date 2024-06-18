package keeper

import (
	"github.com/sagaxyz/ssc/x/billing/types"
)

var _ types.QueryServer = Keeper{}
