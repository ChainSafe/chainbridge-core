package config

import (
	"fmt"

	"github.com/ChainSafe/chainbridge-core/config/relayer"
	"github.com/spf13/viper"
)

type Config struct {
	RelayerConfig relayer.RelayerConfig    `mapstructure:"relayer" json:"relayer"`
	ChainConfigs  []map[string]interface{} `mapstructure:"chains" json:"chains"`
}

func GetConfig(path string) (Config, error) {
	config := Config{}

	viper.SetConfigFile(path)
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	err = config.RelayerConfig.Validate()
	if err != nil {
		return config, err
	}
	for _, chain := range config.ChainConfigs {
		if chain["type"] == "" || chain["type"] == nil {
			return config, fmt.Errorf("Chain 'type' must be provided for every configured chain")
		}
	}

	return config, nil
}
