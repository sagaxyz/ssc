package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdkcodec "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/sagaxyz/ssc/testutil/keeper"
	"github.com/sagaxyz/ssc/x/epochs/keeper"
	"github.com/sagaxyz/ssc/x/epochs/types"
)

func TestEpochInfo(t *testing.T) {
	var (
		expRes *types.QueryEpochsInfoResponse
	)

	// setup
	epochsKeeper, ctx := keepertest.EpochsKeeper(t)
	day := newEpochInfo(ctx, types.DayEpochID, time.Hour*24)
	hour := newEpochInfo(ctx, types.HourEpochID, time.Hour)
	minute := newEpochInfo(ctx, types.MinuteEpochID, time.Minute)
	week := newEpochInfo(ctx, types.WeekEpochID, time.Hour*24*7)

	testCases := []struct {
		name     string
		malleate func() error
		expPass  bool
	}{
		{
			"default EpochInfos",
			func() error {
				expRes = &types.QueryEpochsInfoResponse{
					Epochs: []types.EpochInfo{day, hour, minute, week},
				}
				return nil
			},
			true,
		},
		{
			"set epoch info",
			func() error {
				quarter := newEpochInfo(ctx, "quarter", time.Hour*24*7*13)
				err := epochsKeeper.AddEpochInfo(ctx, quarter)
				if err != nil {
					return err
				}
				expRes = &types.QueryEpochsInfoResponse{
					Epochs: []types.EpochInfo{day, hour, minute, quarter, week},
				}
				return nil
			},
			true,
		},
	}
	for _, tc := range testCases {
		err := tc.malleate()
		require.NoError(t, err)
		res, err := queryClient(*epochsKeeper, ctx).EpochInfos(ctx, &types.QueryEpochsInfoRequest{})
		if tc.expPass {
			require.NoError(t, err)
			require.Equal(t, expRes, res)
		} else {
			require.NotNil(t, err)
		}
	}
}

func newEpochInfo(ctx sdk.Context, epochType string, duration time.Duration) types.EpochInfo {
	epochInfo := types.EpochInfo{
		Identifier:              epochType,
		StartTime:               time.Time{},
		Duration:                duration,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 1,
		CurrentEpochStartTime:   time.Time{},
		EpochCountingStarted:    false,
	}
	epochInfo.StartTime = ctx.BlockTime()
	epochInfo.CurrentEpochStartHeight = ctx.BlockHeight()
	return epochInfo
}

func queryClient(epochsKeeper keeper.Keeper, ctx sdk.Context) types.QueryClient {
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, sdkcodec.NewInterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQuerier(epochsKeeper))
	return types.NewQueryClient(queryHelper)
}

func TestCurrentEpoch(t *testing.T) {
	var (
		req    *types.QueryCurrentEpochRequest
		expRes *types.QueryCurrentEpochResponse
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"unknown identifier",
			func() {
				defaultCurrentEpoch := int64(0)
				req = &types.QueryCurrentEpochRequest{Identifier: "second"}
				expRes = &types.QueryCurrentEpochResponse{
					CurrentEpoch: defaultCurrentEpoch,
				}
			},
			false,
		},
		{
			"week - default currentEpoch",
			func() {
				defaultCurrentEpoch := int64(0)
				req = &types.QueryCurrentEpochRequest{Identifier: types.WeekEpochID}
				expRes = &types.QueryCurrentEpochResponse{
					CurrentEpoch: defaultCurrentEpoch,
				}
			},
			true,
		},
		{
			"day - default currentEpoch",
			func() {
				defaultCurrentEpoch := int64(0)
				req = &types.QueryCurrentEpochRequest{Identifier: types.DayEpochID}
				expRes = &types.QueryCurrentEpochResponse{
					CurrentEpoch: defaultCurrentEpoch,
				}
			},
			true,
		},
	}
	for _, tc := range testCases {
		tc.malleate()
		epochsKeeper, ctx := keepertest.EpochsKeeper(t)
		res, err := keeper.NewQuerier(*epochsKeeper).CurrentEpoch(ctx, req)
		if tc.expPass {
			require.NoError(t, err)
			require.Equal(t, expRes, res)
		} else {
			require.NotNil(t, err)
		}
	}
}
