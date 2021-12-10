package account

import (
	"github.com/spf13/cobra"
)

var AccountRootCMD = &cobra.Command{
	Use:   "accounts",
	Short: "Set of commands for managing accounts",
	Long:  "Set of commands for managing accounts",
}

func init() {
	AccountRootCMD.AddCommand(importPrivKeyCmd)
	AccountRootCMD.AddCommand(generateKeyPairCmd)
	AccountRootCMD.AddCommand(transferBaseCurrencyCmd)
}
