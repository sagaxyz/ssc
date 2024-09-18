package keeper

import (
	"encoding/json"
	"fmt"

	cosmossdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/escrow/types"
)

func (k Keeper) ChainletAccountStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ChainletAccKey))
}
func (k Keeper) UserAccountStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.UserAccKey))
}

func (k Keeper) NewChainletAccount(ctx sdk.Context, address sdk.AccAddress, chainId string, depositAmount sdk.Coin) error {
	store := k.ChainletAccountStore(ctx)
	key := []byte(chainId)
	if store.Has(key) {
		return fmt.Errorf("chainlet account already exists")
	}

	if depositAmount.Denom != k.GetParams(ctx).SupportedDenom {
		return cosmossdkerrors.Wrapf(types.ErrInvalidDenom, "invalid denom %s, only %s is allowed", depositAmount.Denom, k.GetParams(ctx).SupportedDenom)
	}
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, sdk.NewCoins(depositAmount))
	if err != nil {
		return cosmossdkerrors.Wrapf(types.ErrInsufficientBalance, "insufficient balance in account %s for deposit: %s required", address, depositAmount.String())
	}

	shares, err := sdk.ParseDecCoin(depositAmount.String())
	if err != nil {
		return err
	}

	chainlet := types.ChainletAccount{
		ChainId: chainId,
		Balance: depositAmount,
		Shares:  shares.Amount,
		Funders: map[string]types.Funder{},
	}

	chainlet.Funders[address.String()] = types.Funder{
		Shares: shares.Amount,
	}
	err = k.SetChainletAccount(ctx, chainlet)
	if err != nil {
		return err
	}

	ctx.Logger().Info(fmt.Sprintf("successfully created a new chainlet escrow account for %s. balance: %s", chainlet.ChainId, chainlet.Balance))
	return nil
}

func (k Keeper) deposit(ctx sdk.Context, address sdk.AccAddress, amount sdk.Coin, chainId string) error {
	chainlet, err := k.GetKprChainletAccount(ctx, chainId)
	if err != nil {
		return err
	}

	if chainlet.Balance.Amount.Equal(math.ZeroInt()) {
		chainlet.Shares = math.LegacyZeroDec()
		chainlet.Funders = make(map[string]types.Funder)
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, sdk.NewCoins(amount))
	if err != nil {
		return fmt.Errorf("failed to send coins from account to module: %v", err)
	}

	// Calculation of shares: https://docs.cosmos.network/v0.47/modules/staking#how-shares-are-calculated
	// S_j = S * T_j / T
	// Implementation Example: https://github.com/cosmos/cosmos-sdk/blob/bfba5491f39f0e0af100480a3194a30c2dc4b9c3/x/staking/types/validator.go#L326
	var newShares math.LegacyDec
	if chainlet.Balance.IsPositive() {
		newShares = chainlet.Shares.MulInt(amount.Amount).QuoInt(chainlet.Balance.Amount)
	} else if chainlet.Balance.IsZero() {
		coinShares, err := sdk.ParseDecCoin(amount.String())
		if err != nil {
			return fmt.Errorf("failed to parse deposit amount as DecCoin: %v", err)
		}
		newShares = coinShares.Amount
	}

	addrString := address.String()

	_, ok := chainlet.Funders[addrString]

	if ok {
		chainlet.Funders[addrString] = types.Funder{
			Shares: chainlet.Funders[addrString].Shares.Add(newShares),
		}
	} else {
		chainlet.Funders[addrString] = types.Funder{
			Shares: newShares,
		}
	}

	chainlet.Shares = chainlet.Shares.Add(newShares)
	chainlet.Balance = chainlet.Balance.Add(amount)
	err = k.SetChainletAccount(ctx, chainlet)
	if err != nil {
		return err
	}

	err = k.billingKeeper.BillAndRestartChainlet(ctx, chainId)
	if err != nil {
		return err
	}

	return ctx.EventManager().EmitTypedEvent(&types.EventDeposit{ //nolint:errcheck
		User:     addrString,
		Chainlet: chainlet.ChainId,
		Amount:   amount.String(),
		NewTotal: chainlet.Balance.String(),
	})
}

