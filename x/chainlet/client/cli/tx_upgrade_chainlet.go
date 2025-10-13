package cli

import (
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/spf13/cobra"
)

func CmdUpgradeChainlet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade-chainlet <chain-id> <stack-version> [height-delta] [channel-id] [unbonding-period]",
		Short: "Broadcast message upgradeChainlet",
		Args:  cobra.RangeArgs(2, 5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argChainId := args[0]
			argStackVersion := args[1]
			var heightDelta uint64
			if len(args) > 2 {
				argHeightDelta := args[2]
				heightDelta, err = strconv.ParseUint(argHeightDelta, 10, 64)
				if err != nil {
					return err
				}
			}
			var argChannelID string
			if len(args) > 3 {
				argChannelID = args[3]
			}
			var unbondingPeriod *time.Duration
			if len(args) > 4 {
				var d time.Duration
				d, err = time.ParseDuration(args[4])
				if err != nil {
					return
				}
				unbondingPeriod = &d
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpgradeChainlet(
				clientCtx.GetFromAddress().String(),
				argChainId,
				argStackVersion,
				heightDelta,
				argChannelID,
				unbondingPeriod,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
