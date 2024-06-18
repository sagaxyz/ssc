package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/spf13/cobra"
)

func CmdCreateChainletStack() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-chainlet-stack [display-name] [description] [image] [version] [checksum] [epochfee] [epochlength] [upfrontfee]",
		Short: "Broadcast message create-chainlet-stack",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDisplayName := args[0]
			argDescription := args[1]
			argImage := args[2]
			argVersion := args[3]
			argChecksum := args[4]
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateChainletStack(
				clientCtx.GetFromAddress().String(),
				argDisplayName,
				argDescription,
				argImage,
				argVersion,
				argChecksum,
				types.ChainletStackFees{
					EpochFee:    args[5],
					EpochLength: args[6],
					SetupFee:    args[7],
				},
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
