package optimismclient

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ChainSafe/chainbridge-core/config"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/spf13/viper"
)

const DefaultGasLimit = 6721975
const DefaultGasPrice = 15000000
const DefaultGasMultiplier = 1
const DefaultBlockConfirmations = 10

type OptimismConfig struct {
	SharedEVMConfig  config.SharedEVMConfig // Pass an Optimism verifier replica for the endpoint to handle all chain listening and reads
	kp               *secp256k1.Keypair
	VerifyRollup     bool
	VerifierEndpoint string // This is the endpoint for the Optimism verifier and is purely used for verifying transactions
}

type RawOptimismConfig struct {
	config.RawSharedEVMConfig `mapstructure:",squash"`
	VerifyRollup              bool   `mapstructure:"verifyRollup"`
	VerifierEndpoint          string `mapstructure:"verifierEndpoint"`
}

func NewConfig() *OptimismConfig {
	return &OptimismConfig{}
}

func GetConfig(path string, name string) (*RawOptimismConfig, error) {
	config := &RawOptimismConfig{}

	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read in the config file, error: %w", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config into struct, error: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func ParseConfig(rawConfig *RawOptimismConfig) (*OptimismConfig, error) {

	cfg, err := rawConfig.RawSharedEVMConfig.ParseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to parse shared evm config, error: %w", err)
	}

	config := &OptimismConfig{
		SharedEVMConfig:  *cfg,
		VerifyRollup:     rawConfig.VerifyRollup,
		VerifierEndpoint: rawConfig.VerifierEndpoint,
	}

	return config, nil
}

func (c *RawOptimismConfig) ToJSON(file string) *os.File {
	var (
		newFile *os.File
		err     error
	)

	var raw []byte
	if raw, err = json.Marshal(&c); err != nil {
		fmt.Println("error marshalling json", "err", err)
		os.Exit(1)
	}

	newFile, err = os.Create(file)
	if err != nil {
		fmt.Println("error creating config file", "err", err)
	}
	_, err = newFile.Write(raw)
	if err != nil {
		fmt.Println("error writing to config file", "err", err)
	}

	if err := newFile.Close(); err != nil {
		fmt.Println("failed to unmarshal config into struct", "err", err)
	}
	return newFile
}
