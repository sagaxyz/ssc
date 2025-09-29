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

func (k *Keeper) setPendingInit(ctx sdk.Context, chainId string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletInit)
	store.Set([]byte(chainId), k.cdc.MustMarshal(&types.PendingInit{}))
}

func (k *Keeper) addConsumer(ctx sdk.Context, chainId string, spawnTime time.Time) (string, error) {
	revision := ibcclienttypes.ParseChainID(chainId)
	msg := &ccvprovidertypes.MsgCreateConsumer{
		Submitter: authtypes.NewModuleAddress(types.ModuleName).String(),
		ChainId:   chainId,
		InitializationParameters: &ccvprovidertypes.ConsumerInitializationParameters{
			InitialHeight:                     ibcclienttypes.NewHeight(revision, 1),
			SpawnTime:                         spawnTime,
			UnbondingPeriod:                   ccvtypes.DefaultConsumerUnbondingPeriod, //TODO
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

	// Enqueue an empty VSC packet
	valUpdateID := k.providerKeeper.GetValidatorSetUpdateId(ctx)
	packet := ccvtypes.NewValidatorSetChangePacketData(nil, valUpdateID, nil)
	k.providerKeeper.AppendPendingVSCPackets(ctx, chainId, packet)
	k.providerKeeper.IncrementValidatorSetUpdateId(ctx)

	k.setPendingInit(ctx, chainId)
	return res.ConsumerId, nil
}

// Forces sending queued VSC packets of new chainlets without waiting for the the provider epoch to end.
func (k *Keeper) InitConsumers(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletInit)

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		chainId := string(iterator.Key())

		// Check if the consumer exists yet
		_, consumerRegistered := k.providerKeeper.GetConsumerClientId(ctx, chainId)
		if !consumerRegistered {
			continue
		}

		// Check if the channel is open
		channelId, found := k.providerKeeper.GetConsumerIdToChannelId(ctx, chainId)
		if !found {
			continue
		}

		// Send the queued VSC packet immediately
		ctx.Logger().Info(fmt.Sprintf("force-sending queued VSC packets to a new chainlet %s", chainId))
		err := k.providerKeeper.SendVSCPacketsToChain(ctx, chainId, channelId)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("force-sending VSC packets for %s failed: %s", chainId, err))
			//NOTE: We can ignore the error as it only delays the packet
		}

		defer store.Delete(iterator.Key())
	}
}
