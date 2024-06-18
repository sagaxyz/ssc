package keeper

import (
	"context"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/exp/slices"

	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

const SagaAddress = "saga1h8r6gm4jehflfn2nn7mtw53l37skrke5kyax8l"

func (k msgServer) UpgradeChainlet(goCtx context.Context, msg *types.MsgUpgradeChainlet) (*types.MsgUpgradeChainletResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := msg.ValidateBasic()
	if err != nil {
		return &types.MsgUpgradeChainletResponse{}, err
	}

	ogChainlet, err := k.Chainlet(ctx, msg.ChainId)
	if err != nil {
		return &types.MsgUpgradeChainletResponse{}, err
	}

	if !slices.Contains(ogChainlet.Maintainers, msg.Creator) && msg.Creator != SagaAddress {
		return nil, fmt.Errorf("address %s not whitelisted for creating or updating chainlet stacks", msg.Creator)
	}
	majorUpgrade, err := versions.CheckUpgrade(ogChainlet.ChainletStackVersion, msg.StackVersion)
	if err != nil {
		return nil, err
	}
	// Only this Saga-controlled address is allowed to perform (manual) major upgrades until they're automated using IBC
	if majorUpgrade && msg.Creator != SagaAddress {
		return nil, errors.New("major upgrades not implemented")
	}

	err = k.UpgradeChainletStackVersion(ctx, msg.ChainId, msg.StackVersion)
	if err != nil {
		return nil, fmt.Errorf("error while updating chainlet: %s", err)
	}

	return &types.MsgUpgradeChainletResponse{}, ctx.EventManager().EmitTypedEvent(&types.EventUpdateChainlet{
		ChainId:      msg.ChainId,
		StackVersion: msg.StackVersion,
	})
}
