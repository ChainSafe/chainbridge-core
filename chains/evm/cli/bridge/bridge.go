package bridge

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/spf13/cobra"
)

var BridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Set of commands for interacting with a bridge",
	Long:  "Set of commands for interacting with a bridge",
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
	BridgeCmd.AddCommand(cancelProposalCmd)
	BridgeCmd.AddCommand(queryProposalCmd)
	BridgeCmd.AddCommand(queryResourceCmd)
	BridgeCmd.AddCommand(registerGenericResourceCmd)
	BridgeCmd.AddCommand(registerResourceCmd)
	BridgeCmd.AddCommand(setBurnCmd)
	BridgeCmd.AddCommand(voteProposalCmd)
}
