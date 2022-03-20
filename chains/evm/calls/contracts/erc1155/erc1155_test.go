package erc1155_test

import (
	"math/big"
	"testing"

	erc1155 "github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/erc1155"
	mock_calls "github.com/ChainSafe/chainbridge-core/chains/evm/calls/mock"
	mock_transactor "github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/signAndSend"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ERC1155ContractCallsTestSuite struct {
	suite.Suite
	gomockController                   *gomock.Controller
	mockContractCallerDispatcherClient *mock_calls.MockContractCallerDispatcher
	mockTransactor                     *mock_transactor.MockTransactor
	erc1155contract                    *erc1155.ERC1155Contract
}

var (
	testContractAddress   = "0x5f75ce92326e304962b22749bd71e36976171285"
	testInteractorAddress = "0x8362bbbd6d987895E2A4630a55e69Dd8C7b9f87B"
)

func TestRunERC1155ContractCallsTestSuite(t *testing.T) {
	suite.Run(t, new(ERC1155ContractCallsTestSuite))
}

func (s *ERC1155ContractCallsTestSuite) SetupSuite()    {}
func (s *ERC1155ContractCallsTestSuite) TearDownSuite() {}
func (s *ERC1155ContractCallsTestSuite) SetupTest() {
	s.gomockController = gomock.NewController(s.T())
	s.mockContractCallerDispatcherClient = mock_calls.NewMockContractCallerDispatcher(s.gomockController)
	s.mockTransactor = mock_transactor.NewMockTransactor(s.gomockController)
	s.erc1155contract = erc1155.NewErc1155Contract(
		s.mockContractCallerDispatcherClient, common.HexToAddress(testContractAddress), s.mockTransactor,
	)
}
func (s *ERC1155ContractCallsTestSuite) TearDownTest() {}

func (s *ERC1155ContractCallsTestSuite) TestErc1155Contract_GetBalance_Success() {
	s.mockContractCallerDispatcherClient.EXPECT().From().Return(common.HexToAddress(testInteractorAddress))
	s.mockContractCallerDispatcherClient.EXPECT().CallContract(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5}, nil)
	res, err := s.erc1155contract.GetBalance(common.HexToAddress(testInteractorAddress), big.NewInt(0))
	s.Equal(
		big.NewInt(5),
		res,
	)
	s.Nil(err)
}

func (s *ERC1155ContractCallsTestSuite) TestErc1155Contract_MintTokens_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{1, 2, 3, 4, 5}, nil)
	res, err := s.erc1155contract.Mint(common.HexToAddress(testInteractorAddress), big.NewInt(0), big.NewInt(10), []byte("0x"), signAndSend.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{1, 2, 3, 4, 5},
		res,
	)
	s.Nil(err)
}

func (s *ERC1155ContractCallsTestSuite) TestErc1155Contract_ApproveTokens_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{1, 2, 3, 4, 5, 6, 7, 8, 9}, nil)
	res, err := s.erc1155contract.Approve(common.HexToAddress(testInteractorAddress), true, signAndSend.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{1, 2, 3, 4, 5, 6, 7, 8, 9},
		res,
	)
	s.Nil(err)
}

func (s *ERC1155ContractCallsTestSuite) TestErc1155Contract_MinterRole_Success() {
	s.mockContractCallerDispatcherClient.EXPECT().From().Return(common.HexToAddress(testInteractorAddress))
	s.mockContractCallerDispatcherClient.EXPECT().CallContract(
		gomock.Any(),
		gomock.Any(),
		nil,
	).Return([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10}, nil)
	res, err := s.erc1155contract.MinterRole()
	s.Equal(
		[32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10},
		res,
	)
	s.Nil(err)
}

func (s *ERC1155ContractCallsTestSuite) TestErc1155Contract_AddMinter_Success() {
	s.mockContractCallerDispatcherClient.EXPECT().From().Return(common.HexToAddress(testInteractorAddress))
	s.mockContractCallerDispatcherClient.EXPECT().CallContract(
		gomock.Any(),
		gomock.Any(),
		nil,
	).Return([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 25}, nil)
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{1, 2, 3}, nil)
	res, err := s.erc1155contract.AddMinter(common.HexToAddress(testInteractorAddress), signAndSend.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{1, 2, 3},
		res,
	)
	s.Nil(err)
}
