package keeper_test

import (
	"testing"

	chainlettypes "github.com/sagaxyz/saga-sdk/x/chainlet/types"
	keepertest "github.com/sagaxyz/ssc/testutil/keeper"

	"github.com/stretchr/testify/require"
)

func TestKeeper_SetPort(t *testing.T) {
	k, ctx := keepertest.ChainletKeeper(t)

	// Test setting a port
	portID := chainlettypes.PortID
	k.SetPort(ctx, portID)

	// Verify the port was set
	retrievedPort := k.GetPort(ctx)
	require.Equal(t, portID, retrievedPort)
}

func TestKeeper_GetPort(t *testing.T) {
	k, ctx := keepertest.ChainletKeeper(t)

	// Test getting port when not set (should return empty string)
	port := k.GetPort(ctx)
	require.Empty(t, port)

	// Set a port
	portID := chainlettypes.PortID
	k.SetPort(ctx, portID)

	// Verify we can retrieve it
	retrievedPort := k.GetPort(ctx)
	require.Equal(t, portID, retrievedPort)
}

func TestKeeper_PortPersistence(t *testing.T) {
	k, ctx := keepertest.ChainletKeeper(t)

	// Set a port
	portID := chainlettypes.PortID
	k.SetPort(ctx, portID)

	// Create a new context (simulating a new block)
	ctx2 := ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// Verify port persists across contexts
	retrievedPort := k.GetPort(ctx2)
	require.Equal(t, portID, retrievedPort)
}

func TestKeeper_PortUpdate(t *testing.T) {
	k, ctx := keepertest.ChainletKeeper(t)

	// Set initial port
	initialPort := "chainlet"
	k.SetPort(ctx, initialPort)
	require.Equal(t, initialPort, k.GetPort(ctx))

	// Update to a different port
	newPort := chainlettypes.PortID
	k.SetPort(ctx, newPort)
	require.Equal(t, newPort, k.GetPort(ctx))
}

