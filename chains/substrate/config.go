package substrate

import (
	"strconv"

	"github.com/ChainSafe/chainbridge-core/chains"
)

type SubstrateConfig struct {
	ChainConfig     chains.GeneralChainConfig
	StartBlock      uint64
	UseExtendedCall bool
}

func ParseConfig(rawConfig *chains.RawChainConfig) *SubstrateConfig {
	config := &SubstrateConfig{
		ChainConfig:     rawConfig.GeneralChainConfig,
		StartBlock:      0,
		UseExtendedCall: false,
	}

	if blk, ok := rawConfig.Opts["startBlock"]; ok {
		res, err := strconv.ParseUint(blk, 10, 32)
		if err != nil {
			panic(err)
		}
		config.StartBlock = res
	}

	if b, ok := rawConfig.Opts["useExtendedCall"]; ok {
		res, err := strconv.ParseBool(b)
		if err != nil {
			panic(err)
		}
		config.UseExtendedCall = res
	}

	return config
}
