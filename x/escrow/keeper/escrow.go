package keeper

import (
	"bytes"
	"fmt"

	cosmossdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/x/escrow/types"
)

// ---------- helpers: get/set state rows ----------

func (k Keeper) getChainlet(ctx sdk.Context, chainID string) (types.ChainletAccount, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.ChainletKey(chainID))
	if bz == nil {
		return types.ChainletAccount{}, false
	}
	var acc types.ChainletAccount
	k.cdc.MustUnmarshal(bz, &acc)
	return acc, true
}

func (k Keeper) setChainlet(ctx sdk.Context, acc types.ChainletAccount) {
	bz := k.cdc.MustMarshal(&acc)
	ctx.KVStore(k.storeKey).Set(types.ChainletKey(acc.ChainId), bz)
}

func (k Keeper) getPool(ctx sdk.Context, chainID, denom string) (types.DenomPool, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.PoolKey(chainID, denom))
	if bz == nil {
		return types.DenomPool{}, false
	}
	var p types.DenomPool
	k.cdc.MustUnmarshal(bz, &p)
	return p, true
}

func (k Keeper) setPool(ctx sdk.Context, p types.DenomPool) {
	bz := k.cdc.MustMarshal(&p)
	ctx.KVStore(k.storeKey).Set(types.PoolKey(p.ChainId, p.Denom), bz)
}

func (k Keeper) getFunder(ctx sdk.Context, chainID, denom, addr string) (types.Funder, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.FunderKey(chainID, denom, addr))
	if bz == nil {
		return types.Funder{}, false
	}
	var f types.Funder
	k.cdc.MustUnmarshal(bz, &f)
	return f, true
}

func (k Keeper) setFunder(ctx sdk.Context, chainID, denom, addr string, f types.Funder) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.FunderKey(chainID, denom, addr), k.cdc.MustMarshal(&f))
	// optional reverse index for "my positions"
	store.Set(types.ByFunderKey(addr, chainID, denom), []byte{})
}

func (k Keeper) deleteFunder(ctx sdk.Context, chainID, denom, addr string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.FunderKey(chainID, denom, addr))
	// optional reverse index cleanup
	store.Delete(types.ByFunderKey(addr, chainID, denom))
}

// ---------- params / validation ----------

func (k Keeper) assertSupportedDenom(ctx sdk.Context, denom string) error {
	params := k.GetParams(ctx)
	// assuming Params.supportedDenoms []string
	for _, d := range params.SupportedDenoms {
		if d == denom {
			return nil
		}
	}
	return cosmossdkerrors.Wrapf(types.ErrInvalidDenom, "unsupported denom %s", denom)
}

// ---------- public API ----------

// Creates the chainlet head (if not exists) and bootstraps the pool/funder with the initial deposit.
func (k Keeper) NewChainletAccount(ctx sdk.Context, addr sdk.AccAddress, chainID string, deposit sdk.Coin) error {
	denom := deposit.Denom
	if err := k.assertSupportedDenom(ctx, denom); err != nil {
		return err
	}

	// Create head if not exists
	if _, ok := k.getChainlet(ctx, chainID); ok {
		return fmt.Errorf("chainlet account already exists")
	}
	k.setChainlet(ctx, types.ChainletAccount{ChainId: chainID})

	// move funds into module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.NewCoins(deposit)); err != nil {
		return cosmossdkerrors.Wrapf(types.ErrInsufficientBalance, "insufficient balance in %s for deposit: %s required", addr, deposit.String())
	}

	// pool math (bootstrap -> 1:1 shares)
	decCoin, err := sdk.ParseDecCoin(deposit.String())
	if err != nil {
		return err
	}
	newShares := decCoin.Amount

	pool := types.DenomPool{
		ChainId: chainID,
		Denom:   denom,
		Balance: deposit,
		Shares:  newShares,
	}
	k.setPool(ctx, pool)

	k.setFunder(ctx, chainID, denom, addr.String(), types.Funder{Shares: newShares})

	ctx.Logger().Info(fmt.Sprintf("created chainlet %s and pool %s; balance=%s", chainID, denom, deposit.String()))
	return ctx.EventManager().EmitTypedEvent(&types.EventDeposit{
		User:     addr.String(),
		Chainlet: chainID,
		Denom:    denom,
		Amount:   deposit.String(),
		NewTotal: pool.Balance.String(),
	})
}

