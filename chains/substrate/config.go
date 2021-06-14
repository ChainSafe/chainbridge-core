package substrate

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains"
	"github.com/spf13/viper"
)

type SubstrateConfig struct {
	GeneralChainConfig chains.GeneralChainConfig
	StartBlock         *big.Int
	UseExtendedCall    bool
}

type RawSubstrateConfig struct {
	GeneralChainConfig chains.GeneralChainConfig
	StartBlock         int64 `mapstructure:"startBlock"`
	UseExtendedCall    bool  `mapstructure:"useExtendedCall"`
}

func GetConfig(path string, name string) (*SubstrateConfig, error) {
	var config RawSubstrateConfig

	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read in the config file, error: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config into struct, error: %w", err)
	}

	if err := config.GeneralChainConfig.Validate(); err != nil {
		return nil, err
	}

	parsedCfg := parseConfig(&config)

	return parsedCfg, nil
}

func parseConfig(rawConfig *RawSubstrateConfig) *SubstrateConfig {
	config := &SubstrateConfig{
		GeneralChainConfig: rawConfig.GeneralChainConfig,
		StartBlock:         big.NewInt(rawConfig.StartBlock),
		UseExtendedCall:    rawConfig.UseExtendedCall,
	}
	return config
}
