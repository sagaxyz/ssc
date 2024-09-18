package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sagaxyz/ssc/x/escrow/types"
	"github.com/spf13/cobra"
)

func CmdBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance [address]",
		Short: "Query balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return fmt.Errorf("failed to get client query context: %v", err)
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryBalanceRequest{
				Address: reqAddress,
			}

			res, err := queryClient.Balance(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
