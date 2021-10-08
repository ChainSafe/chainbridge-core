package calls_test

import (
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	mock_utils "github.com/ChainSafe/chainbridge-core/chains/evm/calls/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type GetSolidityFunctionSigTestSuite struct {
	suite.Suite
	gomockController *gomock.Controller
	clientMock       *mock_utils.MockChainClient
}

func TestRunGetSolidityFunctionSigTestSuite(t *testing.T) {
	suite.Run(t, new(GetSolidityFunctionSigTestSuite))
}

func (s *GetSolidityFunctionSigTestSuite) SetupSuite()    {}
func (s *GetSolidityFunctionSigTestSuite) TearDownSuite() {}
func (s *GetSolidityFunctionSigTestSuite) SetupTest() {
	s.gomockController = gomock.NewController(s.T())
}
func (s *GetSolidityFunctionSigTestSuite) TearDownTest() {}

func (s *GetSolidityFunctionSigTestSuite) TestReturnsValidSolidityFunctionSig() {
	sig := calls.GetSolidityFunctionSig([]byte("store(bytes32)"))

	s.Equal(sig, [4]byte{0x65, 0x4c, 0xf8, 0x8c})
}
