package chain

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

const DefaultGasLimit = 6721975
const DefaultGasPrice = 15000000
const DefaultGasMultiplier = 1
const DefaultBlockConfirmations = 10

type OptimismConfig struct {
	EVMConfig        EVMConfig // Contains configuration of Optimism l2geth client which is used to for sequencing transactions to Optimism
	VerifyRollup     bool
	VerifierEndpoint string // This is the endpoint for an Optimism replica which is read-only and used only for verifying transactions
}

type RawOptimismConfig struct {
	RawEVMConfig     `mapstructure:",squash"`
	VerifyRollup     bool   `mapstructure:"verifyRollup"`
	VerifierEndpoint string `mapstructure:"verifierEndpoint"`
}

func (c *RawOptimismConfig) Validate() error {
	if err := c.RawEVMConfig.Validate(); err != nil {
		return err
	}
	return nil
}

// NewOptimismConfig decodes and validates an instance of an OptimismConfig from
// raw chain config
func NewOptimismConfig(chainConfig map[string]interface{}) (*OptimismConfig, error) {
	log.Debug().Msg("got into optimism config")
	var c RawOptimismConfig
	err := mapstructure.Decode(chainConfig, &c)
	if err != nil {
		return nil, err
	}
	log.Debug().Msg("successfully decoded")
	err = c.Validate()
	if err != nil {
		return nil, err
	}

	evmCfg, err := c.RawEVMConfig.ParseConfig()
	if err != nil {
		return nil, err
	}

	config := &OptimismConfig{
		EVMConfig:        *evmCfg,
		VerifyRollup:     c.VerifyRollup,
		VerifierEndpoint: c.VerifierEndpoint,
	}

	return config, nil
}
