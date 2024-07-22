package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/sagaxyz/ssc/x/chainlet/types"
	"github.com/spf13/cobra"
)

func CmdLaunchChainlet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "launch-chainlet [maintainers] [stack-name] [stack-version] [name] [denom] [params]",
		Short: "Broadcast message launch-chainlet",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argMaintainers := args[0] // looks like 'address1,address2,address3...'
			argStackName := args[1]
			argStackVersion := args[2]
			argName := args[3]
			argDenom := args[4]
			argParams := args[5] // looks like '{"bondDemon":"asaga","denom":"asaga",...}'

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var params types.ChainletParams
			err = json.Unmarshal([]byte(argParams), &params)
			if err != nil {
				return err
			}

			maintainers := strings.Split(argMaintainers, ",")

			evmChainId, _ := cmd.Flags().GetInt64("evm-chain-id")
			networkVersion, _ := cmd.Flags().GetInt64("network-version")
			if evmChainId < 1 {
				return fmt.Errorf("invalid evm chain id %d", evmChainId)
			}
			if networkVersion < 1 {
				return fmt.Errorf("invalid network version %d", networkVersion)
			}
			argChainId := generateChainId(argName, evmChainId, networkVersion)
			msg := types.NewMsgLaunchChainlet(
				clientCtx.GetFromAddress().String(),
				maintainers,
				argStackName,
				argStackVersion,
				argName,
				argChainId,
				argDenom,
				params,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	now := time.Now()
	cmd.Flags().Int64("evm-chain-id", now.UTC().UnixMicro(), "evm chain id")
	cmd.Flags().Int64("network-version", 1, "network version")

	return cmd
}

func generateChainId(name string, evm, version int64) string {
	return fmt.Sprintf("%s_%d-%d", name, evm, version)
}
