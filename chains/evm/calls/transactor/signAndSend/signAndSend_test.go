package signAndSend_test

import (
	"math/big"
	"testing"

	erc20 "github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/erc20"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmtransaction"
	mock_calls "github.com/ChainSafe/chainbridge-core/chains/evm/calls/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	mock_transactor "github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/signAndSend"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type TransactorTestSuite struct {
	suite.Suite
	gomockController                   *gomock.Controller
	mockContractCallerDispatcherClient *mock_calls.MockContractCallerDispatcher
	mockTransactor                     *mock_transactor.MockTransactor
	erc20ContractAddress               common.Address
	erc20Contract                      *erc20.ERC20Contract
	mockGasPricer                      *mock_calls.MockGasPricer
}

var (
	erc20ContractAddress = "0x829bd824b016326a401d083b33d092293333a830"
)

func TestERC20TestSuite(t *testing.T) {
	suite.Run(t, new(TransactorTestSuite))
}

func (s *TransactorTestSuite) SetupSuite()    {}
func (s *TransactorTestSuite) TearDownSuite() {}
func (s *TransactorTestSuite) SetupTest() {
	s.gomockController = gomock.NewController(s.T())
	s.mockContractCallerDispatcherClient = mock_calls.NewMockContractCallerDispatcher(s.gomockController)
	s.erc20ContractAddress = common.HexToAddress(erc20ContractAddress)
	s.erc20Contract = erc20.NewERC20Contract(
		s.mockContractCallerDispatcherClient, common.HexToAddress(erc20ContractAddress), s.mockTransactor,
	)
	s.mockTransactor = mock_transactor.NewMockTransactor(s.gomockController)
	s.mockGasPricer = mock_calls.NewMockGasPricer(s.gomockController)
}

func (s *TransactorTestSuite) TestTransactor_SignAndSend_Success() {
	var byteData = []byte{47, 47, 241, 93, 159, 45, 240, 254, 210, 199, 118, 72, 222, 88, 96, 164, 204, 80, 140, 208, 129, 140, 133, 184, 184, 161, 171, 76, 238, 239, 141, 152, 28, 137, 86, 166, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 60, 48, 181, 109, 237, 4, 127, 230, 34, 95, 112, 4, 234, 75, 225, 174, 112, 201, 2, 106}

	s.mockContractCallerDispatcherClient.EXPECT().LockNonce()
	s.mockContractCallerDispatcherClient.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockGasPricer.EXPECT().GasPrice().Return([]*big.Int{big.NewInt(1)}, nil)
	s.mockContractCallerDispatcherClient.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Return(common.Hash{1, 2, 3, 4, 5}, nil)
	s.mockContractCallerDispatcherClient.EXPECT().WaitAndReturnTxReceipt(gomock.Any()).Return(&types.Receipt{}, nil)
	s.mockContractCallerDispatcherClient.EXPECT().UnsafeIncreaseNonce().Return(nil)
	s.mockContractCallerDispatcherClient.EXPECT().UnlockNonce()

	txFabric := evmtransaction.NewTransaction
	var trans = signAndSend.NewSignAndSendTransactor(
		txFabric,
		s.mockGasPricer,
		s.mockContractCallerDispatcherClient,
	)
	txHash, err := trans.Transact(
		&common.Address{},
		byteData,
		transactor.TransactOptions{},
	)

	s.Nil(err)
	// without prepare flag omitted SignAndSendTransactor is used and output is normal tx hash
	s.Equal("0x0102030405000000000000000000000000000000000000000000000000000000", txHash.String())
}
