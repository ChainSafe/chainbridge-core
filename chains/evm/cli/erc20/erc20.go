package erc20

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/spf13/cobra"
)

var ERC20Cmd = &cobra.Command{
	Use:   "erc20",
	Short: "Set of commands for interacting with an ERC20 contract",
	Long:  "Set of commands for interacting with an ERC20 contract",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		// fetch global flag values
		url, gasLimit, gasPrice, senderKeyPair, prepare, err = flags.GlobalFlagValues(cmd)
		if err != nil {
			return fmt.Errorf("could not get global flags: %v", err)
		}
		return nil
	},
}

func init() {
	ERC20Cmd.AddCommand(addMinterCmd)
	ERC20Cmd.AddCommand(getAllowanceCmd)
	ERC20Cmd.AddCommand(approveCmd)
	ERC20Cmd.AddCommand(balanceCmd)
	ERC20Cmd.AddCommand(depositCmd)
	ERC20Cmd.AddCommand(mintCmd)
}
