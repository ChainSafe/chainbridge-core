package chain

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/mitchellh/mapstructure"
)

type EVMConfig struct {
	GeneralChainConfig GeneralChainConfig
	Bridge             string
	Erc20Handler       string
	Erc721Handler      string
	GenericHandler     string
	MaxGasPrice        *big.Int
	GasMultiplier      *big.Float
	GasLimit           *big.Int
	StartBlock         *big.Int
	BlockConfirmations uint8
}

type RawEVMConfig struct {
	GeneralChainConfig `mapstructure:",squash"`
	Bridge             string  `mapstructure:"bridge"`
	Erc20Handler       string  `mapstructure:"erc20Handler"`
	Erc721Handler      string  `mapstructure:"erc721Handler"`
	GenericHandler     string  `mapstructure:"genericHandler"`
	MaxGasPrice        int64   `mapstructure:"maxGasPrice"`
	GasMultiplier      float64 `mapstructure:"gasMultiplier"`
	GasLimit           int64   `mapstructure:"gasLimit"`
	StartBlock         int64   `mapstructure:"startBlock"`
	BlockConfirmations uint8   `mapstructure:"blockConfirmations"`
}

func (c *RawEVMConfig) Validate() error {
	if err := c.GeneralChainConfig.Validate(); err != nil {
		return err
	}
	if c.Bridge == "" {
		return fmt.Errorf("required field chain.Bridge empty for chain %v", *c.Id)
	}
	return nil
}

func NewEVMConfig(chainConfig map[string]interface{}) (*EVMConfig, error) {
	var c RawEVMConfig
	err := mapstructure.Decode(chainConfig, &c)
	if err != nil {
		return nil, err
	}

	err = c.Validate()
	if err != nil {
		return nil, err
	}

	c.GeneralChainConfig.ParseFlags()
	config := &EVMConfig{
		GeneralChainConfig: c.GeneralChainConfig,
		Erc20Handler:       c.Erc20Handler,
		Erc721Handler:      c.Erc721Handler,
		GenericHandler:     c.GenericHandler,
		Bridge:             c.Bridge,
		GasLimit:           big.NewInt(consts.DefaultGasLimit),
		MaxGasPrice:        big.NewInt(consts.DefaultGasPrice),
		GasMultiplier:      big.NewFloat(consts.DefaultGasMultiplier),
		StartBlock:         big.NewInt(c.StartBlock),
		BlockConfirmations: consts.DefaultBlockConfirmations,
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
		config.BlockConfirmations = c.BlockConfirmations
	}

	return config, nil
}
