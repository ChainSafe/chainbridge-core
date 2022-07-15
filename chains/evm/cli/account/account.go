package account

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/spf13/cobra"
)

var AccountRootCMD = &cobra.Command{
	Use:   "accounts",
	Short: "Set of commands for managing accounts",
	Long:  "Set of commands for managing accounts",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		// fetch global flag values
		url, gasLimit, gasPrice, senderKeyPair, _, err = flags.GlobalFlagValues(cmd)
		if err != nil {
			return fmt.Errorf("could not get global flags: %v", err)
		}
		return nil
	},
}

func init() {
	AccountRootCMD.AddCommand(generateKeyPairCmd)
	AccountRootCMD.AddCommand(transferBaseCurrencyCmd)
}
