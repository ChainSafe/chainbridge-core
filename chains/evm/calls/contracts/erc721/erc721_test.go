package erc721_test

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/erc721"
	mock_calls "github.com/ChainSafe/chainbridge-core/chains/evm/calls/mock"
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	mock_transactor "github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/mock"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ERC721CallsTestSuite struct {
	suite.Suite
	gomockController                   *gomock.Controller
	clientMock                         *mock_calls.MockClientDispatcher
	mockContractCallerDispatcherClient *mock_calls.MockContractCallerDispatcher
	mockTransactor                     *mock_transactor.MockTransactor
	erc721ContractAddress              common.Address
	erc721Contract                     *erc721.ERC721Contract
}

func TestRunERC721CallsTestSuite(t *testing.T) {
	suite.Run(t, new(ERC721CallsTestSuite))
}

func (s *ERC721CallsTestSuite) SetupSuite()    {}
func (s *ERC721CallsTestSuite) TearDownSuite() {}
func (s *ERC721CallsTestSuite) SetupTest() {
	s.gomockController = gomock.NewController(s.T())
	s.clientMock = mock_calls.NewMockClientDispatcher(s.gomockController)
	s.mockContractCallerDispatcherClient = mock_calls.NewMockContractCallerDispatcher(s.gomockController)
	s.mockTransactor = mock_transactor.NewMockTransactor(s.gomockController)
	s.erc721ContractAddress = common.HexToAddress("0x9A0E6F91E6031C08326764655432f8F9c180fBa0")
	s.erc721Contract = erc721.NewErc721Contract(s.mockContractCallerDispatcherClient, s.erc721ContractAddress, s.mockTransactor)
}
func (s *ERC721CallsTestSuite) TearDownTest() {}

func (s *ERC721CallsTestSuite) TestERC721Contract_PackMethod_ValidRequest_Success() {
	res, err := s.erc721Contract.PackMethod("approve", common.Address{}, big.NewInt(10))
	s.Equal(
		common.Hex2Bytes("095ea7b30000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a"),
		res,
	)
	s.Nil(err)
}

func (s *ERC721CallsTestSuite) TestERC721Contract_PackMethod_InvalidNumberOfArguments_Fail() {
	res, err := s.erc721Contract.PackMethod("approve", common.Address{})
	s.Equal(
		[]byte{},
		res,
	)
	s.Error(err)
}

func (s *ERC721CallsTestSuite) TestERC721Contract_PackMethod_NotExistingMethod_Fail() {
	res, err := s.erc721Contract.PackMethod("fail", common.Address{})
	s.Equal(
		[]byte{},
		res,
	)
	s.Error(err)
}

func (s *ERC721CallsTestSuite) TestERC721Contract_UnpackResult_InvalidData_Fail() {
	rawData := common.Hex2Bytes("095ea7b30000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a")
	res, err := s.erc721Contract.UnpackResult("approve", rawData)
	s.NotNil(err)
	s.Nil(res)
}

func (s *ERC721CallsTestSuite) TestERC721Contract_Approve_Success() {
	s.mockTransactor.EXPECT().Transact(
		&s.erc721ContractAddress,
		gomock.Any(),
		transactor.TransactOptions{},
	).Return(&common.Hash{}, nil)

	res, err := s.erc721Contract.Approve(
		big.NewInt(1),
		common.HexToAddress("0x9FD320F352539E8A0E9be4B63c91395575420Aac"),
		transactor.TransactOptions{},
	)

	s.Nil(err)
	s.NotNil(res)
}
