package types

import (
	cosmossdkerrors "cosmossdk.io/errors"
)

var (
	ErrInternalFailure        = cosmossdkerrors.Register(ModuleName, 7700, "internal failure")
	ErrNoRecords              = cosmossdkerrors.Register(ModuleName, 7701, "no records found")
	ErrJSONMarhsal            = cosmossdkerrors.Register(ModuleName, 7702, "failed to marshal json")
	ErrDuplicateRecord        = cosmossdkerrors.Register(ModuleName, 7703, "duplicate record")
	ErrInternalBillingFailure = cosmossdkerrors.Register(ModuleName, 7704, "internal failure")
)
