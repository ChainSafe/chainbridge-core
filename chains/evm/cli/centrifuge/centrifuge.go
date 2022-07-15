package centrifuge

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/spf13/cobra"
)

var CentrifugeCmd = &cobra.Command{
	Use:   "centrifuge",
	Short: "Set of commands for interacting with a cetrifuge asset store contract",
	Long:  "Set of commands for interacting with a cetrifuge asset store contract",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		// fetch global flag values
		url, _, gasPrice, senderKeyPair, prepare, err = flags.GlobalFlagValues(cmd)
		if err != nil {
			return fmt.Errorf("could not get global flags: %v", err)
		}
		return nil
	},
}

func init() {
	CentrifugeCmd.AddCommand(deployCmd)
	CentrifugeCmd.AddCommand(getHashCmd)
}
