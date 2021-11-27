package chain

import (
	"math/big"

	"github.com/mitchellh/mapstructure"
)

type SubstrateConfig struct {
	GeneralChainConfig GeneralChainConfig
	StartBlock         *big.Int
	UseExtendedCall    bool
}

type RawSubstrateConfig struct {
	GeneralChainConfig `mapstructure:",squash"`
	StartBlock         int64 `mapstructure:"startBlock"`
	UseExtendedCall    bool  `mapstructure:"useExtendedCall"`
}

func NewSubstrateConfig(chainCOnfig map[string]interface{}) (*SubstrateConfig, error) {
	var c RawSubstrateConfig
	err := mapstructure.Decode(chainCOnfig, &c)
	if err != nil {
		return nil, err
	}

	err = c.Validate()
	if err != nil {
		return nil, err
	}

	c.GeneralChainConfig.ParseFlags()
	config := &SubstrateConfig{
		GeneralChainConfig: c.GeneralChainConfig,
		StartBlock:         big.NewInt(c.StartBlock),
		UseExtendedCall:    c.UseExtendedCall,
	}

	return config, nil
}
