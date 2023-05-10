package monitored_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmtransaction"
	mock_calls "github.com/ChainSafe/chainbridge-core/chains/evm/calls/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	mock_transactor "github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/monitored"
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
	mockGasPricer                      *mock_calls.MockGasPricer
}

func TestMonitoredTransactorTestSuite(t *testing.T) {
	suite.Run(t, new(TransactorTestSuite))
}

func (s *TransactorTestSuite) SetupTest() {
	s.gomockController = gomock.NewController(s.T())
	s.mockContractCallerDispatcherClient = mock_calls.NewMockContractCallerDispatcher(s.gomockController)
	s.mockTransactor = mock_transactor.NewMockTransactor(s.gomockController)
	s.mockGasPricer = mock_calls.NewMockGasPricer(s.gomockController)
}

func (s *TransactorTestSuite) TestTransactor_SignAndSend_Success() {
	var byteData = []byte{47, 47, 241, 93, 159, 45, 240, 254, 210, 199, 118, 72, 222, 88, 96, 164, 204, 80, 140, 208, 129, 140, 133, 184, 184, 161, 171, 76, 238, 239, 141, 152, 28, 137, 86, 166, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 60, 48, 181, 109, 237, 4, 127, 230, 34, 95, 112, 4, 234, 75, 225, 174, 112, 201, 2, 106}

	s.mockContractCallerDispatcherClient.EXPECT().LockNonce()
	s.mockContractCallerDispatcherClient.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockGasPricer.EXPECT().GasPrice(gomock.Any()).Return([]*big.Int{big.NewInt(1)}, nil)
	s.mockContractCallerDispatcherClient.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Return(common.Hash{1, 2, 3, 4, 5}, nil)
	s.mockContractCallerDispatcherClient.EXPECT().UnsafeIncreaseNonce().Return(nil)
	s.mockContractCallerDispatcherClient.EXPECT().UnlockNonce()

	t := monitored.NewMonitoredTransactor(
		evmtransaction.NewTransaction,
		s.mockGasPricer,
		s.mockContractCallerDispatcherClient,
		big.NewInt(1000),
		big.NewInt(15))
	txHash, err := t.Transact(
		&common.Address{},
		byteData,
		transactor.TransactOptions{},
	)

	s.Nil(err)
	s.Equal("0x0102030405000000000000000000000000000000000000000000000000000000", txHash.String())
}

func (s *TransactorTestSuite) TestTransactor_SignAndSend_Fail() {
	var byteData = []byte{47, 47, 241, 93, 159, 45, 240, 254, 210, 199, 118, 72, 222, 88, 96, 164, 204, 80, 140, 208, 129, 140, 133, 184, 184, 161, 171, 76, 238, 239, 141, 152, 28, 137, 86, 166, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 60, 48, 181, 109, 237, 4, 127, 230, 34, 95, 112, 4, 234, 75, 225, 174, 112, 201, 2, 106}

	s.mockContractCallerDispatcherClient.EXPECT().LockNonce()
	s.mockContractCallerDispatcherClient.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockGasPricer.EXPECT().GasPrice(gomock.Any()).Return([]*big.Int{big.NewInt(1)}, nil)
	s.mockContractCallerDispatcherClient.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Return(common.Hash{}, fmt.Errorf("error"))
	s.mockContractCallerDispatcherClient.EXPECT().UnlockNonce()

	t := monitored.NewMonitoredTransactor(
		evmtransaction.NewTransaction,
		s.mockGasPricer,
		s.mockContractCallerDispatcherClient,
		big.NewInt(1000),
		big.NewInt(15))
	_, err := t.Transact(
		&common.Address{},
		byteData,
		transactor.TransactOptions{},
	)

	s.NotNil(err)
}

