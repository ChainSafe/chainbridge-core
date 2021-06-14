package chains

import (
	"encoding/json"
	"fmt"
	"os"

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

func NewConfig() *Config {
	return &Config{
		Chains: []RawChainConfig{},
	}
}

func GetConfig(path string, name string) (*Config, error) {
	var config Config

	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read in the config file, error: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config into struct, error: %w", err)
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) validate() error {
	for _, chain := range c.Chains {
		chainId := string(chain.Id)
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

func (c *Config) ToJSON(file string) *os.File {
	var (
		newFile *os.File
		err     error
	)

	var raw []byte
	if raw, err = json.Marshal(*c); err != nil {
		fmt.Println("error marshalling json", "err", err)
		os.Exit(1)
	}

	newFile, err = os.Create(file)
	if err != nil {
		fmt.Println("error creating config file", "err", err)
	}
	_, err = newFile.Write(raw)
	if err != nil {
		fmt.Println("error writing to config file", "err", err)

	}

	if err := newFile.Close(); err != nil {
		fmt.Println("failed to unmarshal config into struct", "err", err)
	}
	return newFile
}
