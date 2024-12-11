package types

// DONTCOVER

import (
	cosmossdkerrors "cosmossdk.io/errors"
)

// x/gmp module sentinel errors
var (
	ErrSample = cosmossdkerrors.Register(ModuleName, 1100, "sample error")
)