// Deposit into a specific {chainID, denom} pool.
func (k Keeper) deposit(ctx sdk.Context, addr sdk.AccAddress, chainID string, amount sdk.Coin) error {
	denom := amount.Denom
	if err := k.assertSupportedDenom(ctx, denom); err != nil {
		return err
	}
	// ensure head exists
	if _, ok := k.getChainlet(ctx, chainID); !ok {
		return cosmossdkerrors.Wrapf(types.ErrChainletAccountNotFound, "chainlet %s not found", chainID)
	}

	// move funds in
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.NewCoins(amount)); err != nil {
		return fmt.Errorf("failed to send coins from account to module: %w", err)
	}

	pool, ok := k.getPool(ctx, chainID, denom)
	if !ok {
		// bootstrap pool
		dec, err := sdk.ParseDecCoin(amount.String())
		if err != nil {
			return err
		}
		pool = types.DenomPool{
			ChainId: chainID,
			Denom:   denom,
			Balance: amount,
			Shares:  dec.Amount,
		}
	} else {
		// S_j = S * T_j / T
		var newShares math.LegacyDec
		if pool.Balance.IsPositive() {
			newShares = pool.Shares.MulInt(amount.Amount).QuoInt(pool.Balance.Amount)
		} else {
			dec, err := sdk.ParseDecCoin(amount.String())
			if err != nil {
				return err
			}
			newShares = dec.Amount
		}
		pool.Shares = pool.Shares.Add(newShares)
		pool.Balance = pool.Balance.Add(amount)

		// update funder
		f, exists := k.getFunder(ctx, chainID, denom, addr.String())
		if exists {
			f.Shares = f.Shares.Add(newShares)
		} else {
			f.Shares = newShares
		}
		k.setFunder(ctx, chainID, denom, addr.String(), f)
	}

	// if pool was bootstrapped above, also set initial funder
	if _, exists := k.getFunder(ctx, chainID, denom, addr.String()); !exists {
		dec, _ := sdk.ParseDecCoin(amount.String())
		k.setFunder(ctx, chainID, denom, addr.String(), types.Funder{Shares: dec.Amount})
	}

	k.setPool(ctx, pool)

	// optional: billing hook
	if err := k.billingKeeper.BillAndRestartChainlet(ctx, chainID); err != nil {
		return err
	}

	return ctx.EventManager().EmitTypedEvent(&types.EventDeposit{
		User:     addr.String(),
		Chainlet: chainID,
		Denom:    denom,
		Amount:   amount.String(),
		NewTotal: pool.Balance.String(),
	})
}

// Withdraw the entire position for {chainID, denom} (or adapt to partials if needed).
// Withdraw all positions (all denoms) a user has on a given chainlet.
func (k Keeper) WithdrawAll(ctx sdk.Context, addr sdk.AccAddress, chainID string) error {
	store := ctx.KVStore(k.storeKey)
	addrStr := addr.String()

	// Open a substore scoped to this funder
	pstore := prefix.NewStore(store, types.ByFunderPrefix(addrStr))
	it := pstore.Iterator(nil, nil)
	defer it.Close()

	total := sdk.NewCoins()
	foundAny := false

	for ; it.Valid(); it.Next() {
		// key = "{chainId}/{denom}"
		parts := bytes.SplitN(it.Key(), []byte{'/'}, 2)
		if len(parts) != 2 {
			continue
		}
		if string(parts[0]) != chainID {
			continue
		}
		denom := string(parts[1])

		pool, ok := k.getPool(ctx, chainID, denom)
		if !ok {
			continue
		}
		f, exists := k.getFunder(ctx, chainID, denom, addrStr)
		if !exists || f.Shares.IsZero() {
			continue
		}

		coinOut, err := k.withdrawOne(ctx, addr, &pool, chainID, denom, f)
		if err != nil {
			return err
		}
		if !coinOut.IsZero() {
			total = total.Add(coinOut)
			foundAny = true
			k.setPool(ctx, pool)
			_ = ctx.EventManager().EmitTypedEvent(&types.EventWithdraw{
				User:      addrStr,
				Chainlet:  chainID,
				Denom:     denom,
				Remaining: pool.Balance.String(),
			})
		}
	}

	if !foundAny {
		return cosmossdkerrors.Wrapf(types.ErrFunderNotFound,
			"no positions for %s on chainlet %s", addrStr, chainID)
	}

	return nil
}

// OPTIONAL: allow withdrawing a single denom position, if you keep a denom-specific Msg later.
func (k Keeper) WithdrawDenom(ctx sdk.Context, addr sdk.AccAddress, chainID, denom string) error {
	addrStr := addr.String()

	pool, ok := k.getPool(ctx, chainID, denom)
	if !ok {
		return cosmossdkerrors.Wrapf(types.ErrChainletAccountNotFound, "pool %s/%s not found", chainID, denom)
	}
	f, exists := k.getFunder(ctx, chainID, denom, addrStr)
	if !exists || f.Shares.IsZero() {
		return cosmossdkerrors.Wrap(types.ErrFunderNotFound, addrStr)
	}

	_, err := k.withdrawOne(ctx, addr, &pool, chainID, denom, f)
	if err != nil {
		return err
	}

	k.setPool(ctx, pool)
	_ = ctx.EventManager().EmitTypedEvent(&types.EventWithdraw{
		User:      addrStr,
		Chainlet:  chainID,
		Denom:     denom,
		Remaining: pool.Balance.String(),
	})
	return nil
}

