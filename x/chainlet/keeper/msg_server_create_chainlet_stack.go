package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dactypes "github.com/sagaxyz/saga-sdk/x/acl/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (k msgServer) CreateChainletStack(goCtx context.Context, msg *types.MsgCreateChainletStack) (*types.MsgCreateChainletStackResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return &types.MsgCreateChainletStackResponse{}, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	p := k.GetParams(ctx)
	if p.ChainletStackProtections {
		dacAddr := dactypes.NewAddress(dactypes.AddressFormat_ADDRESS_BECH32, msg.Creator)
		if !k.dacKeeper.Allowed(ctx, dacAddr) {
			return nil, fmt.Errorf("address %s not allowed to create chainlet stacks", msg.Creator)
		}
	}

	metaData := types.ChainletStackParams{
		Image:    msg.Image,
		Version:  msg.Version,
		Checksum: msg.Checksum,
		Enabled: true,
	}
	metaDataUpsert := []types.ChainletStackParams{metaData}
	chainletStack := types.ChainletStack{
		Creator:     msg.Creator,
		DisplayName: msg.DisplayName,
		Description: msg.Description,
		Versions:    metaDataUpsert,
		Fees:        msg.Fees,
	}
	err = k.NewChainletStack(ctx, chainletStack)
	if err != nil {
		return nil, fmt.Errorf("error while adding chainlet stack: %s", err)
	}

	return &types.MsgCreateChainletStackResponse{}, ctx.EventManager().EmitTypedEvent(&types.EventNewChainletStack{
		Creator: msg.Creator,
		Name:    msg.DisplayName,
		Version: msg.Version,
	})
}
