package epochs_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	keepertest "github.com/sagaxyz/ssc/testutil/keeper"
	"github.com/sagaxyz/ssc/x/epochs"
	"github.com/sagaxyz/ssc/x/epochs/types"
)

func TestEpochsExportGenesis(t *testing.T) {
	epochsKeeper, ctx := keepertest.EpochsKeeper(t)
	chainStartTime := ctx.BlockTime()
	chainStartHeight := ctx.BlockHeight()

	genesis := epochs.ExportGenesis(ctx, *epochsKeeper)
	require.Len(t, genesis.Epochs, 4)

	require.Equal(t, genesis.Epochs[0].Identifier, types.DayEpochID)
	require.Equal(t, genesis.Epochs[0].StartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[0].Duration, time.Hour*24)
	require.Equal(t, genesis.Epochs[0].CurrentEpoch, int64(0))
	require.Equal(t, genesis.Epochs[0].CurrentEpochStartHeight, chainStartHeight)
	require.Equal(t, genesis.Epochs[0].CurrentEpochStartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[0].EpochCountingStarted, false)

	require.Equal(t, genesis.Epochs[1].Identifier, types.HourEpochID)
	require.Equal(t, genesis.Epochs[1].StartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[1].Duration, time.Hour*1)
	require.Equal(t, genesis.Epochs[1].CurrentEpoch, int64(0))
	require.Equal(t, genesis.Epochs[1].CurrentEpochStartHeight, chainStartHeight)
	require.Equal(t, genesis.Epochs[1].CurrentEpochStartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[1].EpochCountingStarted, false)

	require.Equal(t, genesis.Epochs[2].Identifier, types.MinuteEpochID)
	require.Equal(t, genesis.Epochs[2].StartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[2].Duration, time.Minute)
	require.Equal(t, genesis.Epochs[2].CurrentEpoch, int64(0))
	require.Equal(t, genesis.Epochs[2].CurrentEpochStartHeight, chainStartHeight)
	require.Equal(t, genesis.Epochs[2].CurrentEpochStartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[2].EpochCountingStarted, false)

	require.Equal(t, genesis.Epochs[3].Identifier, types.WeekEpochID)
	require.Equal(t, genesis.Epochs[3].StartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[3].Duration, time.Hour*24*7)
	require.Equal(t, genesis.Epochs[3].CurrentEpoch, int64(0))
	require.Equal(t, genesis.Epochs[3].CurrentEpochStartHeight, chainStartHeight)
	require.Equal(t, genesis.Epochs[3].CurrentEpochStartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[3].EpochCountingStarted, false)
}

func TestEpochsInitGenesis(t *testing.T) {
	epochsKeeper, ctx := keepertest.EpochsKeeper(t)

	// On init genesis, default epochs information is set
	// To check init genesis again, should make it fresh status
	epochInfos := epochsKeeper.AllEpochInfos(ctx)
	for _, epochInfo := range epochInfos {
		epochsKeeper.DeleteEpochInfo(ctx, epochInfo.Identifier)
	}

	now := time.Now().UTC()
	ctx = ctx.WithBlockHeight(1)
	ctx = ctx.WithBlockTime(now)

	// test genesisState validation
	genesisState := types.GenesisState{
		Epochs: []types.EpochInfo{
			{
				Identifier:              "monthly",
				StartTime:               time.Time{},
				Duration:                time.Hour * 24,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: ctx.BlockHeight(),
				CurrentEpochStartTime:   time.Time{},
				EpochCountingStarted:    true,
			},
			{
				Identifier:              "monthly",
				StartTime:               time.Time{},
				Duration:                time.Hour * 24,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: ctx.BlockHeight(),
				CurrentEpochStartTime:   time.Time{},
				EpochCountingStarted:    true,
			},
		},
	}
	require.EqualError(t, genesisState.Validate(), "epoch identifier should be unique")

	genesisState = types.GenesisState{
		Epochs: []types.EpochInfo{
			{
				Identifier:              types.DayEpochID,
				StartTime:               time.Time{},
				Duration:                time.Hour * 24,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: ctx.BlockHeight(),
				CurrentEpochStartTime:   time.Time{},
				EpochCountingStarted:    true,
			},
		},
	}

	epochs.InitGenesis(ctx, *epochsKeeper, genesisState)
	epochInfo := epochsKeeper.GetEpochInfo(ctx, types.DayEpochID)
	require.NotNil(t, epochInfo)
	require.Equal(t, epochInfo.Identifier, types.DayEpochID)
	require.Equal(t, epochInfo.StartTime.UTC().String(), now.UTC().String())
	require.Equal(t, epochInfo.Duration, time.Hour*24)
	require.Equal(t, epochInfo.CurrentEpoch, int64(0))
	require.Equal(t, epochInfo.CurrentEpochStartHeight, ctx.BlockHeight())
	require.Equal(t, epochInfo.CurrentEpochStartTime.UTC().String(), time.Time{}.String())
	require.Equal(t, epochInfo.EpochCountingStarted, true)
}
