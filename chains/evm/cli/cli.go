package cli

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/account"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/admin"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/deploy"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/erc20"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/erc721"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

var (
	// Flags for all EVM CLI commands
	UrlFlagName                = "url"
	GasLimitFlagName           = "gasLimit"
	GasPriceFlagName           = "gasPrice"
	NetworkIdFlagName          = "networkid"
	PrivateKeyFlagName         = "privateKey"
	JsonWalletFlagName         = "jsonWallet"
	JsonWalletPasswordFlagName = "jsonWalletPassword"
)

func BindEVMCLIFlags(evmRootCLI *cobra.Command) {
	evmRootCLI.PersistentFlags().String(UrlFlagName, "ws://localhost:8545", "node url")
	evmRootCLI.PersistentFlags().Uint64(GasLimitFlagName, 6721975, "gasLimit used in transactions")
	evmRootCLI.PersistentFlags().Uint64(GasPriceFlagName, 0, "used as upperLimitGasPrice for transactions if not 0. Transactions gasPrice is defined by estimating it on network for pre London fork networks and by estimating BaseFee and MaxTipFeePerGas in post London networks")
	evmRootCLI.PersistentFlags().Uint64(NetworkIdFlagName, 0, "networkid")
	evmRootCLI.PersistentFlags().String(PrivateKeyFlagName, "", "Private key to use")
	evmRootCLI.PersistentFlags().String(JsonWalletFlagName, "", "Encrypted JSON wallet")
	evmRootCLI.PersistentFlags().String(JsonWalletPasswordFlagName, "", "Password for encrypted JSON wallet")

	_ = viper.BindPFlag(UrlFlagName, evmRootCLI.PersistentFlags().Lookup(UrlFlagName))
	_ = viper.BindPFlag(GasLimitFlagName, evmRootCLI.PersistentFlags().Lookup(GasLimitFlagName))
	_ = viper.BindPFlag(GasPriceFlagName, evmRootCLI.PersistentFlags().Lookup(GasPriceFlagName))
	_ = viper.BindPFlag(NetworkIdFlagName, evmRootCLI.PersistentFlags().Lookup(NetworkIdFlagName))
	_ = viper.BindPFlag(PrivateKeyFlagName, evmRootCLI.PersistentFlags().Lookup(PrivateKeyFlagName))
	_ = viper.BindPFlag(JsonWalletFlagName, evmRootCLI.PersistentFlags().Lookup(JsonWalletFlagName))
	_ = viper.BindPFlag(JsonWalletPasswordFlagName, evmRootCLI.PersistentFlags().Lookup(JsonWalletPasswordFlagName))

}

func init() {
	// persistent flags
	// to be used across all evm-cli commands (i.e. global)
	BindEVMCLIFlags(EvmRootCLI)

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

	EvmRootCLI.AddCommand(account.AccountRootCMD)

	// utils
	EvmRootCLI.AddCommand(utils.UtilsCmd)
}
