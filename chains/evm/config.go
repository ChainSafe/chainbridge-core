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
const DefaultBlockConfirmations = 10

// Chain specific options
// var (
// 	BridgeOpt             = "bridge"
// 	Erc20HandlerOpt       = "erc20Handler"
// 	Erc721HandlerOpt      = "erc721Handler"
// 	GenericHandlerOpt     = "genericHandler"
// 	MaxGasPriceOpt        = "maxGasPrice"
// 	GasLimitOpt           = "gasLimit"
// 	GasMultiplier         = "gasMultiplier"
// 	HttpOpt               = "http"
// 	StartBlockOpt         = "startBlock"
// 	BlockConfirmationsOpt = "blockConfirmations"
// 	EGSApiKey             = "egsApiKey"
// 	EGSSpeed              = "egsSpeed"
// )

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
	StartBlock         *big.Int
	BlockConfirmations *big.Int
	EgsApiKey          string // API key for ethgasstation to query gas prices
	EgsSpeed           string // The speed which a transaction should be processed: average, fast, fastest. Default: fast
}

type RawEVMConfig struct {
	chains.GeneralChainConfig `mapstructure:",squash"`
	Bridge                    string  `mapstructure:"bridge"`
	Erc20Handler              string  `mapstructure:"erc20Handler"`
	Erc721Handler             string  `mapstructure:"erc721Handler"`
	GenericHandler            string  `mapstructure:"genericHandler"`
	MaxGasPrice               int64   `mapstructure:"maxGasPrice"`
	GasMultiplier             float64 `mapstructure:"gasMultiplier"`
	GasLimit                  int64   `mapstructure:"gasLimit"`
	Http                      bool    `mapstructure:"http"`
	StartBlock                int64   `mapstructure:"startBlock"`
	BlockConfirmations        int64   `mapstructure:"blockConfirmations"`
	EgsApiKey                 string  `mapstructure:"egsApiKey"`
	EgsSpeed                  string  `mapstructure:"egsSpeed"`
}

func GetConfig(path string, name string) (*EVMConfig, error) {
	config := &RawEVMConfig{}

	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read in the config file, error: %w", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config into struct, error: %w", err)
	}

	if err := config.GeneralChainConfig.Validate(); err != nil {
		return nil, err
	}

	cfg, err := ParseConfig(config)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func ParseConfig(rawConfig *RawEVMConfig) (*EVMConfig, error) {

	config := &EVMConfig{
		GeneralChainConfig: rawConfig.GeneralChainConfig,
		Erc20Handler:       rawConfig.Erc20Handler,
		Erc721Handler:      rawConfig.Erc721Handler,
		GenericHandler:     rawConfig.GenericHandler,
		GasLimit:           big.NewInt(DefaultGasLimit),
		MaxGasPrice:        big.NewInt(DefaultGasPrice),
		GasMultiplier:      big.NewFloat(DefaultGasMultiplier),
		Http:               rawConfig.Http,
		StartBlock:         big.NewInt(rawConfig.StartBlock),
		BlockConfirmations: big.NewInt(DefaultBlockConfirmations),
		EgsApiKey:          "",
		EgsSpeed:           "",
	}

	if rawConfig.Bridge != "" {
		config.Bridge = rawConfig.Bridge
	} else {
		return nil, fmt.Errorf("must provide opts.bridge field for ethereum config")
	}

	if rawConfig.GasLimit != 0 {
		config.GasLimit = big.NewInt(rawConfig.GasLimit)
	}

	if rawConfig.MaxGasPrice != 0 {
		config.MaxGasPrice = big.NewInt(rawConfig.MaxGasPrice)
	}

	if rawConfig.GasMultiplier != 0 {
		config.GasMultiplier = big.NewFloat(rawConfig.GasMultiplier)
	}

	if rawConfig.BlockConfirmations != 0 {
		config.BlockConfirmations = big.NewInt(rawConfig.BlockConfirmations)
	}

	return config, nil
}

// func ParseConfig(rawConfig *chains.RawChainConfig) (*EVMConfig, error) {
// 	config := &EVMConfig{
// 		GeneralChainConfig: rawConfig.GeneralChainConfig,
// 		Bridge:             "",
// 		Erc20Handler:       "",
// 		Erc721Handler:      "",
// 		GenericHandler:     "",
// 		GasLimit:           big.NewInt(DefaultGasLimit),
// 		MaxGasPrice:        big.NewInt(DefaultGasPrice),
// 		GasMultiplier:      big.NewFloat(DefaultGasMultiplier),
// 		Http:               false,
// 		StartBlock:         big.NewInt(0),
// 		BlockConfirmations: big.NewInt(DefaultBlockConfirmations),
// 		EgsApiKey:          "",
// 		EgsSpeed:           "",
// 	}

