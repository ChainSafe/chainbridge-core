package substrate

import (
	"fmt"

	"github.com/ChainSafe/chainbridge-core/chains"
	"github.com/spf13/viper"
)

// This second struct exists for parsing fields into types that are not primitive
// Currently only used to place the generic chain config options under ChainConfig field
type SubstrateConfig struct {
	ChainConfig     chains.ChainConfig
	From            string
	StartBlock      uint32
	UseExtendedCall bool
}

type RawSubstrateConfig struct {
	chains.ChainConfig `mapstructure:",squash"`
	From               string `mapstructure:"from"`
	StartBlock         uint32 `mapstructure:"startBlock"`
	UseExtendedCall    bool   `mapstructure:"useExtendedCall"`
}

func GetConfig(path string, name string) (*SubstrateConfig, error) {
	var rawConfig RawSubstrateConfig

	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read in the config file, error: %w", err)
	}

	err = viper.Unmarshal(&rawConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config into struct, error: %w", err)
	}

	config := parseConfig(&rawConfig)

	return config, nil
}

func parseConfig(rawConfig *RawSubstrateConfig) *SubstrateConfig {
	return &SubstrateConfig{
		ChainConfig:     rawConfig.ChainConfig,
		From:            rawConfig.From,
		StartBlock:      rawConfig.StartBlock,
		UseExtendedCall: rawConfig.UseExtendedCall,
	}
}
