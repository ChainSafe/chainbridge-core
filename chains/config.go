package chains

import (
	"fmt"
)

type GeneralChainConfig struct {
	Name     string `mapstructure:"name"`
	Type     string `mapstructure:"type"`
	Id       *uint8 `mapstructure:"id"`
	Endpoint string `mapstructure:"endpoint"`
	From     string `mapstructure:"from"`
}

type SharedEVMConfig struct {
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

func (c *SharedEVMConfig) Validate() error {
	if err := c.GeneralChainConfig.Validate(); err != nil {
		return err
	}
	if c.Bridge == "" {
		return fmt.Errorf("required field chain.Bridge empty for chain %v", *c.Id)
	}
	return nil
}