func (s *TransactorTestSuite) TestTransactor_MonitoredTransaction_SuccessfulExecution() {
	var byteData = []byte{47, 47, 241, 93, 159, 45, 240, 254, 210, 199, 118, 72, 222, 88, 96, 164, 204, 80, 140, 208, 129, 140, 133, 184, 184, 161, 171, 76, 238, 239, 141, 152, 28, 137, 86, 166, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 60, 48, 181, 109, 237, 4, 127, 230, 34, 95, 112, 4, 234, 75, 225, 174, 112, 201, 2, 106}

	// Sending transaction
	s.mockContractCallerDispatcherClient.EXPECT().LockNonce()
	s.mockContractCallerDispatcherClient.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockGasPricer.EXPECT().GasPrice(gomock.Any()).Return([]*big.Int{big.NewInt(1)}, nil)
	s.mockContractCallerDispatcherClient.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Return(common.Hash{1, 2, 3, 4, 5}, nil)
	s.mockContractCallerDispatcherClient.EXPECT().UnsafeIncreaseNonce().Return(nil)
	s.mockContractCallerDispatcherClient.EXPECT().UnlockNonce()

	ctx, cancel := context.WithCancel(context.Background())
	t := monitored.NewMonitoredTransactor(
		evmtransaction.NewTransaction,
		s.mockGasPricer,
		s.mockContractCallerDispatcherClient,
		big.NewInt(1000),
		big.NewInt(15))

	go t.Monitor(ctx, time.Millisecond*50, time.Minute, time.Millisecond)
	hash, err := t.Transact(
		&common.Address{},
		byteData,
		transactor.TransactOptions{},
	)
	// Transaction executed
	s.mockContractCallerDispatcherClient.EXPECT().TransactionReceipt(gomock.Any(), *hash).Return(&types.Receipt{
		Status: types.ReceiptStatusSuccessful,
	}, nil)
	s.Nil(err)

	time.Sleep(time.Millisecond * 150)
	cancel()
}

func (s *TransactorTestSuite) TestTransactor_MonitoredTransaction_TxTimeout() {
	var byteData = []byte{47, 47, 241, 93, 159, 45, 240, 254, 210, 199, 118, 72, 222, 88, 96, 164, 204, 80, 140, 208, 129, 140, 133, 184, 184, 161, 171, 76, 238, 239, 141, 152, 28, 137, 86, 166, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 60, 48, 181, 109, 237, 4, 127, 230, 34, 95, 112, 4, 234, 75, 225, 174, 112, 201, 2, 106}

	// Sending transaction
	s.mockContractCallerDispatcherClient.EXPECT().LockNonce()
	s.mockContractCallerDispatcherClient.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockGasPricer.EXPECT().GasPrice(gomock.Any()).Return([]*big.Int{big.NewInt(1)}, nil)
	s.mockContractCallerDispatcherClient.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Return(common.Hash{1, 2, 3, 4, 5}, nil)
	s.mockContractCallerDispatcherClient.EXPECT().UnsafeIncreaseNonce().Return(nil)
	s.mockContractCallerDispatcherClient.EXPECT().UnlockNonce()

	ctx, cancel := context.WithCancel(context.Background())
	t := monitored.NewMonitoredTransactor(
		evmtransaction.NewTransaction,
		s.mockGasPricer,
		s.mockContractCallerDispatcherClient,
		big.NewInt(1000),
		big.NewInt(15))

	go t.Monitor(ctx, time.Millisecond*50, time.Millisecond, time.Millisecond)
	hash, err := t.Transact(
		&common.Address{},
		byteData,
		transactor.TransactOptions{},
	)
	s.mockContractCallerDispatcherClient.EXPECT().TransactionReceipt(gomock.Any(), *hash).Return(nil, fmt.Errorf("not found"))
	s.Nil(err)

	time.Sleep(time.Millisecond * 150)
	cancel()
}

func (s *TransactorTestSuite) TestTransactor_MonitoredTransaction_TransactionResent() {
	var byteData = []byte{47, 47, 241, 93, 159, 45, 240, 254, 210, 199, 118, 72, 222, 88, 96, 164, 204, 80, 140, 208, 129, 140, 133, 184, 184, 161, 171, 76, 238, 239, 141, 152, 28, 137, 86, 166, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 60, 48, 181, 109, 237, 4, 127, 230, 34, 95, 112, 4, 234, 75, 225, 174, 112, 201, 2, 106}

	// Sending transaction
	s.mockContractCallerDispatcherClient.EXPECT().LockNonce()
	s.mockContractCallerDispatcherClient.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockGasPricer.EXPECT().GasPrice(gomock.Any()).Return([]*big.Int{big.NewInt(10)}, nil)
	s.mockContractCallerDispatcherClient.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Return(common.Hash{1, 2, 3, 4, 5}, nil)
	s.mockContractCallerDispatcherClient.EXPECT().UnsafeIncreaseNonce().Return(nil)
	s.mockContractCallerDispatcherClient.EXPECT().UnlockNonce()

	// Resending transaction
	s.mockContractCallerDispatcherClient.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Return(common.Hash{1, 2, 3, 4, 5}, nil)

	ctx, cancel := context.WithCancel(context.Background())
	t := monitored.NewMonitoredTransactor(
		evmtransaction.NewTransaction,
		s.mockGasPricer,
		s.mockContractCallerDispatcherClient,
		big.NewInt(1000),
		big.NewInt(15))

	go t.Monitor(ctx, time.Millisecond*50, time.Minute, time.Millisecond)
	hash, err := t.Transact(
		&common.Address{},
		byteData,
		transactor.TransactOptions{},
	)
	s.Nil(err)

	s.mockContractCallerDispatcherClient.EXPECT().TransactionReceipt(gomock.Any(), *hash).Return(nil, fmt.Errorf("not found"))
	s.mockContractCallerDispatcherClient.EXPECT().TransactionReceipt(gomock.Any(), common.Hash{1, 2, 3, 4, 5}).Return(&types.Receipt{
		Status: types.ReceiptStatusFailed,
	}, nil)

	time.Sleep(time.Millisecond * 125)
	cancel()
}

