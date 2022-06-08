package relayer

import (
	"fmt"

	"github.com/rs/zerolog"
)

type RelayerConfig struct {
	OpenTelemetryCollectorURL string
	LogLevel                  zerolog.Level
	LogFile                   string
}

type RawRelayerConfig struct {
	OpenTelemetryCollectorURL string `mapstructure:"OpenTelemetryCollectorURL" json:"opentelemetryCollectorURL"`
	LogLevel                  string `mapstructure:"LogLevel" json:"logLevel" default:"info"`
	LogFile                   string `mapstructure:"LogFile" json:"logFile" default:"out.log"`
}

func (c *RawRelayerConfig) Validate() error {
	return nil
}

// NewRelayerConfig parses RawRelayerConfig into RelayerConfig
func NewRelayerConfig(rawConfig RawRelayerConfig) (RelayerConfig, error) {
	config := RelayerConfig{}
	err := rawConfig.Validate()
	if err != nil {
		return config, err
	}

	logLevel, err := zerolog.ParseLevel(rawConfig.LogLevel)
	if err != nil {
		return config, fmt.Errorf("unknown log level: %s", rawConfig.LogLevel)
	}
	config.LogLevel = logLevel

	config.LogFile = rawConfig.LogFile
	config.OpenTelemetryCollectorURL = rawConfig.OpenTelemetryCollectorURL

	return config, nil
}
