package cli

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/admin"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/deploy"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/erc20"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/erc721"
	"github.com/spf13/cobra"
)

// BindCLI is public function to be invoked in example-app's cobra command
func BindCLI(cli *cobra.Command) {
	cli.AddCommand(evmRootCLI)
}

var evmRootCLI = &cobra.Command{
	Use:   "evm-cli",
	Short: "EVM CLI",
	Long:  "Root command for starting EVM CLI",
}

func init() {
	// persistent flags
	// to be used across all evm-cli commands (i.e. global)
	evmRootCLI.PersistentFlags().String("url", "ws://localhost:8545", "node url")
	evmRootCLI.PersistentFlags().Uint64("gasLimit", 6721975, "gasLimit used in transactions")
	evmRootCLI.PersistentFlags().Uint64("gasPrice", 20000000000, "gasPrice used for transactions")
	evmRootCLI.PersistentFlags().Uint64("networkID", 0, "networkid")
	evmRootCLI.PersistentFlags().String("privateKey", "ws://localhost:8545", "Private key to use")
	evmRootCLI.PersistentFlags().String("jsonWallet", "ws://localhost:8545", "Encrypted JSON wallet")
	evmRootCLI.PersistentFlags().String("jsonWalletPassword", "ws://localhost:8545", "Password for encrypted JSON wallet")

	// add commands to evm-cli root
	// deploy
	evmRootCLI.AddCommand(deploy.DeployEVM)

	// admin
	evmRootCLI.AddCommand(admin.AddAdminCmd)
	evmRootCLI.AddCommand(admin.AddRelayerCmd)
	evmRootCLI.AddCommand(admin.IsRelayerCmd)
	evmRootCLI.AddCommand(admin.PauseCmd)
	evmRootCLI.AddCommand(admin.RemoveAdminCmd)
	evmRootCLI.AddCommand(admin.RemoveRelayerCmd)
	evmRootCLI.AddCommand(admin.SetFeeCmd)
	evmRootCLI.AddCommand(admin.SetThresholdCmd)
	evmRootCLI.AddCommand(admin.UnpauseCmd)
	evmRootCLI.AddCommand(admin.WithdrawCmd)

	// bridge
	evmRootCLI.AddCommand(bridge.CancelProposalCmd)
	evmRootCLI.AddCommand(bridge.QueryProposalCmd)
	evmRootCLI.AddCommand(bridge.QueryResourceCmd)
	evmRootCLI.AddCommand(bridge.RegisterGenericResourceCmd)

	// erc20
	evmRootCLI.AddCommand(erc20.AddMinterCmd)

	// erc721
	evmRootCLI.AddCommand(erc721.AddMinterCmd)
}

/*
func EVMCLI(cli *cobra.Command) *cobra.Command {
	evmRootCLI.Flags().String("url", "ws://localhost:8545", "node url")
	evmRootCLI.Flags().Uint64("gasLimit", 6721975, "gasLimit used in transactions")
	evmRootCLI.Flags().Uint64("gasPrice", 20000000000, "gasPrice used for transactions")
	evmRootCLI.Flags().Uint64("networkID", 0, "networkid")
	evmRootCLI.Flags().String("privateKey", "ws://localhost:8545", "Private key to usel")
	evmRootCLI.Flags().String("jsonWallet", "ws://localhost:8545", "Encrypted JSON wallet")
	evmRootCLI.Flags().String("jsonWalletPassword", "ws://localhost:8545", "Password for encrypted JSON wallet")

	evmRootCLI.AddCommand(deploy.DeployEVM)
	evmRootCLI.AddCommand(bridge.CancelProposalEVM)

	cli.AddCommand(evmRootCLI)

	return evmRootCLI
}
*/
