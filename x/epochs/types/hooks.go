package types

import (
	"errors"
	fmt "fmt"
	"cosmossdk.io/store/types"
	"runtime"
	"runtime/debug"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type EpochHooks interface {
	// the first block whose timestamp is after the duration is counted as the end of the epoch
	AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error
	// new epoch is next block of epoch end block
	BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error
}

var _ EpochHooks = MultiEpochHooks{}

// combine multiple gamm hooks, all hook functions are run in array sequence.
type MultiEpochHooks []EpochHooks

func NewMultiEpochHooks(hooks ...EpochHooks) MultiEpochHooks {
	return hooks
}

// AfterEpochEnd is called when epoch is going to be ended, epochNumber is the number of epoch that is ending.
func (h MultiEpochHooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	for i := range h {
		panicCatchingEpochHook(ctx, h[i].AfterEpochEnd, epochIdentifier, epochNumber)
	}
	return nil
}

// BeforeEpochStart is called when epoch is going to be started, epochNumber is the number of epoch that is starting.
func (h MultiEpochHooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	for i := range h {
		panicCatchingEpochHook(ctx, h[i].BeforeEpochStart, epochIdentifier, epochNumber)
	}
	return nil
}

func panicCatchingEpochHook(
	ctx sdk.Context,
	hookFn func(ctx sdk.Context, epochIdentifier string, epochNumber int64) error,
	epochIdentifier string,
	epochNumber int64,
) {
	wrappedHookFn := func(ctx sdk.Context) error {
		return hookFn(ctx, epochIdentifier, epochNumber)
	}
	// TODO: Thread info for which hook this is, may be dependent on larger hook system refactoring
	err := applyFuncIfNoError(ctx, wrappedHookFn)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("error in epoch hook %v", err))
	}
}

func applyFuncIfNoError(ctx sdk.Context, f func(ctx sdk.Context) error) (err error) {
	// Add a panic safeguard
	defer func() {
		if recoveryError := recover(); recoveryError != nil {
			printPanicRecoveryError(ctx, recoveryError)
			err = errors.New("panic occurred during execution")
		}
	}()
	// makes a new cache context, which all state changes get wrapped inside of.
	cacheCtx, write := ctx.CacheContext()
	err = f(cacheCtx)
	if err != nil {
		ctx.Logger().Error(err.Error())
	} else {
		// no error, write the output of f
		write()
		ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
	}
	return err
}

func printPanicRecoveryError(ctx sdk.Context, recoveryError interface{}) {
	errStackTrace := string(debug.Stack())
	switch e := recoveryError.(type) {
	case types.ErrorOutOfGas:
		ctx.Logger().Debug("out of gas error inside panic recovery block: " + e.Descriptor)
		return
	case string:
		ctx.Logger().Error("Recovering from (string) panic: " + e)
	case runtime.Error:
		ctx.Logger().Error("recovered (runtime.Error) panic: " + e.Error())
	case error:
		ctx.Logger().Error("recovered (error) panic: " + e.Error())
	default:
		ctx.Logger().Error("recovered (default) panic. Could not capture logs in ctx, see stdout")
		fmt.Println("Recovering from panic ", recoveryError)
		debug.PrintStack()
		return
	}
	ctx.Logger().Error("stack trace: " + errStackTrace)
}
