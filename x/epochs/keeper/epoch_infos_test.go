package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	keepertest "github.com/sagaxyz/ssc/testutil/keeper"
	"github.com/sagaxyz/ssc/x/epochs/types"
)

func TestEpochInfos(t *testing.T) {

	epochInfo := types.EpochInfo{
		Identifier:            "monthly",
		StartTime:             time.Time{},
		Duration:              time.Hour * 24 * 30,
		CurrentEpoch:          0,
		CurrentEpochStartTime: time.Time{},
		EpochCountingStarted:  false,
	}

	epochsKeeper, ctx := keepertest.EpochsKeeper(t)

	err := epochsKeeper.AddEpochInfo(ctx, epochInfo)
	if err != nil {
		t.FailNow()
	}

	epochInfoSaved := epochsKeeper.GetEpochInfo(ctx, "monthly")
	require.NotNil(t, epochInfoSaved)
	require.Equal(t, epochInfo, epochInfoSaved)

	allEpochs := epochsKeeper.AllEpochInfos(ctx)
	require.Len(t, allEpochs, 5)
	require.Equal(t, allEpochs[0].Identifier, types.DayEpochID) // alphabetical order
	require.Equal(t, allEpochs[1].Identifier, types.HourEpochID)
	require.Equal(t, allEpochs[2].Identifier, types.MinuteEpochID)
	require.Equal(t, allEpochs[3].Identifier, "monthly")
	require.Equal(t, allEpochs[4].Identifier, types.WeekEpochID)
}
