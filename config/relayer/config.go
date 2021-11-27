package relayer

type RelayerConfig struct {
	OpenTelemetryCollectorURL string `mapstructure:"OpenTelemetryCollectorURL"`
}

func (c *RelayerConfig) Validate() error {
	return nil
}
