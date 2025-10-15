package keeper

import (
	"context"
	"errors"

	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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

	accAddr, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		err = cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
		return
	}
	valAddr := sdk.ValAddress(accAddr)

	data := types.Data{
		Updated:   ctx.BlockTime(),
		Addresses: msg.Peers,
	}
	k.StoreData(ctx, msg.ChainId, valAddr.String(), data)

	err = ctx.EventManager().EmitTypedEvent(&types.EventUpdatedChainlet{
		ChainId: msg.ChainId,
	})
	if err != nil {
		return
	}

	resp = &types.MsgSetPeersResponse{}
	return
}
