package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ccvprovidertypes "github.com/cosmos/interchain-security/v5/x/ccv/provider/types"
	ccvtypes "github.com/cosmos/interchain-security/v5/x/ccv/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
	// Methods imported from account should be defined here
}

type BankKeeper interface {
}

type ProviderKeeper interface {
	HandleConsumerAdditionProposal(ctx sdk.Context, prop *ccvprovidertypes.MsgConsumerAddition) error
	AppendPendingVSCPackets(ctx sdk.Context, chainID string, newPackets ...ccvtypes.ValidatorSetChangePacketData)
	GetValidatorSetUpdateId(ctx sdk.Context) (validatorSetUpdateId uint64)
	IncrementValidatorSetUpdateId(ctx sdk.Context)
	GetChainToChannel(ctx sdk.Context, chainID string) (string, bool)
	SendVSCPacketsToChain(ctx sdk.Context, chainID string, channelID string)
	GetConsumerClientId(ctx sdk.Context, chainID string) (string, bool)
}

type BillingKeeper interface {
	BillAccount(ctx sdk.Context, amount sdk.Coin, chainlet Chainlet, epochIdentifier, memo string) error
	PayEpochFeeToValidator(ctx sdk.Context, epochFee sdk.Coins, fromModuleName string, valAddr sdk.AccAddress, memo string) (err error)
}

type EscrowKeeper interface {
	NewChainletAccount(ctx sdk.Context, address sdk.AccAddress, chainId string, depositAmount sdk.Coin) error
}

type AclKeeper interface {
	Allowed(ctx sdk.Context, addr sdk.AccAddress) bool
}
