package types

// DONTCOVER

import (
	cosmossdkerrors "cosmossdk.io/errors"
)

// x/escrow module sentinel errors
var (
	ErrInvalidCoin             = cosmossdkerrors.Register(ModuleName, 5699, "Invalid coin amount.")
	ErrFunderNotFound          = cosmossdkerrors.Register(ModuleName, 5700, "Funder not found.")
	ErrBankFailure             = cosmossdkerrors.Register(ModuleName, 5701, "Bank failure.")
	ErrInvalidDenom            = cosmossdkerrors.Register(ModuleName, 5702, "Invalid denom.")
	ErrInsufficientBalance     = cosmossdkerrors.Register(ModuleName, 5703, "Insufficient balance.")
	ErrChainletAccountNotFound = cosmossdkerrors.Register(ModuleName, 5704, "Chainlet account not found.")
	ErrUnauthorized            = cosmossdkerrors.Register(ModuleName, 5705, "Unauthorized action.")
	ErrInvalidParams           = cosmossdkerrors.Register(ModuleName, 5706, "Invalid parameters.")
)
