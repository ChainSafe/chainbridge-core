package bridge

import (
	"github.com/spf13/cobra"
)

var BridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Bridge-related instructions",
	Long:  "Bridge-related instructions",
}

func init() {
	BridgeCmd.AddCommand(cancelProposalCmd)
	BridgeCmd.AddCommand(queryProposalCmd)
	BridgeCmd.AddCommand(queryResourceCmd)
	BridgeCmd.AddCommand(registerGenericResourceCmd)
	BridgeCmd.AddCommand(registerResourceCmd)
	BridgeCmd.AddCommand(setBurnCmd)
}
