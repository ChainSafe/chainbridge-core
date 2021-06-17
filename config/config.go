package config

import (
	"fmt"
	"math/big"
)

const DefaultGasLimit = 6721975
const DefaultGasPrice = 20000000000
const DefaultGasMultiplier = 1
const DefaultBlockConfirmations = 10

type GeneralChainConfig struct {
	Name     string `mapstructure:"name"`
	Type     string `mapstructure:"type"`
	Id       *uint8 `mapstructure:"id"`
	Endpoint string `mapstructure:"endpoint"`
	From     string `mapstructure:"from"`
}

type SharedEVMConfig struct {
	GeneralChainConfig GeneralChainConfig
	Bridge             string
	Erc20Handler       string
	Erc721Handler      string
	GenericHandler     string
	MaxGasPrice        *big.Int
	GasMultiplier      *big.Float
	GasLimit           *big.Int
	Http               bool
	StartBlock         *big.Int
	BlockConfirmations *big.Int
}

type RawSharedEVMConfig struct {
	GeneralChainConfig `mapstructure:",squash"`
	Bridge             string  `mapstructure:"bridge"`
	Erc20Handler       string  `mapstructure:"erc20Handler"`
	Erc721Handler      string  `mapstructure:"erc721Handler"`
	GenericHandler     string  `mapstructure:"genericHandler"`
	MaxGasPrice        int64   `mapstructure:"maxGasPrice"`
	GasMultiplier      float64 `mapstructure:"gasMultiplier"`
	GasLimit           int64   `mapstructure:"gasLimit"`
	Http               bool    `mapstructure:"http"`
	StartBlock         int64   `mapstructure:"startBlock"`
	BlockConfirmations int64   `mapstructure:"blockConfirmations"`
}

func (c *GeneralChainConfig) Validate() error {
	// viper defaults to 0 for not specified ints, but we must have a valid chain id
	// Previous method of checking used a string cast like below
	//chainId := string(c.Id)
	if c.Id == nil {
		return fmt.Errorf("required field chain.Id empty for chain %v", c.Id)
	}
	if c.Type == "" {
		return fmt.Errorf("required field chain.Type empty for chain %v", *c.Id)
	}
	if c.Endpoint == "" {
		return fmt.Errorf("required field chain.Endpoint empty for chain %v", *c.Id)
	}
	if c.Name == "" {
		return fmt.Errorf("required field chain.Name empty for chain %v", *c.Id)
	}
	if c.From == "" {
		return fmt.Errorf("required field chain.From empty for chain %v", *c.Id)
	}
	return nil
}

func (c *RawSharedEVMConfig) Validate() error {
	if err := c.GeneralChainConfig.Validate(); err != nil {
		return err
	}
	if c.Bridge == "" {
		return fmt.Errorf("required field chain.Bridge empty for chain %v", *c.Id)
	}
	return nil
}

func (c *RawSharedEVMConfig) ParseConfig() (*SharedEVMConfig, error) {

	config := &SharedEVMConfig{
		GeneralChainConfig: c.GeneralChainConfig,
		Erc20Handler:       c.Erc20Handler,
		Erc721Handler:      c.Erc721Handler,
		GenericHandler:     c.GenericHandler,
		GasLimit:           big.NewInt(DefaultGasLimit),
		MaxGasPrice:        big.NewInt(DefaultGasPrice),
		GasMultiplier:      big.NewFloat(DefaultGasMultiplier),
		Http:               c.Http,
		StartBlock:         big.NewInt(c.StartBlock),
		BlockConfirmations: big.NewInt(DefaultBlockConfirmations),
	}

	if c.Bridge != "" {
		config.Bridge = c.Bridge
	} else {
		return nil, fmt.Errorf("must provide opts.bridge field for ethereum config")
	}

	if c.GasLimit != 0 {
		config.GasLimit = big.NewInt(c.GasLimit)
	}

	if c.MaxGasPrice != 0 {
		config.MaxGasPrice = big.NewInt(c.MaxGasPrice)
	}

	if c.GasMultiplier != 0 {
		config.GasMultiplier = big.NewFloat(c.GasMultiplier)
	}

	if c.BlockConfirmations != 0 {
		config.BlockConfirmations = big.NewInt(c.BlockConfirmations)
	}

	return config, nil
}
