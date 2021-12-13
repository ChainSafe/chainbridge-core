package contracts

import (
	"errors"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	mock_calls "github.com/ChainSafe/chainbridge-core/chains/evm/calls/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	mock_transactor "github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/mock"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"math/big"
	"strings"
	"testing"
)

type ContractTestSuite struct {
	suite.Suite
	gomockController                   *gomock.Controller
	mockContractCallerDispatcherClient *mock_calls.MockContractCallerDispatcher
	mockTransactor                     *mock_transactor.MockTransactor
	contract                           Contract
}

func TestRunContractTestSuite(t *testing.T) {
	suite.Run(t, new(ContractTestSuite))
}

func (s *ContractTestSuite) SetupSuite()    {}
func (s *ContractTestSuite) TearDownSuite() {}
func (s *ContractTestSuite) SetupTest() {
	s.gomockController = gomock.NewController(s.T())
	s.mockContractCallerDispatcherClient = mock_calls.NewMockContractCallerDispatcher(s.gomockController)
	s.mockTransactor = mock_transactor.NewMockTransactor(s.gomockController)
	// Use ERC721 contract ABI inside the contract test
	a, _ := abi.JSON(strings.NewReader(consts.ERC721PresetMinterPauserABI))
	b := common.FromHex(consts.ERC721PresetMinterPauserBin)
	s.contract = NewContract(
		common.Address{}, a, b, s.mockContractCallerDispatcherClient, s.mockTransactor,
	)
}
func (s *ContractTestSuite) TearDownTest() {}

func (s *ContractTestSuite) TestContract_PackMethod_ValidRequest_Success() {
	res, err := s.contract.PackMethod("approve", common.Address{}, big.NewInt(10))
	s.Equal(
		common.Hex2Bytes("095ea7b30000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a"),
		res,
	)
	s.Nil(err)
}

func (s *ContractTestSuite) TestContract_PackMethod_InvalidRequest_Fail() {
	res, err := s.contract.PackMethod("invalid_method", common.Address{}, big.NewInt(10))
	s.Equal([]byte{}, res)
	s.NotNil(err)
}

func (s *ContractTestSuite) TestContract_UnpackResult_InvalidRequest_Fail() {
	rawInvalidApproveData := common.Hex2Bytes("095ea7b30000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a")
	res, err := s.contract.UnpackResult("approve", rawInvalidApproveData)
	s.NotNil(err)
	s.Nil(res)
}

func (s *ContractTestSuite) TestContract_ExecuteTransaction_ValidRequest_Success() {
	s.mockTransactor.EXPECT().Transact(
		&common.Address{},
		gomock.Any(),
		transactor.TransactOptions{},
	).Return(&common.Hash{}, nil)
	hash, err := s.contract.ExecuteTransaction(
		"approve",
		transactor.TransactOptions{}, common.Address{}, big.NewInt(10),
	)
	s.Nil(err)
	s.NotNil(hash)
}

func (s *ContractTestSuite) TestContract_ExecuteTransaction_TransactError_Fail() {
	s.mockTransactor.EXPECT().Transact(
		&common.Address{},
		gomock.Any(),
		transactor.TransactOptions{},
	).Return(nil, errors.New("error"))
	hash, err := s.contract.ExecuteTransaction(
		"approve",
		transactor.TransactOptions{}, common.Address{}, big.NewInt(10),
	)
	s.Nil(hash)
	s.Error(err, "error")
}

func (s *ContractTestSuite) TestContract_ExecuteTransaction_InvalidRequest_Fail() {
	hash, err := s.contract.ExecuteTransaction(
		"approve",
		transactor.TransactOptions{}, common.Address{}, // missing one argument
	)
	s.Nil(hash)
	s.Error(err, "error")
}

func (s *ContractTestSuite) TestContract_CallContract_CallContractError_Fail() {
	s.mockContractCallerDispatcherClient.EXPECT().CallContract(
		gomock.Any(),
		gomock.Any(),
		nil,
	).Return(nil, errors.New("error"))
	s.mockContractCallerDispatcherClient.EXPECT().From().Times(1).Return(common.Address{})

	res, err := s.contract.CallContract("ownerOf", big.NewInt(0))
	if err != nil {
		return
	}
	s.Nil(res)
	s.Error(err, "error")
}

func (s *ContractTestSuite) TestContract_CallContract_InvalidRequest_Fail() {
	res, err := s.contract.CallContract("invalidMethod", big.NewInt(0))
	if err != nil {
		return
	}
	s.Nil(res)
	s.Error(err, "error")
}

func (s *ContractTestSuite) TestContract_CallContract_MissingContract_Fail() {
	s.mockContractCallerDispatcherClient.EXPECT().CallContract(
		gomock.Any(),
		gomock.Any(),
		nil,
	).Return(nil, errors.New("error"))
	s.mockContractCallerDispatcherClient.EXPECT().From().Times(1).Return(common.Address{})
	res, err := s.contract.CallContract("ownerOf", big.NewInt(0))
	if err != nil {
		return
	}
	s.Nil(res)
	s.Error(err, "error")
}

func (s *ContractTestSuite) TestContract_DeployContract_InvalidRequest_Fail() {
	res, err := s.contract.DeployContract("invalid_param")
	s.Equal(common.Address{}, res)
	s.Error(err, "error")
}

func (s *ContractTestSuite) TestContract_DeployContract_TransactionError_Fail() {
	s.mockTransactor.EXPECT().Transact(
		nil, gomock.Any(), gomock.Any(),
	).Times(1).Return(&common.Hash{}, errors.New("error"))
	res, err := s.contract.DeployContract("TestERC721", "TST721", "")
	s.Equal(common.Address{}, res)
	s.Error(err, "error")
}

func (s *ContractTestSuite) TestContract_DeployContract_GetTxByHashError_Fail() {
	s.mockTransactor.EXPECT().Transact(
		nil, gomock.Any(), gomock.Any(),
	).Times(1).Return(&common.Hash{}, nil)
	s.mockContractCallerDispatcherClient.EXPECT().GetTransactionByHash(
		common.Hash{},
	).Return(nil, false, errors.New("error"))

	res, err := s.contract.DeployContract("TestERC721", "TST721", "")
	s.Equal(common.Address{}, res)
	s.Error(err, "error")
}
