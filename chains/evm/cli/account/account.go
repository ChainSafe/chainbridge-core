package account

import (
	"github.com/spf13/cobra"
)

var AccountRootCMD = &cobra.Command{
	Use:   "accounts",
	Short: "Account instructions",
	Long:  "Account instructions",
}

func init() {
	AccountRootCMD.AddCommand(importPrivKeyCmd)
	AccountRootCMD.AddCommand(generateKeyPairCmd)
	AccountRootCMD.AddCommand(transferBaseCurrencyCmd)
}
