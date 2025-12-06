package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/chainlet module sentinel errors
var (
	ErrInvalidDenom            = sdkerrors.Register(ModuleName, 6900, "invalid denom")
	ErrInvalidCoin             = sdkerrors.Register(ModuleName, 6901, "invalid coin")
	ErrInvalidChainletStack    = sdkerrors.Register(ModuleName, 6902, "invalid chainlet stack")
	ErrInvalidChainId          = sdkerrors.Register(ModuleName, 6903, "invalid chain id")
	ErrBillingFailure          = sdkerrors.Register(ModuleName, 6904, "billing failure")
	ErrChainletCreationFailure = sdkerrors.Register(ModuleName, 6905, "failed to create chainlet")
	ErrChainletExists          = sdkerrors.Register(ModuleName, 6906, "chainlet already exists")
	ErrJSONMarhsal             = sdkerrors.Register(ModuleName, 6907, "error marshalling json")
	ErrChainletStartFailure    = sdkerrors.Register(ModuleName, 6908, "failed to start chainlet")
	ErrTooManyChainlets        = sdkerrors.Register(ModuleName, 6909, "chainlet limit exceeded")
	ErrUnauthorized            = sdkerrors.Register(ModuleName, 6910, "not authorized to launch a service chainlet")
	ErrInvalidPacketTimeout    = sdkerrors.Register(ModuleName, 6911, "invalid packet timeout")
	ErrInvalidVersion          = sdkerrors.Register(ModuleName, 6912, "invalid version")
	ErrInvalidFees             = sdkerrors.Register(ModuleName, 6913, "invalid fees")
	ErrDuplicateDenom          = sdkerrors.Register(ModuleName, 6914, "duplicate denom in fees")
	ErrNoUpgradeInProgress     = sdkerrors.Register(ModuleName, 6915, "no upgrade in progress")
)
