package chains

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Chains []RawChainConfig `mapstructure:"chains"`
}

type RawChainConfig struct {
	GeneralChainConfig `mapstructure:",squash"`
	Opts               map[string]string
}

type GeneralChainConfig struct {
	Name     string `mapstructure:"name"`
	Type     string `mapstructure:"type"`
	Id       uint8  `mapstructure:"id"`
	Endpoint string `mapstructure:"endpoint"`
	From     string `mapstructure:"from"`
}

func GetConfig(path string, name string) (*Config, error) {
	var config Config

	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read in the config file, error: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config into struct, error: %w", err)
	}

	err = config.validate()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) validate() error {
	for _, chain := range c.Chains {
		chainId := strconv.Itoa(int(chain.Id))
		if chain.Type == "" {
			return fmt.Errorf("required field chain.Type empty for chain %v", chain.Id)
		}
		if chain.Endpoint == "" {
			return fmt.Errorf("required field chain.Endpoint empty for chain %v", chain.Id)
		}
		if chain.Name == "" {
			return fmt.Errorf("required field chain.Name empty for chain %v", chain.Id)
		}
		if chainId == "" {
			return fmt.Errorf("required field chain.Id empty for chain %v", chain.Id)
		}
		if chain.From == "" {
			return fmt.Errorf("required field chain.From empty for chain %v", chain.Id)
		}
	}
	return nil
}
