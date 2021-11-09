package calls_test

import (
	"errors"

	calls "github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	mock_calls "github.com/ChainSafe/chainbridge-core/chains/evm/calls/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"math/big"

	"testing"
)

type DeployTestSuite struct {
	suite.Suite
	mockClientDeployer *mock_calls.MockClientDeployer
	mockgasPricer      *mock_calls.MockGasPricer
}

func TestRunDeployTestSuite(t *testing.T) {
	suite.Run(t, new(DeployTestSuite))
}

func (s *DeployTestSuite) SetupSuite() {
	gomockController := gomock.NewController(s.T())
	s.mockClientDeployer = mock_calls.NewMockClientDeployer(gomockController)
	s.mockgasPricer = mock_calls.NewMockGasPricer(gomockController)
}
func (s *DeployTestSuite) TearDownSuite() {}
func (s *DeployTestSuite) SetupTest()     {}
func (s *DeployTestSuite) TearDownTest()  {}

func (s *DeployTestSuite) TestDeployErc20NonceUnlockCallWithErrorThrown() {
	s.mockClientDeployer.EXPECT().LockNonce().Times(1)
	s.mockClientDeployer.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockgasPricer.EXPECT().GasPrice().Return([]*big.Int{big.NewInt(10)}, nil)
	s.mockClientDeployer.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Times(1).Return(common.Hash{}, nil)
	s.mockClientDeployer.EXPECT().WaitAndReturnTxReceipt(gomock.Any()).Times(1).Return(nil, nil)
	s.mockClientDeployer.EXPECT().From().Times(1).Return(common.Address{})
	s.mockClientDeployer.EXPECT().UnsafeIncreaseNonce().Times(1)
	s.mockClientDeployer.EXPECT().CodeAt(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
	s.mockClientDeployer.EXPECT().UnlockNonce().Times(1)

	_, _ = calls.DeployErc20(
		s.mockClientDeployer,
		evmtransaction.NewTransaction,
		s.mockgasPricer,
		"TEST",
		"TST")
}

func (s *DeployTestSuite) TestDeployErc20NonceUnlockCallWithoutErrorsThrown() {
	s.mockClientDeployer.EXPECT().LockNonce().Times(1)
	s.mockClientDeployer.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockgasPricer.EXPECT().GasPrice().Return([]*big.Int{big.NewInt(10)}, nil)
	s.mockClientDeployer.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Times(1).Return(common.Hash{}, errors.New("error"))
	s.mockClientDeployer.EXPECT().UnlockNonce().Times(1)
	s.mockClientDeployer.EXPECT().WaitAndReturnTxReceipt(gomock.Any()).Times(0)
	s.mockClientDeployer.EXPECT().UnsafeIncreaseNonce().Times(0)

	_, _ = calls.DeployErc20(
		s.mockClientDeployer,
		evmtransaction.NewTransaction,
		s.mockgasPricer,
		"TEST",
		"TST")
}

func (s *DeployTestSuite) TestDeployBridgeNonceUnlockCall() {
	s.mockClientDeployer.EXPECT().LockNonce().Times(1)
	s.mockClientDeployer.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockgasPricer.EXPECT().GasPrice().Return([]*big.Int{big.NewInt(10)}, nil)
	s.mockClientDeployer.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Times(1).Return(common.Hash{}, errors.New("error"))
	s.mockClientDeployer.EXPECT().UnlockNonce().Times(1)
	s.mockClientDeployer.EXPECT().WaitAndReturnTxReceipt(gomock.Any()).Times(0)
	s.mockClientDeployer.EXPECT().UnsafeIncreaseNonce().Times(0)

	toAddress := common.HexToAddress("0xtest1")

	_, _ = calls.DeployBridge(
		s.mockClientDeployer,
		evmtransaction.NewTransaction,
		s.mockgasPricer,
		0x1,
		[]common.Address{toAddress},
		big.NewInt(2),
		big.NewInt(10))
}

func (s *DeployTestSuite) TestDeployBridgeNonceUnlockCallWithoutErrorsThrown() {
	s.mockClientDeployer.EXPECT().LockNonce().Times(1)
	s.mockClientDeployer.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockgasPricer.EXPECT().GasPrice().Return([]*big.Int{big.NewInt(10)}, nil)
	s.mockClientDeployer.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Times(1).Return(common.Hash{}, errors.New("error"))
	s.mockClientDeployer.EXPECT().UnlockNonce().Times(1)
	s.mockClientDeployer.EXPECT().WaitAndReturnTxReceipt(gomock.Any()).Times(0)
	s.mockClientDeployer.EXPECT().UnsafeIncreaseNonce().Times(0)

	toAddress := common.HexToAddress("0xtest1")

	_, _ = calls.DeployBridge(
		s.mockClientDeployer,
		evmtransaction.NewTransaction,
		s.mockgasPricer,
		0x1,
		[]common.Address{toAddress},
		big.NewInt(2),
		big.NewInt(10))
}

func (s *DeployTestSuite) TestDeployErc20HandlerNonceUnlockCallWithErrorThrown() {
	s.mockClientDeployer.EXPECT().LockNonce().Times(1)
	s.mockClientDeployer.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockgasPricer.EXPECT().GasPrice().Return([]*big.Int{big.NewInt(10)}, nil)
	s.mockClientDeployer.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Times(1).Return(common.Hash{}, errors.New("error"))
	s.mockClientDeployer.EXPECT().UnlockNonce().Times(1)
	s.mockClientDeployer.EXPECT().WaitAndReturnTxReceipt(gomock.Any()).Times(0)
	s.mockClientDeployer.EXPECT().UnsafeIncreaseNonce().Times(0)

	toAddress := common.HexToAddress("0xtest1")

	_, _ = calls.DeployErc20Handler(
		s.mockClientDeployer,
		evmtransaction.NewTransaction,
		s.mockgasPricer,
		toAddress)
}

func (s *DeployTestSuite) TestDeployErc20HandlerNonceUnlockCallWithoutErrorsThrown() {
	s.mockClientDeployer.EXPECT().LockNonce().Times(1)
	s.mockClientDeployer.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockgasPricer.EXPECT().GasPrice().Return([]*big.Int{big.NewInt(10)}, nil)
	s.mockClientDeployer.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Times(1).Return(common.Hash{}, errors.New("error"))
	s.mockClientDeployer.EXPECT().UnlockNonce().Times(1)
	s.mockClientDeployer.EXPECT().WaitAndReturnTxReceipt(gomock.Any()).Times(0)
	s.mockClientDeployer.EXPECT().UnsafeIncreaseNonce().Times(0)

	toAddress := common.HexToAddress("0xtest1")

	_, _ = calls.DeployErc20Handler(
		s.mockClientDeployer,
		evmtransaction.NewTransaction,
		s.mockgasPricer,
		toAddress)
}
