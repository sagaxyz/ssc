package types

// DONTCOVER

import (
	cosmossdkerrors "cosmossdk.io/errors"
)

// x/gmp module sentinel errors
var (
	ErrSample               = cosmossdkerrors.Register(ModuleName, 1100, "sample error")
	ErrInvalidPacketTimeout = cosmossdkerrors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion       = cosmossdkerrors.Register(ModuleName, 1501, "invalid version")
)
