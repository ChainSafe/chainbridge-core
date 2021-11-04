package relayer

import (
	"fmt"

	"github.com/spf13/viper"
)

type RelayerConfig struct {
	PrometheusPath            string `mapstructure:"PrometheusPath"`
	PrometheusPort            uint64 `mapstructure:"PrometheusPort"`
	OpenTelemetryCollectorURL string `mapstructure:"OpenTelemetryCollectorURL"`
}

func (c *RelayerConfig) Validate() error {
	if c.PrometheusPort < 1 || c.PrometheusPort > 65535 {
		return fmt.Errorf(`PrometheusPort outside of valid range of 1 - 65535`)
	}

	return nil
}

func setDefaultValues() {
	viper.SetDefault("PrometheusPath", "/metrics")
	viper.SetDefault("PrometheusPort", 2112)
}

func GetRelayerConfig(path string) (RelayerConfig, error) {
	config := RelayerConfig{}
	setDefaultValues()

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
