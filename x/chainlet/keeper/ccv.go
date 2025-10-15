package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	ccvprovidertypes "github.com/cosmos/interchain-security/v7/x/ccv/provider/types"
	ccvtypes "github.com/cosmos/interchain-security/v7/x/ccv/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (k *Keeper) setPendingInit(ctx sdk.Context, consumerID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletInit)
	store.Set([]byte(consumerID), k.cdc.MustMarshal(&types.PendingInit{}))
}

func (k *Keeper) addConsumer(ctx sdk.Context, chainId string, spawnTime time.Time, unbondingPeriod time.Duration) (string, error) {
	revision := ibcclienttypes.ParseChainID(chainId)
	msg := &ccvprovidertypes.MsgCreateConsumer{
		Submitter: authtypes.NewModuleAddress(types.ModuleName).String(),
		ChainId:   chainId,
		InitializationParameters: &ccvprovidertypes.ConsumerInitializationParameters{
			InitialHeight:                     ibcclienttypes.NewHeight(revision, 1),
			SpawnTime:                         spawnTime,
			UnbondingPeriod:                   unbondingPeriod,
			CcvTimeoutPeriod:                  ccvtypes.DefaultCCVTimeoutPeriod,
			TransferTimeoutPeriod:             ccvtypes.DefaultTransferTimeoutPeriod,
			ConsumerRedistributionFraction:    "0.0",
			BlocksPerDistributionTransmission: ccvtypes.DefaultBlocksPerDistributionTransmission,
			HistoricalEntries:                 ccvtypes.DefaultHistoricalEntries,
		},
		PowerShapingParameters: &ccvprovidertypes.PowerShapingParameters{
			Top_N:              0, // Start chainlets as opt-in chains
			ValidatorsPowerCap: 32,
		},
		AllowlistedRewardDenoms: &ccvprovidertypes.AllowlistedRewardDenoms{
			Denoms: []string{},
		},
		InfractionParameters: &ccvprovidertypes.InfractionParameters{
			DoubleSign: &ccvprovidertypes.SlashJailParameters{
				SlashFraction: math.LegacyNewDec(0), //TODO increase
				JailDuration:  time.Duration(1<<63 - 1),
				Tombstone:     true,
			},
			Downtime: &ccvprovidertypes.SlashJailParameters{
				SlashFraction: math.LegacyNewDec(0),
				JailDuration:  0, //TODO increase
				Tombstone:     false,
			},
		},
	}
	res, err := k.providerMsgServer.CreateConsumer(ctx, msg)
	if err != nil {
		return "", err
	}

	consumerID := res.ConsumerId

	// Enqueue an empty VSC packet
	valUpdateID := k.providerKeeper.GetValidatorSetUpdateId(ctx)
	packet := ccvtypes.NewValidatorSetChangePacketData(nil, valUpdateID, nil)
	k.providerKeeper.AppendPendingVSCPackets(ctx, consumerID, packet)
	k.providerKeeper.IncrementValidatorSetUpdateId(ctx)

	k.setPendingInit(ctx, consumerID)
	return consumerID, nil
}

// Forces sending queued VSC packets of new chainlets without waiting for the the provider epoch to end.
func (k *Keeper) InitConsumers(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletInit)

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		consumerID := string(iterator.Key())
		ctx.Logger().Debug(fmt.Sprintf("trying to initialize consumer %s", consumerID))

		// Check if the consumer is in the launched phase
		if k.providerKeeper.GetConsumerPhase(ctx, consumerID) != ccvprovidertypes.CONSUMER_PHASE_LAUNCHED {
			ctx.Logger().Debug(fmt.Sprintf("not initializing consumer %s: not launched phase", consumerID))
			continue
		}

		// Check if the channel is open
		channelId, found := k.providerKeeper.GetConsumerIdToChannelId(ctx, consumerID)
		if !found {
			ctx.Logger().Debug(fmt.Sprintf("not initializing consumer %s: channel not found", consumerID))
			continue
		}

		// Send the queued VSC packet immediately
		ctx.Logger().Info(fmt.Sprintf("force-sending queued VSC packets to consumer %s", consumerID))
		err := k.providerKeeper.SendVSCPacketsToChain(ctx, consumerID, channelId)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("force-sending VSC packets for %s failed: %s", consumerID, err))
			//NOTE: We can ignore the error as it only delays the packet
		}

		defer store.Delete(iterator.Key())
	}
}
