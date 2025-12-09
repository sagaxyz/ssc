package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ccvtypes "github.com/cosmos/interchain-security/v7/x/ccv/types"
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

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return &types.MsgUpgradeChainletResponse{}, err
	}
	admin := k.aclKeeper.IsAdmin(ctx, creator)

	//NOTE: Non-CCV chainlets have to be manually upgraded by Saga
	if !slices.Contains(ogChainlet.Maintainers, msg.Creator) && (ogChainlet.IsCCVConsumer || !admin) {
		return nil, fmt.Errorf("address %s is not a chainlet maintainer", msg.Creator)
	}
	breakingUpgrade, err := versions.CheckUpgrade(ogChainlet.ChainletStackVersion, msg.StackVersion)
	if err != nil {
		return nil, err
	}
	if breakingUpgrade {
		if ogChainlet.IsCCVConsumer {
			p := k.GetParams(ctx)
			upgradeDelta := p.UpgradeMinimumHeightDelta + msg.HeightDelta
			height, err := k.sendUpgradePlan(ctx, &ogChainlet, msg.StackVersion, upgradeDelta, msg.ChannelId)
			if err != nil {
				return nil, fmt.Errorf("error sending upgrade: %s", err)
			}

			return &types.MsgUpgradeChainletResponse{
				Height: height,
			}, nil
		} else {
			//NOTE: Non-CCV chainlets have to be manually upgraded by Saga
			if !admin {
				return nil, fmt.Errorf("address %s is not allowed to upgrade this chainlet", msg.Creator)
			}

			// Add as a consumer if upgrade enables CCV
			newStack, err := k.getChainletStackVersion(ctx, ogChainlet.ChainletStackName, msg.StackVersion)
			if err != nil {
				return nil, err
			}
			if newStack.CcvConsumer {
				p := k.GetParams(ctx)

				var unbondingPeriod time.Duration
				if msg.UnbondingPeriod != nil {
					unbondingPeriod = *msg.UnbondingPeriod
				} else {
					unbondingPeriod = ccvtypes.DefaultConsumerUnbondingPeriod
				}
				err := k.EnableConsumer(ctx, ogChainlet.ChainId, ctx.BlockTime().Add(p.LaunchDelay), unbondingPeriod)
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