func (k Keeper) withdraw(ctx sdk.Context, address sdk.AccAddress, chainId string) error {
	chainlet, err := k.GetKprChainletAccount(ctx, chainId)
	if err != nil {
		return err
	}

	sf := ScalingFactor(chainlet)
	addr := address.String()

	c, ok := chainlet.Funders[addr]
	if !ok {
		return cosmossdkerrors.Wrap(types.ErrFunderNotFound, addr)
	}
	tokens := c.Shares.Quo(sf)
	coins := sdk.NewCoin(k.GetParams(ctx).SupportedDenom, tokens.RoundInt())
	newChainletBalance, err := chainlet.Balance.SafeSub(coins)
	if err != nil {
		return err
	}
	ctx.Logger().Info(fmt.Sprintf("sending %s coins to %s", coins.String(), addr))

	toaddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toaddr, sdk.NewCoins(coins))
	if err != nil {
		return err
	}

	delete(chainlet.Funders, addr)
	chainlet.Balance = newChainletBalance

	chainlet.Shares = chainlet.Shares.Sub(c.Shares)

	err = k.SetChainletAccount(ctx, chainlet)
	if err != nil {
		return err
	}

	return ctx.EventManager().EmitTypedEvent(&types.EventWithdraw{
		User:      addr,
		Chainlet:  chainlet.ChainId,
		Amount:    coins.String(),
		Remaining: chainlet.Balance.String(),
	})
}

func (k Keeper) GetKprChainletAccount(ctx sdk.Context, chainId string) (acc types.ChainletAccount, err error) {
	store := k.ChainletAccountStore(ctx)
	key := []byte(chainId)

	if !store.Has(key) {
		return acc, cosmossdkerrors.Wrapf(types.ErrChainletAccountNotFound, "chainlet escrow account for chain %s not found", chainId)
	}

	account := store.Get(key)
	err = json.Unmarshal(account, &acc)
	if err != nil {
		return acc, fmt.Errorf("failed to unmarshal chainlet %s", chainId)
	}
	return
}

func (k Keeper) SetChainletAccount(ctx sdk.Context, chainlet types.ChainletAccount) error {
	store := k.ChainletAccountStore(ctx)
	key := []byte(chainlet.ChainId)

	bz, err := json.Marshal(chainlet)
	if err != nil {
		return err
	}

	store.Set(key, bz)
	return nil
}

func (k Keeper) checkEnoughBalance(ctx sdk.Context, chainId string, amount sdk.Coin) (sdk.Coin, error) {
	acc, err := k.GetKprChainletAccount(ctx, chainId)
	if err != nil {
		return amount, err
	}
	return acc.Balance.SafeSub(amount)
}

// logic for billing
func (k Keeper) BillAccount(ctx sdk.Context, amount sdk.Coin, chainId string, toModule string) error {

	acc, err := k.GetKprChainletAccount(ctx, chainId)
	if err != nil {
		return err
	}
	nb, err := k.checkEnoughBalance(ctx, chainId, amount)
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, toModule, sdk.NewCoins(amount))
	if err != nil {
		return cosmossdkerrors.Wrap(types.ErrBankFailure, err.Error())
	}
	acc.Balance = nb
	err = k.SetChainletAccount(ctx, acc)
	if err != nil {
		return err
	}
	return nil
}

func ScalingFactor(chainlet types.ChainletAccount) math.LegacyDec {
	totalPoolShares := chainlet.Shares
	decCoins, _ := sdk.ParseDecCoin(chainlet.Balance.String())
	totalDeposit := decCoins.Amount
	if totalDeposit.IsZero() {
		return math.LegacyMustNewDecFromStr("1")
	}
	return totalPoolShares.Quo(totalDeposit)
}

func InverseScalingFactor(chainlet types.ChainletAccount) math.LegacyDec {
	return math.LegacyNewDec(int64(1)).Quo(ScalingFactor(chainlet))
}
