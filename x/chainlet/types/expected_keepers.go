package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	icatypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v10/modules/core/03-connection/types"
	ccvprovidertypes "github.com/cosmos/interchain-security/v7/x/ccv/provider/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
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
	// GetAllValidators(sdk.Context) []stakingtypes.Validator
	GetAllValidators(ctx context.Context) (validators []stakingtypes.Validator, err error)
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

type ICAKeeper interface {
	GetInterchainAccountAddress(sdk.Context, string, string) (string, bool)
	GetOpenActiveChannel(sdk.Context, string, string) (string, bool)
	SendTx(sdk.Context, *capabilitytypes.Capability, string, string, icatypes.InterchainAccountPacketData, uint64) (uint64, error)
}

type ClientKeeper interface {
	GetClientState(sdk.Context, string) (ibcexported.ClientState, bool)
}
type ChannelKeeper interface {
	GetChannel(sdk.Context, string, string) (ibcchanneltypes.Channel, bool)
}
type ConnectionKeeper interface {
	GetConnection(sdk.Context, string) (ibcconnectiontypes.ConnectionEnd, bool)
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
	IsAdmin(ctx sdk.Context, addr sdk.AccAddress) bool
}
