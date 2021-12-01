package chain_test

import (
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-core/config/chain"
	"github.com/stretchr/testify/suite"
)

type NewSubstrateConfigTestSuite struct {
	suite.Suite
}

func TestRunNewSubstrateConfigTestSuite(t *testing.T) {
	suite.Run(t, new(NewSubstrateConfigTestSuite))
}

func (s *NewSubstrateConfigTestSuite) SetupSuite()    {}
func (s *NewSubstrateConfigTestSuite) TearDownSuite() {}
func (s *NewSubstrateConfigTestSuite) SetupTest()     {}
func (s *NewSubstrateConfigTestSuite) TearDownTest()  {}

func (s *NewSubstrateConfigTestSuite) Test_FailedDecode() {
	_, err := chain.NewSubstrateConfig(map[string]interface{}{
		"startBlock": "invalid",
	})

	s.NotNil(err)
}

func (s *NewSubstrateConfigTestSuite) Test_FailedGeneralConfigValidation() {
	_, err := chain.NewSubstrateConfig(map[string]interface{}{})

	s.NotNil(err)
}

func (s *NewSubstrateConfigTestSuite) Test_ValidConfig() {
	rawConfig := map[string]interface{}{
		"id":       1,
		"endpoint": "ws://domain.com",
		"name":     "evm1",
		"from":     "address",
		"bridge":   "bridgeAddress",
	}

	actualConfig, err := chain.NewSubstrateConfig(rawConfig)

	id := new(uint8)
	*id = 1
	s.Nil(err)
	s.Equal(*actualConfig, chain.SubstrateConfig{
		GeneralChainConfig: chain.GeneralChainConfig{
			Name:     "evm1",
			From:     "address",
			Endpoint: "ws://domain.com",
			Id:       id,
		},
		StartBlock:      big.NewInt(0),
		UseExtendedCall: false,
	})
}
