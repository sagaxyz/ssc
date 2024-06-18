package types

// DONTCOVER

import (
	cosmossdkerrors "cosmossdk.io/errors"
)

// x/chainlet module sentinel errors
var (
	ErrInvalidEpoch            = cosmossdkerrors.Register(ModuleName, 6900, "invalid epoch length")
	ErrInvalidCoin             = cosmossdkerrors.Register(ModuleName, 6901, "invalid coin")
	ErrInvalidChainletStack    = cosmossdkerrors.Register(ModuleName, 6902, "invalid chainlet stack")
	ErrInvalidChainId          = cosmossdkerrors.Register(ModuleName, 6903, "invalid chain id")
	ErrBillingFailure          = cosmossdkerrors.Register(ModuleName, 6904, "billing failure")
	ErrChainletCreationFailure = cosmossdkerrors.Register(ModuleName, 6905, "failed to create chainlet")
	ErrChainletExists          = cosmossdkerrors.Register(ModuleName, 6906, "chainlet already exists")
	ErrJSONMarhsal             = cosmossdkerrors.Register(ModuleName, 6907, "error marshalling json")
	ErrChainletStartFailure    = cosmossdkerrors.Register(ModuleName, 6908, "failed to start chainlet")
	ErrTooManyChainlets        = cosmossdkerrors.Register(ModuleName, 6909, "chainlet limit exceeded")
)
