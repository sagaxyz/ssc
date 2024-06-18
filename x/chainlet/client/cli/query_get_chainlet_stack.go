package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetChainletStack() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-chainlet-stack [display-name]",
		Short: "Query get-chainlet-stack",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqDisplayName := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetChainletStackRequest{

				DisplayName: reqDisplayName,
			}

			res, err := queryClient.GetChainletStack(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
