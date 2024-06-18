package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sagaxyz/ssc/x/billing/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetBillingHistory() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-billing-history [chain-id]",
		Short: "Query get-billing-history",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqChainId := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetBillingHistoryRequest{
				ChainId:    reqChainId,
				Pagination: pageReq,
			}

			res, err := queryClient.GetBillingHistory(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "get-billing-history")

	return cmd
}