// 	if contract, ok := rawConfig.Opts[BridgeOpt]; ok && contract != "" {
// 		config.Bridge = contract
// 		delete(rawConfig.Opts, BridgeOpt)
// 	} else {
// 		return nil, fmt.Errorf("must provide opts.bridge field for ethereum config")
// 	}

// 	if contract, ok := rawConfig.Opts[Erc20HandlerOpt]; ok {
// 		config.Erc20Handler = contract
// 		delete(rawConfig.Opts, Erc20HandlerOpt)
// 	}

// 	if contract, ok := rawConfig.Opts[Erc721HandlerOpt]; ok {
// 		config.Erc721Handler = contract
// 		delete(rawConfig.Opts, Erc721HandlerOpt)
// 	}

// 	if contract, ok := rawConfig.Opts[GenericHandlerOpt]; ok {
// 		config.GenericHandler = contract
// 		delete(rawConfig.Opts, GenericHandlerOpt)
// 	}

// 	if gasPrice, ok := rawConfig.Opts[MaxGasPriceOpt]; ok {
// 		price := big.NewInt(0)
// 		_, pass := price.SetString(gasPrice, 10)
// 		if pass {
// 			config.MaxGasPrice = price
// 			delete(rawConfig.Opts, MaxGasPriceOpt)
// 		} else {
// 			return nil, errors.New("unable to parse max gas price")
// 		}
// 	}

// 	if gasLimit, ok := rawConfig.Opts[GasLimitOpt]; ok {
// 		limit := big.NewInt(0)
// 		_, pass := limit.SetString(gasLimit, 10)
// 		if pass {
// 			config.GasLimit = limit
// 			delete(rawConfig.Opts, GasLimitOpt)
// 		} else {
// 			return nil, errors.New("unable to parse gas limit")
// 		}
// 	}

// 	if gasMultiplier, ok := rawConfig.Opts[GasMultiplier]; ok {
// 		multiplier := big.NewFloat(1)
// 		_, pass := multiplier.SetString(gasMultiplier)
// 		if pass {
// 			config.GasMultiplier = multiplier
// 			delete(rawConfig.Opts, GasMultiplier)
// 		} else {
// 			return nil, errors.New("unable to parse gasMultiplier to float")
// 		}
// 	}

// 	if HTTP, ok := rawConfig.Opts[HttpOpt]; ok && HTTP == "true" {
// 		config.Http = true
// 		delete(rawConfig.Opts, HttpOpt)
// 	} else if HTTP, ok := rawConfig.Opts[HttpOpt]; ok && HTTP == "false" {
// 		config.Http = false
// 		delete(rawConfig.Opts, HttpOpt)
// 	}

// 	if startBlock, ok := rawConfig.Opts[StartBlockOpt]; ok && startBlock != "" {
// 		block := big.NewInt(0)
// 		_, pass := block.SetString(startBlock, 10)
// 		if pass {
// 			config.startBlock = block
// 			delete(rawConfig.Opts, StartBlockOpt)
// 		} else {
// 			return nil, fmt.Errorf("unable to parse %s", StartBlockOpt)
// 		}
// 	}

// 	if blockConfirmations, ok := rawConfig.Opts[BlockConfirmationsOpt]; ok && blockConfirmations != "" {
// 		val := big.NewInt(DefaultBlockConfirmations)
// 		_, pass := val.SetString(blockConfirmations, 10)
// 		if pass {
// 			config.blockConfirmations = val
// 			delete(rawConfig.Opts, BlockConfirmationsOpt)
// 		} else {
// 			return nil, fmt.Errorf("unable to parse %s", BlockConfirmationsOpt)
// 		}
// 	} else {
// 		config.blockConfirmations = big.NewInt(DefaultBlockConfirmations)
// 		delete(rawConfig.Opts, BlockConfirmationsOpt)
// 	}

// 	if gsnApiKey, ok := rawConfig.Opts[EGSApiKey]; ok && gsnApiKey != "" {
// 		config.egsApiKey = gsnApiKey
// 		delete(rawConfig.Opts, EGSApiKey)
// 	}

// 	// TODO: change speed enums to be in separate package for querying egs
// 	if speed, ok := rawConfig.Opts[EGSSpeed]; ok && speed == "average" || speed == "fast" || speed == "fastest" {
// 		config.egsSpeed = speed
// 		delete(rawConfig.Opts, EGSSpeed)
// 	} else {
// 		// Default to "fast"
// 		config.egsSpeed = "fast"
// 		delete(rawConfig.Opts, EGSSpeed)
// 	}

// 	if len(rawConfig.Opts) != 0 {
// 		return nil, fmt.Errorf("unknown Opts Encountered: %v", rawConfig.Opts)
// 	}

// 	return config, nil
// }
