package cli

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/deploy"
	"github.com/spf13/cobra"
)

func BindCLI(cli *cobra.Command) {
	cli.AddCommand(evmRootCLI)
}

var evmRootCLI = &cobra.Command{
	Use: "evm-cli",
	Short: "root command for starting evm cli",
	Long: "root command for starting evm cli",
}

func EVMCLI() {
	evmRootCLI.Flags().String("url", "ws://localhost:8545", "node url")
	evmRootCLI.Flags().Uint64("gasLimit", 6721975, "gasLimit used in transactions")
	evmRootCLI.Flags().Uint64("gasPrice", 20000000000, "gasPrice used for transactions")
	evmRootCLI.Flags().Uint64("networkID", 0, "networkid")
	evmRootCLI.Flags().String("privateKey", "ws://localhost:8545", "Private key to usel")
	evmRootCLI.Flags().String("jsonWallet", "ws://localhost:8545", "Encrypted JSON wallet")
	evmRootCLI.Flags().String("jsonWalletPassword", "ws://localhost:8545", "Password for encrypted JSON wallet")


	evmRootCLI.AddCommand(deploy.DeployEVM)
}
