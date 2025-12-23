package v3

import (
	"context"
	"fmt"
	"time"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	transferkeeper "github.com/cosmos/ibc-go/v10/modules/apps/transfer/keeper"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
)

const Name = "2-to-3"

func UpgradeHandler(mm *module.Manager, configurator module.Configurator, transferKeeper transferkeeper.Keeper, bankKeeper bankkeeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}

		sdkCtx := sdk.UnwrapSDKContext(ctx)

		// curated list of addresses with stuck balances
		addresses := []string{
			"saga1vhj6jw40shhms5ssxy7j4k8ejrvlwt9f3awg2c",
			"saga1pdegmtwnwr8j6fq65zlesudlqejy99dh3ava9s",
			"saga1x7pfl6dcueacye7zqk9egk04yju78jja3hcmx4",
			"saga19tzflvk3pwd6shh83df44c56rfdjc5czrzt4ra",
			"saga1ru9et4vqjrelpq5ze8y825v7ws08e5z9drevdd",
			"saga1kkay5yc0nccqxmgurkc34zgnete6eh462m2hu9",
		}

		firstChannel := "channel-36"
		lastChannel := "channel-4"

		if sdkCtx.ChainID() == "ssc-staging-1" {
			firstChannel = "channel-0"
			lastChannel = "channel-0"
		}

		if sdkCtx.ChainID() == "sscd-a" {
			firstChannel = "channel-0"
			lastChannel = "channel-1"
		}

		// get all the balances and forward them to the same addresses but on channel-4
		for _, address := range addresses {
			addr, err := sdk.AccAddressFromBech32(address)
			if err != nil {
				return nil, err
			}
			balances := bankKeeper.GetAllBalances(sdkCtx, addr)
			for _, balance := range balances {
				msg := &transfertypes.MsgTransfer{
					SourcePort:       "transfer",
					SourceChannel:    firstChannel,
					Token:            balance,
					Sender:           addr.String(),
					Receiver:         addr.String(),
					TimeoutTimestamp: uint64(sdkCtx.BlockTime().Add(time.Second * 600).UnixNano()), // allow 1 day for the transfer
					Memo:             fmt.Sprintf(`{"forward":{"port":"transfer","channel":"%s","receiver":"%s"}}`, lastChannel, address),
				}
				res, err := transferKeeper.Transfer(sdkCtx, msg)
				if err != nil {
					return nil, fmt.Errorf("failed to send transfer: %w, address: %s, coin: %s", err, address, balance.String())
				}

				sdkCtx.Logger().Info("transferred", "address", address, "coin", balance.String(), "sequence", res.Sequence)
			}
		}

		return newVM, nil
	}
}
