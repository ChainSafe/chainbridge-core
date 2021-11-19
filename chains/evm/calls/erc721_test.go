package calls_test

import (
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	mock_utils "github.com/ChainSafe/chainbridge-core/chains/evm/calls/mock"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ERC721CallsTestSuite struct {
	suite.Suite
	gomockController *gomock.Controller
	clientMock       *mock_utils.MockClientDispatcher
}

func TestRunERC721CallsTestSuite(t *testing.T) {
	suite.Run(t, new(ERC721CallsTestSuite))
}

func (s *ERC721CallsTestSuite) SetupSuite()    {}
func (s *ERC721CallsTestSuite) TearDownSuite() {}
func (s *ERC721CallsTestSuite) SetupTest() {
	s.gomockController = gomock.NewController(s.T())
	s.clientMock = mock_utils.NewMockClientDispatcher(s.gomockController)
}
func (s *ERC721CallsTestSuite) TearDownTest() {}

func (s *ERC721CallsTestSuite) TestERC721Approve_ValidRequest_Success() {
	abi, res, err := calls.PackERC721Method("approve", common.Address{}, big.NewInt(10))
	s.Equal(
		common.Hex2Bytes("095ea7b30000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a"),
		res,
	)
	s.Nil(err)
	s.NotNil(abi)
}

func (s *ERC721CallsTestSuite) TestERC721Approve_InvalidNumberOfArguments_Fail() {
	abi, res, err := calls.PackERC721Method("approve", common.Address{})
	s.Equal(
		[]byte{},
		res,
	)
	s.Error(err)
	s.NotNil(abi)
}

func (s *ERC721CallsTestSuite) TestERC721Approve_NotExistingMethod_Fail() {
	abi, res, err := calls.PackERC721Method("fail", common.Address{})
	s.Equal(
		[]byte{},
		res,
	)
	s.Error(err)
	s.NotNil(abi)
}
