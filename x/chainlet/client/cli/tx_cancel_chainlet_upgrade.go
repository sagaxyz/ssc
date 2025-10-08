package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/spf13/cobra"
)

func CmdCancelChainletUpgrade() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-chainlet-upgrade <chain-id> <stack-version> <channel-id>",
		Short: "Broadcast message upgradeChainlet",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argChainId := args[0]
			argStackVersion := args[1]
			argChannelID := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelChainletUpgrade(
				clientCtx.GetFromAddress().String(),
				argChainId,
				argStackVersion,
				argChannelID,
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
