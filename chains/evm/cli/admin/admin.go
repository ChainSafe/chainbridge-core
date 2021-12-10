package admin

import (
	"github.com/spf13/cobra"
)

var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Set of commands for executing various admin actions",
	Long:  "Set of commands for executing various admin actions",
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
