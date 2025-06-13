package keeper

import (
	"context"
	"fmt"

	cosmossdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (k msgServer) LaunchChainlet(goCtx context.Context, msg *types.MsgLaunchChainlet) (*types.MsgLaunchChainletResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return &types.MsgLaunchChainletResponse{}, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	admin := k.aclKeeper.Admin(ctx, msg.GetSigners()[0])
	if !admin {
		ok, err := types.ValidateNonAdminChainId(msg.ChainId)
		if !ok {
			return &types.MsgLaunchChainletResponse{}, err
		}
	}

	// get total number of chainlets
	chainletCountRes, err := k.GetChainletCount(ctx, &types.QueryGetChainletCountRequest{})
	if err != nil {
		return nil, err
	}

	if chainletCountRes.Count >= k.GetParams(ctx).MaxChainlets {
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

	stack, err := k.getChainletStack(ctx, msg.ChainletStackName)
	if err != nil {
		return &types.MsgLaunchChainletResponse{}, types.ErrInvalidChainletStack
	}
	stackVersion, err := k.getChainletStackVersion(ctx, msg.ChainletStackName, msg.ChainletStackVersion)
	if err != nil {
		return &types.MsgLaunchChainletResponse{}, types.ErrInvalidChainletStack
	}
	if stackVersion.CcvConsumer && !p.EnableCCV {
		return &types.MsgLaunchChainletResponse{}, types.ErrInvalidChainletStack
	}

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
		GenesisValidators:    k.validators(ctx),
		IsServiceChainlet:    msg.IsServiceChainlet,
		IsCCVConsumer:        stackVersion.CcvConsumer,
	}

	// launching a service chainlet means we can skip the billing setup and just create the chainlet
	if msg.IsServiceChainlet {
		if !admin {
			return &types.MsgLaunchChainletResponse{}, types.ErrUnauthorized
		}

		chainlet.Tags = msg.Tags

		err = k.NewChainlet(ctx, chainlet)
		if err != nil {
			return &types.MsgLaunchChainletResponse{}, err
		}

		return &types.MsgLaunchChainletResponse{}, ctx.EventManager().EmitTypedEvent(&types.EventLaunchChainlet{
			ChainName:    chainlet.ChainletName,
			Launcher:     chainlet.Launcher,
			ChainId:      chainlet.ChainId,
			Stack:        chainlet.ChainletStackName,
			StackVersion: chainlet.ChainletStackVersion,
		})
	}

	// logic to launch non-service chainlets
	epochfee, err := sdk.ParseCoinNormalized(stack.Fees.EpochFee)
	if err != nil {
		return &types.MsgLaunchChainletResponse{}, types.ErrInvalidCoin
	}
	setupfee, err := sdk.ParseCoinNormalized(stack.Fees.SetupFee)
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
	err = k.billingKeeper.BillAccount(ctx, totalFee, chainlet, stack.Fees.EpochLength, "launching chainlet")
	if err != nil {
		return &types.MsgLaunchChainletResponse{}, cosmossdkerrors.Wrapf(types.ErrBillingFailure, "failed to bill new account %s", err.Error())
	}

	// Add as a CCV consumer if enabled
	if chainlet.IsCCVConsumer {
		err = k.addConsumer(ctx, chainlet.ChainId, chainlet.SpawnTime)
		if err != nil {
			return nil, err
		}
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

func (k Keeper) validators(ctx sdk.Context) []string {
	validators, err := k.stakingKeeper.GetAllValidators(ctx)
	if err != nil {
		panic(err)
	}
	addresses := make([]string, 0, len(validators))
	for _, val := range validators {
		if val.GetStatus() != stakingtypes.Bonded {
			continue
		}

		//TODO remove: temporary hack to support acc addresses used in SagaOS start.sh
		addr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
		if err != nil {
			panic(err)
		}

		addresses = append(addresses, sdk.AccAddress(addr).String())
	}

	return addresses
}
