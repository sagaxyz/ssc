package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sagaxyz/ssc/x/escrow/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetChainletAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-chainlet-account [chain-id]",
		Short: "Query getChainletAccount",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqChainId := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetChainletAccountRequest{
				ChainId: reqChainId,
			}

			res, err := queryClient.GetChainletAccount(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