// Helper: withdraw ENTIRE position for a single denom.
// Returns the coin paid out for this denom and mutates 'pool' in place.
func (k Keeper) withdrawOne(
	ctx sdk.Context,
	addr sdk.AccAddress,
	pool *types.DenomPool,
	chainID, denom string,
	f types.Funder,
) (sdk.Coin, error) {
	// Defensive checks
	if pool.Shares.IsZero() || f.Shares.IsZero() {
		return sdk.NewCoin(denom, math.ZeroInt()), nil
	}

	sf := ScalingFactor(*pool) // shares -> tokens scale (sdk.Dec)
	if sf.IsZero() {
		return sdk.NewCoin(denom, math.ZeroInt()), nil
	}

	// tokens = floor(f.Shares / sf); floor to avoid overpaying by rounding.
	tokensDec := f.Shares.Quo(sf)
	amt := tokensDec.TruncateInt()
	if !amt.IsPositive() {
		return sdk.NewCoin(denom, math.ZeroInt()), nil
	}

	coin := sdk.NewCoin(denom, amt)

	newBal, err := pool.Balance.SafeSub(coin)
	if err != nil {
		return sdk.NewCoin(denom, math.ZeroInt()), err
	}

	// Transfer
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(coin)); err != nil {
		return sdk.NewCoin(denom, math.ZeroInt()), err
	}

	// State updates (also clears reverse index via deleteFunder if implemented there)
	k.deleteFunder(ctx, chainID, denom, addr.String())
	pool.Balance = newBal
	pool.Shares = pool.Shares.Sub(f.Shares)

	return coin, nil
}

// GetKprChainletAccount (kept for compatibility) — now just returns the head.
// This was previously used to get the full account with funders and pools etc.
func (k Keeper) GetKprChainletAccount(ctx sdk.Context, chainID string) (types.ChainletAccount, error) {
	acc, ok := k.getChainlet(ctx, chainID)
	if !ok {
		return types.ChainletAccount{}, cosmossdkerrors.Wrapf(types.ErrChainletAccountNotFound, "chainlet %s not found", chainID)
	}
	return acc, nil
}

// SetChainletAccount (compat) — uses protobuf codec.
func (k Keeper) SetChainletAccount(ctx sdk.Context, chainlet types.ChainletAccount) error {
	k.setChainlet(ctx, chainlet)
	return nil
}

// Check that a pool has enough balance for a debit.
func (k Keeper) checkEnoughBalance(ctx sdk.Context, chainID string, amount sdk.Coin) (sdk.Coin, error) {
	pool, ok := k.getPool(ctx, chainID, amount.Denom)
	if !ok {
		return amount, cosmossdkerrors.Wrapf(types.ErrChainletAccountNotFound, "pool %s/%s not found", chainID, amount.Denom)
	}
	return pool.Balance.SafeSub(amount)
}

// Module-to-module billing from a specific denom pool.
func (k Keeper) BillAccount(ctx sdk.Context, amount sdk.Coin, chainID, toModule string) error {
	// ensure enough balance
	nb, err := k.checkEnoughBalance(ctx, chainID, amount)
	if err != nil {
		return err
	}
	// move funds
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, toModule, sdk.NewCoins(amount)); err != nil {
		return cosmossdkerrors.Wrap(types.ErrBankFailure, err.Error())
	}
	// update pool
	pool, _ := k.getPool(ctx, chainID, amount.Denom)
	pool.Balance = nb
	k.setPool(ctx, pool)
	return nil
}

// ---------- math ----------

func ScalingFactor(pool types.DenomPool) math.LegacyDec {
	totalPoolShares := pool.Shares
	decCoins, _ := sdk.ParseDecCoin(pool.Balance.String())
	totalDeposit := decCoins.Amount
	if totalDeposit.IsZero() {
		return math.LegacyMustNewDecFromStr("1")
	}
	return totalPoolShares.Quo(totalDeposit)
}

func InverseScalingFactor(pool types.DenomPool) math.LegacyDec {
	return math.LegacyNewDec(int64(1)).Quo(ScalingFactor(pool))
}

// GetChainletWithPools returns the ChainletAccount head and all DenomPool rows
// for the given chainID. If the chainlet doesn't exist, returns NotFound.
func (k Keeper) GetChainletWithPools(
	ctx sdk.Context,
	chainID string,
) (types.ChainletAccount, []*types.DenomPool, error) {
	// Fetch chainlet head
	acc, ok := k.getChainlet(ctx, chainID)
	if !ok {
		return types.ChainletAccount{}, nil, cosmossdkerrors.Wrapf(
			types.ErrChainletAccountNotFound, "chainlet %s not found", chainID,
		)
	}

	// Collect all pools under the chainlet
	store := ctx.KVStore(k.storeKey)
	pfx := prefix.NewStore(store, types.PoolPrefix(chainID))

	it := pfx.Iterator(nil, nil)
	defer it.Close()

	pools := make([]*types.DenomPool, 0, 8) // small default cap; grows as needed
	for ; it.Valid(); it.Next() {
		p := new(types.DenomPool)
		k.cdc.MustUnmarshal(it.Value(), p)
		pools = append(pools, p)
	}

	return acc, pools, nil
}