func (s *TransactorTestSuite) TestTransactor_MonitoredTransaction_MaxGasPriceReached() {
	var byteData = []byte{47, 47, 241, 93, 159, 45, 240, 254, 210, 199, 118, 72, 222, 88, 96, 164, 204, 80, 140, 208, 129, 140, 133, 184, 184, 161, 171, 76, 238, 239, 141, 152, 28, 137, 86, 166, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 60, 48, 181, 109, 237, 4, 127, 230, 34, 95, 112, 4, 234, 75, 225, 174, 112, 201, 2, 106}

	// Sending transaction
	s.mockContractCallerDispatcherClient.EXPECT().LockNonce()
	s.mockContractCallerDispatcherClient.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockGasPricer.EXPECT().GasPrice(gomock.Any()).Return([]*big.Int{big.NewInt(11)}, nil)
	s.mockContractCallerDispatcherClient.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Return(common.Hash{1, 2, 3, 4, 5}, nil)
	s.mockContractCallerDispatcherClient.EXPECT().UnsafeIncreaseNonce().Return(nil)
	s.mockContractCallerDispatcherClient.EXPECT().UnlockNonce()

	// Resending transaction
	s.mockContractCallerDispatcherClient.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Return(common.Hash{1, 2, 3, 4, 5}, nil)

	ctx, cancel := context.WithCancel(context.Background())
	t := monitored.NewMonitoredTransactor(
		evmtransaction.NewTransaction,
		s.mockGasPricer,
		s.mockContractCallerDispatcherClient,
		big.NewInt(10),
		big.NewInt(15))

	go t.Monitor(ctx, time.Millisecond*50, time.Minute, time.Millisecond)
	hash, err := t.Transact(
		&common.Address{},
		byteData,
		transactor.TransactOptions{},
	)
	s.Nil(err)

	s.mockContractCallerDispatcherClient.EXPECT().TransactionReceipt(gomock.Any(), *hash).Return(nil, fmt.Errorf("not found"))
	s.mockContractCallerDispatcherClient.EXPECT().TransactionReceipt(gomock.Any(), common.Hash{1, 2, 3, 4, 5}).Return(&types.Receipt{
		Status: types.ReceiptStatusFailed,
	}, nil)

	time.Sleep(time.Millisecond * 125)
	cancel()
}

func (s *TransactorTestSuite) TestTransactor_IncreaseGas_15PercentIncrease() {
	t := monitored.NewMonitoredTransactor(
		evmtransaction.NewTransaction,
		s.mockGasPricer,
		s.mockContractCallerDispatcherClient,
		big.NewInt(150),
		big.NewInt(15))

	newGas := t.IncreaseGas([]*big.Int{big.NewInt(1), big.NewInt(10), big.NewInt(100)})

	s.Equal(newGas, []*big.Int{big.NewInt(2), big.NewInt(11), big.NewInt(115)})
}

func (s *TransactorTestSuite) TestTransactor_IncreaseGas_MaxGasReached() {
	t := monitored.NewMonitoredTransactor(
		evmtransaction.NewTransaction,
		s.mockGasPricer,
		s.mockContractCallerDispatcherClient,
		big.NewInt(15),
		big.NewInt(15))

	newGas := t.IncreaseGas([]*big.Int{big.NewInt(1), big.NewInt(10), big.NewInt(100)})

	s.Equal(newGas, []*big.Int{big.NewInt(2), big.NewInt(11), big.NewInt(15)})
}
