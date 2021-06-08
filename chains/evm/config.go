package evm

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains"
)

const DefaultGasLimit = 6721975
const DefaultGasPrice = 20000000000
const DefaultGasMultiplier = 1

// Chain specific options
var (
	BridgeOpt         = "bridge"
	Erc20HandlerOpt   = "erc20Handler"
	Erc721HandlerOpt  = "erc721Handler"
	GenericHandlerOpt = "genericHandler"
	MaxGasPriceOpt    = "maxGasPrice"
	GasLimitOpt       = "gasLimit"
	GasMultiplier     = "gasMultiplier"
	HttpOpt           = "http"
	// StartBlockOpt         = "startBlock"
	// BlockConfirmationsOpt = "blockConfirmations"
	// EGSApiKey             = "egsApiKey"
	// EGSSpeed              = "egsSpeed"
)

type EVMConfig struct {
	GeneralChainConfig chains.GeneralChainConfig
	Bridge             string
	Erc20Handler       string
	Erc721Handler      string
	GenericHandler     string
	MaxGasPrice        *big.Int
	GasMultiplier      *big.Float
	GasLimit           *big.Int
	Http               bool
}

func ParseConfig(rawConfig *chains.RawChainConfig) (*EVMConfig, error) {
	config := &EVMConfig{
		GeneralChainConfig: rawConfig.GeneralChainConfig,
		Bridge:             "",
		Erc20Handler:       "",
		Erc721Handler:      "",
		GenericHandler:     "",
		GasLimit:           big.NewInt(DefaultGasLimit),
		MaxGasPrice:        big.NewInt(DefaultGasPrice),
		GasMultiplier:      big.NewFloat(DefaultGasMultiplier),
		Http:               false,
	}

	if contract, ok := rawConfig.Opts[BridgeOpt]; ok && contract != "" {
		config.Bridge = contract
		delete(rawConfig.Opts, BridgeOpt)
	} else {
		return nil, fmt.Errorf("must provide opts.bridge field for ethereum config")
	}

	if contract, ok := rawConfig.Opts[Erc20HandlerOpt]; ok {
		config.Erc20Handler = contract
		delete(rawConfig.Opts, Erc20HandlerOpt)
	}

	if contract, ok := rawConfig.Opts[Erc721HandlerOpt]; ok {
		config.Erc721Handler = contract
		delete(rawConfig.Opts, Erc721HandlerOpt)
	}

	if contract, ok := rawConfig.Opts[GenericHandlerOpt]; ok {
		config.GenericHandler = contract
		delete(rawConfig.Opts, GenericHandlerOpt)
	}

	if gasPrice, ok := rawConfig.Opts[MaxGasPriceOpt]; ok {
		price := big.NewInt(0)
		_, pass := price.SetString(gasPrice, 10)
		if pass {
			config.MaxGasPrice = price
			delete(rawConfig.Opts, MaxGasPriceOpt)
		} else {
			return nil, errors.New("unable to parse max gas price")
		}
	}

	if gasLimit, ok := rawConfig.Opts[GasLimitOpt]; ok {
		limit := big.NewInt(0)
		_, pass := limit.SetString(gasLimit, 10)
		if pass {
			config.GasLimit = limit
			delete(rawConfig.Opts, GasLimitOpt)
		} else {
			return nil, errors.New("unable to parse gas limit")
		}
	}

	if gasMultiplier, ok := rawConfig.Opts[GasMultiplier]; ok {
		multiplier := big.NewFloat(1)
		_, pass := multiplier.SetString(gasMultiplier)
		if pass {
			config.GasMultiplier = multiplier
			delete(rawConfig.Opts, GasMultiplier)
		} else {
			return nil, errors.New("unable to parse gasMultiplier to float")
		}
	}

	if HTTP, ok := rawConfig.Opts[HttpOpt]; ok && HTTP == "true" {
		config.Http = true
		delete(rawConfig.Opts, HttpOpt)
	} else if HTTP, ok := rawConfig.Opts[HttpOpt]; ok && HTTP == "false" {
		config.Http = false
		delete(rawConfig.Opts, HttpOpt)
	}

	if len(rawConfig.Opts) != 0 {
		return nil, fmt.Errorf("unknown Opts Encountered: %v", rawConfig.Opts)
	}

	return config, nil
}
