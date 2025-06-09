package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/spf13/cobra"
)

func CmdUpdateChainletStack() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-chainlet-stack [display-name] [image] [version] [checksum] [ccv-consumer]",
		Short: "Broadcast message update-chainlet-stack",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDisplayName := args[0]
			argImage := args[1]
			argVersion := args[2]
			argChecksum := args[3]
			argCcvConsumer := args[4]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			ccvConsumer, err := strconv.ParseBool(argCcvConsumer)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateChainletStack(
				clientCtx.GetFromAddress().String(),
				argDisplayName,
				argImage,
				argVersion,
				argChecksum,
				ccvConsumer,
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
