package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/spf13/cobra"
)

func CmdCreateChainletStack() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-chainlet-stack [display-name] [description] [image] [version] [checksum] [epochfee] [upfrontfee] [ccv-consumer]",
		Short: "Broadcast message create-chainlet-stack",
		Args:  cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDisplayName := args[0]
			argDescription := args[1]
			argImage := args[2]
			argVersion := args[3]
			argChecksum := args[4]
			argEpochFee := args[5]
			argSetupFee := args[6]
			argCcvConsumer := args[7]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			ccvConsumer, err := strconv.ParseBool(argCcvConsumer)
			if err != nil {
				return err
			}

			epochDenom, err := sdk.ParseCoinNormalized(argEpochFee)
			if err != nil {
				return err
			}

			setupDenom, err := sdk.ParseCoinNormalized(argSetupFee)
			if err != nil {
				return err
			}

			if epochDenom.GetDenom() != setupDenom.GetDenom() {
				return types.ErrMismatchedDenom
			}

			msg := types.NewMsgCreateChainletStack(
				clientCtx.GetFromAddress().String(),
				argDisplayName,
				argDescription,
				argImage,
				argVersion,
				argChecksum,
				types.ChainletStackFees{
					Denom:    epochDenom.GetDenom(),
					EpochFee: argEpochFee,
					SetupFee: argSetupFee,
				},
				ccvConsumer,
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
