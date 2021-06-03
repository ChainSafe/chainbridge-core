package evm

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains"
	"github.com/spf13/viper"
)

const DefaultGasLimit = 6721975
const DefaultGasPrice = 20000000000
const DefaultGasMultiplier = 1

type EVMConfig struct {
	ChainConfig    chains.ChainConfig
	Bridge         string
	Erc20Handler   string
	Erc721Handler  string
	GenericHandler string
	MaxGasPrice    *big.Int
	GasMultiplier  *big.Float
	GasLimit       *big.Int
}

type RawEVMConfig struct {
	chains.ChainConfig `mapstructure:",squash"`
	Bridge             string  `mapstructure:"bridge"`
	Erc20Handler       string  `mapstructure:"erc20Handler"`
	Erc721Handler      string  `mapstructure:"erc721Handler"`
	GenericHandler     string  `mapstructure:"genericHandler"`
	MaxGasPrice        int64   `mapstructure:"maxGasPrice"`
	GasMultiplier      float64 `mapstructure:"gasMultiplier"`
	GasLimit           int64   `mapstructure:"gasLimit"`
}

func GetConfig(path string, name string) (*EVMConfig, error) {
	var rawConfig RawEVMConfig

	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read in the config file, error: %w", err)
	}

	// Set values requiring defaults
	rawConfig.MaxGasPrice = DefaultGasPrice
	rawConfig.GasMultiplier = DefaultGasMultiplier
	rawConfig.GasLimit = DefaultGasLimit

	err = viper.Unmarshal(&rawConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config into struct, error: %w", err)
	}

	config := parseConfig(&rawConfig)

	return config, nil
}

func parseConfig(rawConfig *RawEVMConfig) *EVMConfig {
	return &EVMConfig{
		ChainConfig:    rawConfig.ChainConfig,
		Bridge:         rawConfig.Bridge,
		Erc20Handler:   rawConfig.Erc20Handler,
		Erc721Handler:  rawConfig.Erc721Handler,
		GenericHandler: rawConfig.GenericHandler,
		GasLimit:       big.NewInt(rawConfig.GasLimit),
		MaxGasPrice:    big.NewInt(int64(rawConfig.MaxGasPrice)),
		GasMultiplier:  big.NewFloat(rawConfig.GasMultiplier),
	}
}
