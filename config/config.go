package config

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/config/relayer"
	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

type Config struct {
	RelayerConfig relayer.RelayerConfig
	ChainConfigs  []map[string]interface{}
}

type RawConfig struct {
	RelayerConfig relayer.RawRelayerConfig `mapstructure:"relayer" json:"relayer"`
	ChainConfigs  []map[string]interface{} `mapstructure:"chains" json:"chains"`
}

// GetConfig reads config from file, validates it and parses
// it into config suitable for application
func GetConfig(path string) (Config, error) {
	rawConfig := RawConfig{}
	config := Config{}

	viper.SetConfigFile(path)
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&rawConfig)
	if err != nil {
		return config, err
	}

	if err := defaults.Set(&rawConfig); err != nil {
		return config, err
	}

	relayerConfig, err := relayer.NewRelayerConfig(rawConfig.RelayerConfig)
	if err != nil {
		return config, err
	}
	for _, chain := range rawConfig.ChainConfigs {
		if chain["type"] == "" || chain["type"] == nil {
			return config, fmt.Errorf("Chain 'type' must be provided for every configured chain")
		}
	}

	config.RelayerConfig = relayerConfig
	config.ChainConfigs = rawConfig.ChainConfigs

	return config, nil
}
