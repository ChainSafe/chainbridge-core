package erc20

import (
	"github.com/spf13/cobra"
)

var ERC20Cmd = &cobra.Command{
	Use:   "erc20",
	Short: "ERC20-related instructions",
	Long:  "ERC20-related instructions",
}

func init() {
	ERC20Cmd.AddCommand(addMinterCmd)
	ERC20Cmd.AddCommand(allowanceCmd)
	ERC20Cmd.AddCommand(approveCmd)
	ERC20Cmd.AddCommand(balanceCmd)
	ERC20Cmd.AddCommand(depositCmd)
	ERC20Cmd.AddCommand(mintCmd)
}
