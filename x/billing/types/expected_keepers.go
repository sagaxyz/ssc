package types

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	chainlettypes "github.com/sagaxyz/ssc/x/chainlet/types"
	epochstypes "github.com/sagaxyz/ssc/x/epochs/types"
	escrowtypes "github.com/sagaxyz/ssc/x/escrow/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	ValidateBalance(ctx sdk.Context, addr sdk.AccAddress) error
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool

	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetAccountsBalances(ctx sdk.Context) []banktypes.Balance
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	IterateAccountBalances(ctx sdk.Context, addr sdk.AccAddress, cb func(coin sdk.Coin) (stop bool))
	IterateAllBalances(ctx sdk.Context, cb func(address sdk.AccAddress, coin sdk.Coin) (stop bool))

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

type EscrowKeeper interface {
	BillAccount(ctx sdk.Context, amount sdk.Coin, chainId string, toModule string) error
	GetKprChainletAccount(ctx sdk.Context, chainId string) (acc escrowtypes.ChainletAccount, err error)
}

type StakingKeeper interface {
	GetValidators(ctx sdk.Context, maxRetrieve uint32) (validators []stakingtypes.Validator)
	GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator
}

type EpochsKeeper interface {
	GetEpochInfo(ctx sdk.Context, identifier string) epochstypes.EpochInfo
}

type ChainletKeeper interface {
	ListChainlets(goCtx context.Context, req *chainlettypes.QueryListChainletsRequest) (*chainlettypes.QueryListChainletsResponse, error)
	ListChainletStack(goCtx context.Context, req *chainlettypes.QueryListChainletStackRequest) (*chainlettypes.QueryListChainletStackResponse, error)
	StopChainlet(ctx sdk.Context, chainId string) error
	GetChainlet(goCtx context.Context, req *chainlettypes.QueryGetChainletRequest) (*chainlettypes.QueryGetChainletResponse, error)
	ChainletExists(ctx sdk.Context, chainId string) bool
	StartExistingChainlet(ctx sdk.Context, chainId string) error
	IsChainletStarted(ctx sdk.Context, chainId string) (bool, error)
	GetChainletStackInfo(ctx sdk.Context, chainId string) (*chainlettypes.ChainletStack, error)
	GetChainletInfo(ctx sdk.Context, chainId string) (*chainlettypes.Chainlet, error)
	GetParams(ctx sdk.Context) chainlettypes.Params
}
