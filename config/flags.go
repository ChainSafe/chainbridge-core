package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Flags for running the Chainbridge app
	ConfigFlagName      = "config"
	KeystoreFlagName    = "keystore"
	BlockstoreFlagName  = "blockstore"
	FreshStartFlagName  = "fresh"
	LatestBlockFlagName = "latest"
	TestKeyFlagName     = "testkey"

	// Flags for all EVM CLI commands
	UrlFlagName                = "url"
	GasLimitFlagName           = "gasLimit"
	GasPriceFlagName           = "gasPrice"
	NetworkIdFlagName          = "networkid"
	PrivateKeyFlagName         = "privateKey"
	JsonWalletFlagName         = "jsonWallet"
	JsonWalletPasswordFlagName = "jsonWalletPassword"

	// Flags for all EVM Deploy CLI commands
	BridgeFlagName           = "bridge"
	Erc20HandlerFlagName     = "erc20Handler"
	Erc20FlagName            = "erc20"
	Erc721FlagName           = "erc721"
	DeployAllFlagName        = "all"
	RelayerThresholdFlagName = "relayerThreshold"
	ChainIdFlagName          = "chainId"
	RelayersFlagName         = "relayers"
	FeeFlagName              = "fee"
	BridgeAddressFlagName    = "bridgeAddress"
	Erc20SymbolFlagName      = "erc20Symbol"
	Erc20NameFlagName        = "erc20Name"
)

func BindFlags(rootCMD *cobra.Command) {
	rootCMD.PersistentFlags().String(ConfigFlagName, ".", "Path to JSON configuration files directory")
	viper.BindPFlag(ConfigFlagName, rootCMD.PersistentFlags().Lookup(ConfigFlagName))

	rootCMD.PersistentFlags().String(BlockstoreFlagName, "./lvldbdata", "Specify path for blockstore")
	viper.BindPFlag(BlockstoreFlagName, rootCMD.PersistentFlags().Lookup(BlockstoreFlagName))

	rootCMD.PersistentFlags().Bool(FreshStartFlagName, false, "Disables loading from blockstore at start. Opts will still be used if specified. (default: false)")
	viper.BindPFlag(FreshStartFlagName, rootCMD.PersistentFlags().Lookup(FreshStartFlagName))

	rootCMD.PersistentFlags().Bool(LatestBlockFlagName, false, "Overrides blockstore and start block, starts from latest block (default: false)")
	viper.BindPFlag(LatestBlockFlagName, rootCMD.PersistentFlags().Lookup(LatestBlockFlagName))

	rootCMD.PersistentFlags().String(KeystoreFlagName, "./keys", "Path to keystore directory")
	viper.BindPFlag(KeystoreFlagName, rootCMD.PersistentFlags().Lookup(KeystoreFlagName))

	rootCMD.PersistentFlags().String(TestKeyFlagName, "", "Applies a predetermined test keystore to the chains.")
	viper.BindPFlag(TestKeyFlagName, rootCMD.PersistentFlags().Lookup(TestKeyFlagName))
}

