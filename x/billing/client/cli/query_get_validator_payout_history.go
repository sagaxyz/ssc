package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sagaxyz/ssc/x/billing/types"
	"github.com/spf13/cobra"
)

func CmdGetValidatorPayoutHistory() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-validator-payout-history [validator-address]",
		Short: "Query get-validator-payout-history",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqValidatorAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetValidatorPayoutHistoryRequest{
				ValidatorAddress: reqValidatorAddress,
				Pagination:       pageReq,
			}

			res, err := queryClient.GetValidatorPayoutHistory(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "get-validator-payout-history")

	return cmd
}
