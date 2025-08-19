package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/exp/slices"

	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/sagaxyz/ssc/x/chainlet/types/versions"
)

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

	if !slices.Contains(ogChainlet.Maintainers, msg.Creator) {
		return nil, fmt.Errorf("address %s is not a chainlet maintainer", msg.Creator)
	}
	majorUpgrade, err := versions.CheckUpgrade(ogChainlet.ChainletStackVersion, msg.StackVersion)
	if err != nil {
		return nil, err
	}
	if majorUpgrade {
		currentStack, err := k.getChainletStackVersion(ctx, ogChainlet.ChainletStackName, ogChainlet.ChainletStackVersion)
		if err != nil {
			return nil, err
		}
		if currentStack.CcvConsumer {
			p := k.GetParams(ctx)
			upgradeDelta := p.MinimumUpgradeHeightDelta + msg.HeightDelta
			height, err := k.sendUpgradePlan(ctx, &ogChainlet, msg.StackVersion, upgradeDelta, msg.ChannelID)
			if err != nil {
				return nil, fmt.Errorf("error sending upgrade: %s", err)
			}

			return &types.MsgUpgradeChainletResponse{
				Height: height,
			}, nil
		} else {
			// Add as a consumer if upgrade enables CCV
			newStack, err := k.getChainletStackVersion(ctx, ogChainlet.ChainletStackName, msg.StackVersion)
			if err != nil {
				return nil, err
			}
			if newStack.CcvConsumer {
				p := k.GetParams(ctx)
				err := k.EnableConsumer(ctx, ogChainlet.ChainId, ctx.BlockTime().Add(p.LaunchDelay))
				if err != nil {
					return nil, err
				}
			}
		}
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
