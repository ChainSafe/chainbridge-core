package centrifuge

import (
	"github.com/spf13/cobra"
)

var CentrifugeCmd = &cobra.Command{
	Use:   "centrifuge",
	Short: "Set of commands for interacting with a cetrifuge asset store contract",
	Long:  "Set of commands for interacting with a cetrifuge asset store contract",
}

func init() {
	CentrifugeCmd.AddCommand(deployCmd)
	CentrifugeCmd.AddCommand(getHashCmd)
}
