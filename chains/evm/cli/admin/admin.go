package admin

import (
	"github.com/spf13/cobra"
)

var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin-related instructions",
	Long:  "Admin-related instructions",
}

func init() {
	AdminCmd.AddCommand(addAdminCmd)
	AdminCmd.AddCommand(addRelayerCmd)
	AdminCmd.AddCommand(isRelayerCmd)
	AdminCmd.AddCommand(pauseCmd)
	AdminCmd.AddCommand(removeAdminCmd)
	AdminCmd.AddCommand(removeRelayerCmd)
	AdminCmd.AddCommand(setFeeCmd)
	AdminCmd.AddCommand(setThresholdCmd)
	AdminCmd.AddCommand(getThresholdCmd)
	AdminCmd.AddCommand(unpauseCmd)
	AdminCmd.AddCommand(withdrawCmd)
	AdminCmd.AddCommand(setDepositNonceCmd)
}
