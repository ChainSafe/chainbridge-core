package config

import (
	"fmt"
	"math/big"
)

type SharedSubstrateConfig struct {
	GeneralChainConfig GeneralChainConfig
	StartBlock         *big.Int
	UseExtendedCall    bool
}

type RawSharedSubstrateConfig struct {
	GeneralChainConfig `mapstructure:",squash"`
	StartBlock         int64 `mapstructure:"startBlock"`
	UseExtendedCall    bool  `mapstructure:"useExtendedCall"`
}

func (c *RawSharedSubstrateConfig) ParseConfig() *SharedSubstrateConfig {

	c.GeneralChainConfig.ParseConfig()

	fmt.Println(c.GeneralChainConfig.KeystorePath)
	fmt.Println(c.GeneralChainConfig.BlockstorePath)
	fmt.Println(c.GeneralChainConfig.FreshStart)
	fmt.Println(c.GeneralChainConfig.LatestBlock)

	config := &SharedSubstrateConfig{
		GeneralChainConfig: c.GeneralChainConfig,
		StartBlock:         big.NewInt(c.StartBlock),
		UseExtendedCall:    c.UseExtendedCall,
	}
	return config
}
