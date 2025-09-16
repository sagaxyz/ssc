package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/sagaxyz/ssc/x/escrow/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group escrow queries under a subcommand.
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// core queries
	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdGetChainletAccount())

	// pools & funders
	cmd.AddCommand(CmdGetPools())
	cmd.AddCommand(CmdGetFunders())
	cmd.AddCommand(CmdGetFunder())

	// reverse-index convenience
	cmd.AddCommand(CmdGetFunderBalance())

	// this line is used by starport scaffolding # 1
	return cmd
}

// ---------------------
// Params
// ---------------------
func CmdQueryParams() *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Show the current escrow module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			q := types.NewQueryClient(clientCtx)

			res, err := q.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
}

// ---------------------
// Chainlet head
// ---------------------
func CmdGetChainletAccount() *cobra.Command {
	return &cobra.Command{
		Use:   "chainlet [chain-id]",
		Short: "Get the chainlet head account record",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			chainID := args[0]
			clientCtx := client.GetClientContextFromCmd(cmd)
			q := types.NewQueryClient(clientCtx)

			res, err := q.GetChainletAccount(cmd.Context(), &types.QueryGetChainletAccountRequest{
				ChainId: chainID,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
}

// ---------------------
// Pools (per chainlet)
// ---------------------
func CmdGetPools() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pools [chain-id]",
		Short: "List all denom pools for a chainlet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			chainID := args[0]
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			q := types.NewQueryClient(clientCtx)
			res, err := q.GetPools(cmd.Context(), &types.QueryPoolsRequest{
				ChainId:    chainID,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddPaginationFlagsToCmd(cmd, "pools")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// ---------------------
// Funders (per {chainlet, denom})
// ---------------------
func CmdGetFunders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "funders [chain-id] [denom]",
		Short: "List all funders for a specific {chain-id, denom} pool",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			chainID, denom := args[0], args[1]
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			q := types.NewQueryClient(clientCtx)
			res, err := q.GetFunders(cmd.Context(), &types.QueryFundersRequest{
				ChainId:    chainID,
				Denom:      denom,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddPaginationFlagsToCmd(cmd, "funders")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// ---------------------
// Single funder shares
// ---------------------
func CmdGetFunder() *cobra.Command {
	return &cobra.Command{
		Use:   "funder [chain-id] [denom] [address]",
		Short: "Show a single funder's shares for {chain-id, denom, address}",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			chainID, denom, addr := args[0], args[1], args[2]
			clientCtx := client.GetClientContextFromCmd(cmd)

			q := types.NewQueryClient(clientCtx)
			res, err := q.GetFunder(cmd.Context(), &types.QueryFunderSharesRequest{
				ChainId: chainID,
				Denom:   denom,
				Address: addr,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
}

// ---------------------
// Funder positions (reverse index)
// ---------------------
func CmdGetFunderBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "funder-balance [address]",
		Short: "List all positions (shares) for a funder across chainlets/denoms",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addr := args[0]
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			q := types.NewQueryClient(clientCtx)
			res, err := q.GetFunderBalance(cmd.Context(), &types.QueryFunderPositionsRequest{
				Address:    addr,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddPaginationFlagsToCmd(cmd, "funder-balance")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
