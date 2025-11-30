package types

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	chainlettypes "github.com/sagaxyz/ssc/x/chainlet/types"
	epochstypes "github.com/sagaxyz/ssc/x/epochs/types"
	escrowtypes "github.com/sagaxyz/ssc/x/escrow/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	ValidateBalance(ctx context.Context, addr sdk.AccAddress) error
	HasBalance(ctx context.Context, addr sdk.AccAddress, amt sdk.Coin) bool

	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetAccountsBalances(ctx context.Context) []banktypes.Balance
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	LockedCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins

	IterateAccountBalances(ctx context.Context, addr sdk.AccAddress, cb func(coin sdk.Coin) (stop bool))
	IterateAllBalances(ctx context.Context, cb func(address sdk.AccAddress, coin sdk.Coin) (stop bool))

	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

type EscrowKeeper interface {
	BillAccount(ctx sdk.Context, amount sdk.Coin, chainId string, toModule string) error
	GetChainletWithPools(ctx sdk.Context, chainId string) (acc escrowtypes.ChainletAccount, pool []*escrowtypes.DenomPool, err error)
}

type StakingKeeper interface {
	GetValidators(ctx context.Context, maxRetrieve uint32) (validators []stakingtypes.Validator, err error)
}

type EpochsKeeper interface {
	GetEpochInfo(ctx sdk.Context, identifier string) epochstypes.EpochInfo
}

type ChainletKeeper interface {
	ListChainlets(ctx context.Context, req *chainlettypes.QueryListChainletsRequest) (*chainlettypes.QueryListChainletsResponse, error)
	ListChainletStack(ctx context.Context, req *chainlettypes.QueryListChainletStackRequest) (*chainlettypes.QueryListChainletStackResponse, error)
	StopChainlet(ctx sdk.Context, chainId string) error
	GetChainlet(ctx context.Context, req *chainlettypes.QueryGetChainletRequest) (*chainlettypes.QueryGetChainletResponse, error)
	ChainletExists(ctx sdk.Context, chainId string) bool
	StartExistingChainlet(ctx sdk.Context, chainId string) error
	IsChainletStarted(ctx sdk.Context, chainId string) (bool, error)
	GetChainletStackInfo(ctx sdk.Context, chainId string) (*chainlettypes.ChainletStack, error)
	GetChainletInfo(ctx sdk.Context, chainId string) (*chainlettypes.Chainlet, error)
	GetParams(ctx sdk.Context) chainlettypes.Params
}

type BillingKeeper interface {
	GetPlatformValidators(ctx sdk.Context) []string
}
