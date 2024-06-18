package keeper

import (
	"context"
	"fmt"

	cosmossdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (k msgServer) LaunchChainlet(goCtx context.Context, msg *types.MsgLaunchChainlet) (*types.MsgLaunchChainletResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return &types.MsgLaunchChainletResponse{}, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	ch, err := k.ListChainlets(goCtx, &types.QueryListChainletsRequest{
		Pagination: &query.PageRequest{
			Limit: 1000,
		},
	})
	if err != nil {
		return nil, err
	}
	if len(ch.Chainlets) >= 500 {
		return nil, types.ErrTooManyChainlets
	}

	if k.ChainletExists(ctx, msg.ChainId) {
		return &types.MsgLaunchChainletResponse{}, types.ErrChainletExists
	}

	// pad genesis balances
	for idx, bal := range msg.Params.GenAcctBalances.List {
		amount, err := math.ParseUint(bal.Balance + "000000000000000000")
		if err != nil {
			return nil, err
		}
		if amount.IsZero() {
			msg.Params.GenAcctBalances.List = append(msg.Params.GenAcctBalances.List[:idx], msg.Params.GenAcctBalances.List[idx+1:]...)
			continue
		}
		msg.Params.GenAcctBalances.List[idx].Balance = amount.String()
	}

	p := k.GetParams(ctx)
	chainlet := types.Chainlet{
		SpawnTime:            ctx.BlockTime().Add(p.LaunchDelay),
		Launcher:             msg.Creator,
		Maintainers:          msg.Maintainers,
		ChainletStackName:    msg.ChainletStackName,
		ChainletStackVersion: msg.ChainletStackVersion,
		ChainletName:         msg.ChainletName,
		ChainId:              msg.ChainId,
		Denom:                msg.Denom,
		Params:               msg.Params,
		Status:               types.Status_STATUS_ONLINE,
		AutoUpgradeStack:     !msg.DisableAutomaticStackUpgrades,
	}
	stack, err := k.GetChainletStack(goCtx, &types.QueryGetChainletStackRequest{DisplayName: msg.ChainletStackName})
	if err != nil {
		return &types.MsgLaunchChainletResponse{}, types.ErrInvalidChainletStack
	}
	epochfee, err := sdk.ParseCoinNormalized(stack.ChainletStack.Fees.EpochFee)
	if err != nil {
		return &types.MsgLaunchChainletResponse{}, types.ErrInvalidCoin
	}
	setupfee, err := sdk.ParseCoinNormalized(stack.ChainletStack.Fees.SetupFee)
	if err != nil {
		return &types.MsgLaunchChainletResponse{}, types.ErrInvalidCoin
	}
	owner, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return &types.MsgLaunchChainletResponse{}, err
	}

	multiplier, ok := math.NewIntFromString(k.GetParams(ctx).NEpochDeposit)
	if !ok {
		return &types.MsgLaunchChainletResponse{}, fmt.Errorf("bad multiplier")
	}

	deposit := sdk.Coin{
		Amount: epochfee.Amount.Mul(multiplier),
		Denom:  epochfee.Denom,
	}
	deposit.Add(setupfee)
	err = k.escrowKeeper.NewChainletAccount(ctx, owner, msg.ChainId, deposit)
	if err != nil {
		return &types.MsgLaunchChainletResponse{}, err
	}

	// Bill for the chainlet just after it is launched
	totalFee := epochfee.Add(setupfee)
	err = k.billingKeeper.BillAccount(ctx, totalFee, chainlet, stack.ChainletStack.Fees.EpochLength, "launching chainlet")
	if err != nil {
		return &types.MsgLaunchChainletResponse{}, cosmossdkerrors.Wrapf(types.ErrBillingFailure, fmt.Sprintf("%v", err))
	}

	// Add as a CCV consumer
	err = k.addConsumer(ctx, chainlet.ChainId, chainlet.SpawnTime)
	if err != nil {
		return nil, err
	}

	err = k.NewChainlet(ctx, chainlet)
	if err != nil {
		return nil, err
	}

	return &types.MsgLaunchChainletResponse{}, ctx.EventManager().EmitTypedEvent(&types.EventLaunchChainlet{
		ChainName:    msg.ChainletName,
		Launcher:     msg.Creator,
		ChainId:      msg.ChainId,
		Stack:        msg.ChainletStackName,
		StackVersion: msg.ChainletStackVersion,
	})
}
