package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/sagaxyz/ssc/x/peers/types"
)

func CmdSetPeers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-peers [chain-id] [peer1...peer2]",
		Short: "Set validator peers for a chain ID",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			chainId := args[0]
			peers := args[1:]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetPeers(
				clientCtx.GetFromAddress().String(),
				chainId,
				peers...,
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