func BindEVMCLIFlags(evmRootCLI *cobra.Command) {
	evmRootCLI.PersistentFlags().String(UrlFlagName, "ws://localhost:8545", "node url")
	evmRootCLI.PersistentFlags().Uint64(GasLimitFlagName, 6721975, "gasLimit used in transactions")
	evmRootCLI.PersistentFlags().Uint64(GasPriceFlagName, 20000000000, "gasPrice used for transactions")
	evmRootCLI.PersistentFlags().Uint64(NetworkIdFlagName, 0, "networkid")
	evmRootCLI.PersistentFlags().String(PrivateKeyFlagName, "", "Private key to use")
	evmRootCLI.PersistentFlags().String(JsonWalletFlagName, "", "Encrypted JSON wallet")
	evmRootCLI.PersistentFlags().String(JsonWalletPasswordFlagName, "", "Password for encrypted JSON wallet")

	viper.BindPFlag(UrlFlagName, evmRootCLI.PersistentFlags().Lookup(UrlFlagName))
	viper.BindPFlag(GasLimitFlagName, evmRootCLI.PersistentFlags().Lookup(GasLimitFlagName))
	viper.BindPFlag(GasPriceFlagName, evmRootCLI.PersistentFlags().Lookup(GasPriceFlagName))
	viper.BindPFlag(NetworkIdFlagName, evmRootCLI.PersistentFlags().Lookup(NetworkIdFlagName))
	viper.BindPFlag(PrivateKeyFlagName, evmRootCLI.PersistentFlags().Lookup(PrivateKeyFlagName))
	viper.BindPFlag(JsonWalletFlagName, evmRootCLI.PersistentFlags().Lookup(JsonWalletFlagName))
	viper.BindPFlag(JsonWalletPasswordFlagName, evmRootCLI.PersistentFlags().Lookup(JsonWalletPasswordFlagName))

}

func BindDeployEVMFlags(deployCmd *cobra.Command) {
	deployCmd.Flags().Bool(BridgeFlagName, false, "deploy bridge")
	deployCmd.Flags().Bool(Erc20HandlerFlagName, false, "deploy ERC20 handler")
	//deployCmd.Flags().Bool("erc721Handler", false, "deploy ERC721 handler")
	//deployCmd.Flags().Bool("genericHandler", false, "deploy generic handler")
	deployCmd.Flags().Bool(Erc20FlagName, false, "deploy ERC20")
	deployCmd.Flags().Bool(Erc721FlagName, false, "deploy ERC721")
	deployCmd.Flags().Bool(DeployAllFlagName, false, "deploy all")
	deployCmd.Flags().Uint64(RelayerThresholdFlagName, 1, "number of votes required for a proposal to pass")
	deployCmd.Flags().String(ChainIdFlagName, "1", "chain ID for the instance")
	deployCmd.Flags().StringSlice(RelayersFlagName, []string{}, "list of initial relayers")
	deployCmd.Flags().String(FeeFlagName, "0", "fee to be taken when making a deposit (in ETH, decimas are allowed)")
	deployCmd.Flags().String(BridgeAddressFlagName, "", "bridge contract address. Should be provided if handlers are deployed separately")
	deployCmd.Flags().String(Erc20SymbolFlagName, "", "ERC20 contract symbol")
	deployCmd.Flags().String(Erc20NameFlagName, "", "ERC20 contract name")

	viper.BindPFlag(BridgeFlagName, deployCmd.Flags().Lookup(BridgeFlagName))
	viper.BindPFlag(Erc20HandlerFlagName, deployCmd.Flags().Lookup(Erc20HandlerFlagName))
	viper.BindPFlag(Erc20FlagName, deployCmd.Flags().Lookup(Erc20FlagName))
	viper.BindPFlag(Erc721FlagName, deployCmd.Flags().Lookup(Erc721FlagName))
	viper.BindPFlag(DeployAllFlagName, deployCmd.Flags().Lookup(DeployAllFlagName))
	viper.BindPFlag(RelayerThresholdFlagName, deployCmd.Flags().Lookup(RelayerThresholdFlagName))
	viper.BindPFlag(ChainIdFlagName, deployCmd.Flags().Lookup(ChainIdFlagName))
	viper.BindPFlag(RelayersFlagName, deployCmd.Flags().Lookup(RelayersFlagName))
	viper.BindPFlag(FeeFlagName, deployCmd.Flags().Lookup(FeeFlagName))
	viper.BindPFlag(BridgeAddressFlagName, deployCmd.Flags().Lookup(BridgeAddressFlagName))
	viper.BindPFlag(Erc20SymbolFlagName, deployCmd.Flags().Lookup(Erc20SymbolFlagName))
	viper.BindPFlag(Erc20NameFlagName, deployCmd.Flags().Lookup(Erc20NameFlagName))
}
