package relayer

import (
	"github.com/spf13/viper"
)

type RelayerConfig struct {
	OpenTelemetryCollectorURL string `mapstructure:"OpenTelemetryCollectorURL"`
}

func (c *RelayerConfig) Validate() error {
	return nil
}

func GetRelayerConfig(path string) (RelayerConfig, error) {
	config := RelayerConfig{}

	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	err = config.Validate()
	if err != nil {
		return config, err
	}

	return config, nil
}
