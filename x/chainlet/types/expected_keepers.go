package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v10/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ccvprovidertypes "github.com/cosmos/interchain-security/v7/x/ccv/provider/types"
	ccvtypes "github.com/cosmos/interchain-security/v7/x/ccv/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
	// Methods imported from account should be defined here
}

type BankKeeper interface {
}

type StakingKeeper interface {
	GetAllValidators(ctx context.Context) (validators []stakingtypes.Validator, err error)
}

type ProviderKeeper interface {
	AppendPendingVSCPackets(ctx sdk.Context, consumerID string, newPackets ...ccvtypes.ValidatorSetChangePacketData)
	GetValidatorSetUpdateId(ctx sdk.Context) (validatorSetUpdateId uint64)
	IncrementValidatorSetUpdateId(ctx sdk.Context)
	GetConsumerIdToChannelId(ctx sdk.Context, consumerId string) (string, bool)
	SendVSCPacketsToChain(ctx sdk.Context, consumerID string, channelID string) error
	GetConsumerPhase(ctx sdk.Context, consumerID string) ccvprovidertypes.ConsumerPhase
	GetConsumerClientId(ctx sdk.Context, chainID string) (string, bool)
}

type ProviderMsgServer interface {
	CreateConsumer(goCtx context.Context, msg *ccvprovidertypes.MsgCreateConsumer) (*ccvprovidertypes.MsgCreateConsumerResponse, error)
}

type ClientKeeper interface {
	GetClientLatestHeight(sdk.Context, string) clienttypes.Height
}
type ChannelKeeper interface {
	GetChannel(sdk.Context, string, string) (ibcchanneltypes.Channel, bool)
}
type ConnectionKeeper interface {
	GetConnection(sdk.Context, string) (ibcconnectiontypes.ConnectionEnd, bool)
}

type BillingKeeper interface {
	BillAccount(ctx sdk.Context, amount sdk.Coin, chainlet Chainlet, memo string) error
	PayEpochFeeToValidator(ctx sdk.Context, epochFee sdk.Coins, fromModuleName string, valAddr sdk.AccAddress, memo string) (err error)
}

type EscrowKeeper interface {
	NewChainletAccount(ctx sdk.Context, address sdk.AccAddress, chainId string, depositAmount sdk.Coin) error
	GetSupportedDenoms(ctx sdk.Context) []string
}

type AclKeeper interface {
	Allowed(ctx sdk.Context, addr sdk.AccAddress) bool
	IsAdmin(ctx sdk.Context, addr sdk.AccAddress) bool
}
