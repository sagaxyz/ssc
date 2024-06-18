package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/spf13/cobra"
)

func CmdDisableChainletStackVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable-chainlet-stack-version [display-name] [version]",
		Short: "Disable a specific chainlet stack version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDisplayName := args[0]
			argVersion := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDisableChainletStackVersion(
				clientCtx.GetFromAddress().String(),
				argDisplayName,
				argVersion,
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
