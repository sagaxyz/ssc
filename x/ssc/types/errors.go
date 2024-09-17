package types

// DONTCOVER

import (
	cosmossdkerrors "cosmossdk.io/errors"
)

// x/ssc module sentinel errors
var (
	ErrSample = cosmossdkerrors.Register(ModuleName, 1100, "sample error")
)
