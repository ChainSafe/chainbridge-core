package bridge

import (
	"github.com/spf13/cobra"
)

var BridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Set of commands for operating with a bridge",
	Long:  "Set of commands for operating with a bridge",
}

func init() {
	BridgeCmd.AddCommand(cancelProposalCmd)
	BridgeCmd.AddCommand(queryProposalCmd)
	BridgeCmd.AddCommand(queryResourceCmd)
	BridgeCmd.AddCommand(registerGenericResourceCmd)
	BridgeCmd.AddCommand(registerResourceCmd)
	BridgeCmd.AddCommand(setBurnCmd)
}
