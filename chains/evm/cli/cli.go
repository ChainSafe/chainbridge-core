package cli

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/admin"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/deploy"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/erc20"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/erc721"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/spf13/cobra"
)

// BindCLI is public function to be invoked in example-app's cobra command
func BindCLI(cli *cobra.Command) {
	cli.AddCommand(EvmRootCLI)
}

var EvmRootCLI = &cobra.Command{
	Use:   "evm-cli",
	Short: "EVM CLI",
	Long:  "Root command for starting EVM CLI",
}

func init() {
	// persistent flags
	// to be used across all evm-cli commands (i.e. global)
	flags.BindEVMCLIFlags(EvmRootCLI)

	// add commands to evm-cli root
	// deploy
	EvmRootCLI.AddCommand(deploy.DeployEVM)

	// admin
	EvmRootCLI.AddCommand(admin.AdminCmd)

	// bridge
	EvmRootCLI.AddCommand(bridge.BridgeCmd)

	// erc20
	EvmRootCLI.AddCommand(erc20.ERC20Cmd)

	// erc721
	EvmRootCLI.AddCommand(erc721.ERC721Cmd)
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
