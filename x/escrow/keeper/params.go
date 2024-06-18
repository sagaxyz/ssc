package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sagaxyz/ssc/x/escrow/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramstore.GetParamSetIfExists(ctx, &p)
	if p.Size() > 0 {
		return p
	}
	return types.NewParams()
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
