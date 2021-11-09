package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Flags for running the Chainbridge app
	ChainConfigFlagName   = "config"
	RelayerConfigFlagName = "relayer-config"
	KeystoreFlagName      = "keystore"
	BlockstoreFlagName    = "blockstore"
	FreshStartFlagName    = "fresh"
	LatestBlockFlagName   = "latest"
	TestKeyFlagName       = "testkey"
)

func BindFlags(rootCMD *cobra.Command) {
	rootCMD.PersistentFlags().String(ChainConfigFlagName, ".", "Path to JSON configuration files directory")
	_ = viper.BindPFlag(ChainConfigFlagName, rootCMD.PersistentFlags().Lookup(ChainConfigFlagName))

	rootCMD.PersistentFlags().String(RelayerConfigFlagName, "", "Path to relayer JSON file configuration")
	_ = viper.BindPFlag(RelayerConfigFlagName, rootCMD.PersistentFlags().Lookup(RelayerConfigFlagName))

	rootCMD.PersistentFlags().String(BlockstoreFlagName, "./lvldbdata", "Specify path for blockstore")
	_ = viper.BindPFlag(BlockstoreFlagName, rootCMD.PersistentFlags().Lookup(BlockstoreFlagName))

	rootCMD.PersistentFlags().Bool(FreshStartFlagName, false, "Disables loading from blockstore at start. Opts will still be used if specified. (default: false)")
	_ = viper.BindPFlag(FreshStartFlagName, rootCMD.PersistentFlags().Lookup(FreshStartFlagName))

	rootCMD.PersistentFlags().Bool(LatestBlockFlagName, false, "Overrides blockstore and start block, starts from latest block (default: false)")
	_ = viper.BindPFlag(LatestBlockFlagName, rootCMD.PersistentFlags().Lookup(LatestBlockFlagName))

	rootCMD.PersistentFlags().String(KeystoreFlagName, "./keys", "Path to keystore directory")
	_ = viper.BindPFlag(KeystoreFlagName, rootCMD.PersistentFlags().Lookup(KeystoreFlagName))

	rootCMD.PersistentFlags().String(TestKeyFlagName, "", "Applies a predetermined test keystore to the chains.")
	_ = viper.BindPFlag(TestKeyFlagName, rootCMD.PersistentFlags().Lookup(TestKeyFlagName))
}
