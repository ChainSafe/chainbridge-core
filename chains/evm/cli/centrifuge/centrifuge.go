package centrifuge

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/spf13/cobra"
)

var CentrifugeCmd = &cobra.Command{
	Use:   "centrifuge",
	Short: "Centrifuge related instructions",
	Long:  "Centrifuge related instructions",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		// fetch global flag values
		url, gasLimit, gasPrice, senderKeyPair, err = flags.GlobalFlagValues(cmd)
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
