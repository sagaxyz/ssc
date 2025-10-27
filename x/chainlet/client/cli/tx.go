package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sagaxyz/ssc/x/chainlet/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdCreateChainletStack())
	cmd.AddCommand(CmdLaunchChainlet())
	cmd.AddCommand(CmdUpdateChainletStack())
	cmd.AddCommand(CmdUpgradeChainlet())
	cmd.AddCommand(CmdCancelChainletUpgrade())
	cmd.AddCommand(CmdDisableChainletStackVersion())
	cmd.AddCommand(CmdUpdateChainletStackFees())
	// this line is used by starport scaffolding # 1

	return cmd
}

// token format: "epoch[:setup]" e.g. "10usaga:100usaga" OR just "10usaga" (=> setup=epoch)
func CmdUpdateChainletStackFees() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-stack-fees [stack-name] [fees]",
		Short: "Update fees for a chainlet stack (order preserved). Each fee is epoch[:setup]; if setup omitted, it's copied from epoch.",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			stackName := args[0]

			var raw []string
			if len(args) == 2 {
				for _, t := range strings.Split(args[1], ",") {
					if s := strings.TrimSpace(t); s != "" {
						raw = append(raw, s)
					}
				}
			}
			flagFees, _ := cmd.Flags().GetStringSlice("stack-fee")
			for _, t := range flagFees {
				if s := strings.TrimSpace(t); s != "" {
					raw = append(raw, s)
				}
			}
			if len(raw) == 0 {
				return fmt.Errorf("no fees provided; pass CSV or --stack-fee epoch[:setup]")
			}

			build := func(tok string) (types.ChainletStackFees, error) {
				parts := strings.Split(tok, ":")
				epoch := strings.TrimSpace(parts[0])
				setup := epoch // default: copy epoch
				if len(parts) > 1 && strings.TrimSpace(parts[1]) != "" {
					setup = strings.TrimSpace(parts[1])
				}
				e, err := sdk.ParseCoinNormalized(epoch)
				if err != nil {
					return types.ChainletStackFees{}, fmt.Errorf("bad epoch fee %q: %w", epoch, err)
				}
				s, err := sdk.ParseCoinNormalized(setup)
				if err != nil {
					return types.ChainletStackFees{}, fmt.Errorf("bad setup fee %q: %w", setup, err)
				}
				if s.Denom != e.Denom {
					return types.ChainletStackFees{}, types.ErrInvalidDenom
				}
				return types.ChainletStackFees{
					Denom:    e.Denom,
					EpochFee: epoch,
					SetupFee: setup, // epoch copied if omitted
				}, nil
			}

			fees := make([]types.ChainletStackFees, 0, len(raw))
			for _, tok := range raw {
				f, err := build(tok)
				if err != nil {
					return err
				}
				fees = append(fees, f)
			}

			msg := &types.MsgUpdateChainletStackFees{
				Creator:           clientCtx.GetFromAddress().String(),
				ChainletStackName: stackName,
				Fees:              fees,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().StringSlice("stack-fee", nil, "Add fee option as epoch[:setup], e.g. 10usaga:100usaga. If setup omitted, it defaults to epoch.")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
