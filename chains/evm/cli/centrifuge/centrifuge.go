package centrifuge

import (
	"github.com/spf13/cobra"
)

var CentrifugeCmd = &cobra.Command{
	Use:   "centrifuge",
	Short: "Centrifuge related instructions",
	Long:  "Centrifuge related instructions",
}

func init() {
	CentrifugeCmd.AddCommand(deployCmd)
	CentrifugeCmd.AddCommand(getHashCmd)
}
