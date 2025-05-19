package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/peers/types"
)

func (k msgServer) SetPeers(goCtx context.Context, msg *types.MsgSetPeers) (resp *types.MsgSetPeersResponse, err error) {
	err = msg.ValidateBasic()
	if err != nil {
		return
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err = k.chainletKeeper.Chainlet(ctx, msg.ChainId)
	if err != nil {
		err = errors.New("no such chain ID")
		return
	}

	data := types.Data{
		Updated:   ctx.BlockTime(),
		Addresses: msg.Peers,
	}
	err = k.StoreData(ctx, msg.ChainId, msg.Validator, data)
	if err != nil {
		panic(err) //TODO
	}

	err = ctx.EventManager().EmitTypedEvent(&types.EventUpdatedChainlet{
		ChainId: msg.ChainId,
	})
	if err != nil {
		return
	}

	return
}
