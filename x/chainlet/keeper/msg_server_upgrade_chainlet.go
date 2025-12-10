package keeper

import (
	"context"
	"errors"
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

	isAdmin := k.aclKeeper.IsAdmin(ctx, creator)
	isMaintainer := slices.Contains(ogChainlet.Maintainers, msg.Creator)
	canUpgrade := isMaintainer || (!ogChainlet.IsCCVConsumer && isAdmin) // Non-CCV chainlets have to be manually upgraded by Saga
	if !canUpgrade {
		if ogChainlet.IsCCVConsumer {
			return nil, fmt.Errorf("address %s is not a chainlet maintainer", msg.Creator)
		}
		return nil, fmt.Errorf("address %s is not allowed to upgrade this chainlet (must be maintainer or admin)", msg.Creator)
	}

	newStack, err := k.getChainletStackVersion(ctx, ogChainlet.ChainletStackName, msg.StackVersion)
	if err != nil {
		return nil, err
	}
	breakingUpgrade, err := versions.CheckUpgrade(ogChainlet.ChainletStackVersion, msg.StackVersion)
	if err != nil {
		return nil, err
	}
	if breakingUpgrade {
		if ogChainlet.IsCCVConsumer {
			if !newStack.CcvConsumer {
				return &types.MsgUpgradeChainletResponse{}, errors.New("CCV cannot be disabled")
			}
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
			if !isAdmin {
				return nil, fmt.Errorf("address %s is not allowed to upgrade this chainlet", msg.Creator)
			}

			// Add as a consumer if upgrade enables CCV
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
	} else {
		if newStack.CcvConsumer != ogChainlet.IsCCVConsumer {
			return &types.MsgUpgradeChainletResponse{}, errors.New("changing CCV requires a breaking upgrade")
		}
	}

	err = k.UpgradeChainletStackVersion(ctx, msg.ChainId, msg.StackVersion)
	if err != nil {
		return &types.MsgUpgradeChainletResponse{}, fmt.Errorf("error while updating chainlet: %s", err)
	}

	return &types.MsgUpgradeChainletResponse{}, ctx.EventManager().EmitTypedEvent(&types.EventUpdateChainlet{
		ChainId:      msg.ChainId,
		StackVersion: msg.StackVersion,
	})
}
