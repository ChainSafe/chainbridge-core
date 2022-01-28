package chain_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/config/chain"
	"github.com/stretchr/testify/suite"
)

type NewEVMConfigTestSuite struct {
	suite.Suite
}

func TestRunNewEVMConfigTestSuite(t *testing.T) {
	suite.Run(t, new(NewEVMConfigTestSuite))
}

func (s *NewEVMConfigTestSuite) SetupSuite()    {}
func (s *NewEVMConfigTestSuite) TearDownSuite() {}
func (s *NewEVMConfigTestSuite) SetupTest()     {}
func (s *NewEVMConfigTestSuite) TearDownTest()  {}

func (s *NewEVMConfigTestSuite) Test_FailedDecode() {
	_, err := chain.NewEVMConfig(map[string]interface{}{
		"gasLimit": "invalid",
	})

	s.NotNil(err)
}

func (s *NewEVMConfigTestSuite) Test_FailedGeneralConfigValidation() {
	_, err := chain.NewEVMConfig(map[string]interface{}{})

	s.NotNil(err)
}

func (s *NewEVMConfigTestSuite) Test_FailedEVMConfigValidation() {
	_, err := chain.NewEVMConfig(map[string]interface{}{
		"id":       1,
		"endpoint": "ws://domain.com",
		"name":     "evm1",
		"from":     "address",
	})

	s.NotNil(err)
}

func (s *NewEVMConfigTestSuite) Test_InvalidBlockConfirmation() {
	_, err := chain.NewEVMConfig(map[string]interface{}{
		"id":                 1,
		"endpoint":           "ws://domain.com",
		"name":               "evm1",
		"from":               "address",
		"bridge":             "bridgeAddress",
		"blockConfirmations": -1,
	})

	s.NotNil(err)
	s.Equal(err.Error(), "blockConfirmations has to be >=1")
}

func (s *NewEVMConfigTestSuite) Test_ValidConfig() {
	rawConfig := map[string]interface{}{
		"id":       1,
		"endpoint": "ws://domain.com",
		"name":     "evm1",
		"from":     "address",
		"bridge":   "bridgeAddress",
	}

	actualConfig, err := chain.NewEVMConfig(rawConfig)

	id := new(uint8)
	*id = 1
	s.Nil(err)
	s.Equal(*actualConfig, chain.EVMConfig{
		GeneralChainConfig: chain.GeneralChainConfig{
			Name:     "evm1",
			From:     "address",
			Endpoint: "ws://domain.com",
			Id:       id,
		},
		Bridge:             "bridgeAddress",
		Erc20Handler:       "",
		Erc721Handler:      "",
		GenericHandler:     "",
		GasLimit:           big.NewInt(consts.DefaultGasLimit),
		MaxGasPrice:        big.NewInt(consts.DefaultGasPrice),
		GasMultiplier:      big.NewFloat(consts.DefaultGasMultiplier),
		StartBlock:         big.NewInt(0),
		BlockConfirmations: big.NewInt(consts.DefaultBlockConfirmations),
		BlockRetryInterval: time.Duration(5) * time.Second,
	})
}

func (s *NewEVMConfigTestSuite) Test_ValidConfigWithCustomTxParams() {
	rawConfig := map[string]interface{}{
		"id":                 1,
		"endpoint":           "ws://domain.com",
		"name":               "evm1",
		"from":               "address",
		"bridge":             "bridgeAddress",
		"maxGasPrice":        1000,
		"gasMultiplier":      1000,
		"gasLimit":           1000,
		"startBlock":         1000,
		"blockConfirmations": 10,
		"blockRetryInterval": 10,
	}

	actualConfig, err := chain.NewEVMConfig(rawConfig)

	id := new(uint8)
	*id = 1
	s.Nil(err)
	s.Equal(*actualConfig, chain.EVMConfig{
		GeneralChainConfig: chain.GeneralChainConfig{
			Name:     "evm1",
			From:     "address",
			Endpoint: "ws://domain.com",
			Id:       id,
		},
		Bridge:             "bridgeAddress",
		Erc20Handler:       "",
		Erc721Handler:      "",
		GenericHandler:     "",
		GasLimit:           big.NewInt(1000),
		MaxGasPrice:        big.NewInt(1000),
		GasMultiplier:      big.NewFloat(1000),
		StartBlock:         big.NewInt(1000),
		BlockConfirmations: big.NewInt(10),
		BlockRetryInterval: time.Duration(10) * time.Second,
	})
}
