package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ConfigFlagName      = "config"
	KeystoreFlagName    = "keystore"
	BlockstoreFlagName  = "blockstore"
	FreshStartFlagName  = "fresh"
	LatestBlockFlagName = "latest"
	TestKeyFlagName     = "testkey"
)

func BindFlags(rootCMD *cobra.Command) {
	rootCMD.Flags().String(ConfigFlagName, ".", "Path to JSON configuration files directory")
	viper.BindPFlag(ConfigFlagName, rootCMD.Flags().Lookup(ConfigFlagName))

	rootCMD.Flags().String(BlockstoreFlagName, "./lvldbdata", "Specify path for blockstore")
	viper.BindPFlag(BlockstoreFlagName, rootCMD.Flags().Lookup(BlockstoreFlagName))

	rootCMD.Flags().Bool(FreshStartFlagName, false, "Disables loading from blockstore at start. Opts will still be used if specified. (default: false)")
	viper.BindPFlag(FreshStartFlagName, rootCMD.Flags().Lookup(FreshStartFlagName))

	rootCMD.Flags().Bool(LatestBlockFlagName, false, "Overrides blockstore and start block, starts from latest block (default: false)")
	viper.BindPFlag(LatestBlockFlagName, rootCMD.Flags().Lookup(LatestBlockFlagName))

	rootCMD.Flags().String(KeystoreFlagName, "./keys", "Path to keystore directory")
	viper.BindPFlag(KeystoreFlagName, rootCMD.Flags().Lookup(KeystoreFlagName))

	rootCMD.Flags().String(TestKeyFlagName, "", "Applies a predetermined test keystore to the chains.")
	viper.BindPFlag(TestKeyFlagName, rootCMD.Flags().Lookup(TestKeyFlagName))
}
