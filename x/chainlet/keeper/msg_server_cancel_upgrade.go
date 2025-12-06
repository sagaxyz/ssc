package keeper

import (
	"context"
	"fmt"

	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/exp/slices"

	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (k msgServer) CancelChainletUpgrade(goCtx context.Context, msg *types.MsgCancelChainletUpgrade) (*types.MsgCancelChainletUpgradeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := msg.ValidateBasic()
	if err != nil {
		return &types.MsgCancelChainletUpgradeResponse{}, err
	}

	chainlet, err := k.Chainlet(ctx, msg.ChainId)
	if err != nil {
		return &types.MsgCancelChainletUpgradeResponse{}, err
	}
	if !slices.Contains(chainlet.Maintainers, msg.Creator) {
		return nil, fmt.Errorf("address %s is not a chainlet maintainer", msg.Creator)
	}

	currentStack, err := k.getChainletStackVersion(ctx, chainlet.ChainletStackName, chainlet.ChainletStackVersion)
	if err != nil {
		return nil, err
	}
	if !currentStack.CcvConsumer {
		return nil, fmt.Errorf("not supported for chainlet %s (not a consumer)", chainlet.ChainId)
	}

	// Check if upgrade is in progress before attempting to cancel
	if chainlet.Upgrade == nil {
		return nil, cosmossdkerrors.Wrapf(types.ErrNoUpgradeInProgress, "chainlet %s has no upgrade in progress", msg.ChainId)
	}

	err = k.sendCancelUpgradePlan(ctx, &chainlet, msg.ChannelId)
	if err != nil {
		return nil, fmt.Errorf("error sending cancel upgrade: %s", err)
	}

	return &types.MsgCancelChainletUpgradeResponse{}, nil
}
