package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ccvprovidertypes "github.com/cosmos/interchain-security/v5/x/ccv/provider/types"
	ccvtypes "github.com/cosmos/interchain-security/v5/x/ccv/types"

	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func (k *Keeper) setPendingInit(ctx sdk.Context, chainId string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainletInit)
	store.Set([]byte(chainId), k.cdc.MustMarshal(&types.PendingInit{}))
}

func (k *Keeper) addConsumer(ctx sdk.Context, chainId string, spawnTime time.Time) error {
	revision := ibcclienttypes.ParseChainID(chainId)
	err := k.providerKeeper.HandleConsumerAdditionProposal(ctx, &ccvprovidertypes.MsgConsumerAddition{
		ChainId:                           chainId,
		InitialHeight:                     ibcclienttypes.NewHeight(revision, 1),
		SpawnTime:                         spawnTime,
		UnbondingPeriod:                   ccvtypes.DefaultConsumerUnbondingPeriod,
		CcvTimeoutPeriod:                  ccvtypes.DefaultCCVTimeoutPeriod,
		TransferTimeoutPeriod:             ccvtypes.DefaultTransferTimeoutPeriod,
		ConsumerRedistributionFraction:    "0.0",
		BlocksPerDistributionTransmission: ccvtypes.DefaultBlocksPerDistributionTransmission,
		HistoricalEntries:                 ccvtypes.DefaultHistoricalEntries,
	})
	if err != nil {
		return err
	}

	// Enqueue an empty VSC packet
	valUpdateID := k.providerKeeper.GetValidatorSetUpdateId(ctx)
	packet := ccvtypes.NewValidatorSetChangePacketData(nil, valUpdateID, nil)
	k.providerKeeper.AppendPendingVSCPackets(ctx, chainId, packet)
	k.providerKeeper.IncrementValidatorSetUpdateId(ctx)

	k.setPendingInit(ctx, chainId)
	return nil
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
		channelId, found := k.providerKeeper.GetChainToChannel(ctx, chainId)
		if !found {
			continue
		}

		// Send the queued VSC packet immediately
		ctx.Logger().Info(fmt.Sprintf("force-sending queued VSC packets to a new chainlet %s", chainId))
		k.providerKeeper.SendVSCPacketsToChain(ctx, chainId, channelId)

		// ICA setup for upgrades
		err := k.InitICA(ctx, chainId)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("initializing ICA for chainlet %s failed: %s", chainId, err))
			continue
		}

		defer store.Delete(iterator.Key())
	}
}
