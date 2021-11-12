package calls_test

import (
	"errors"
	"math/big"
	"testing"

	calls "github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	mock_calls "github.com/ChainSafe/chainbridge-core/chains/evm/calls/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type GetSolidityFunctionSigTestSuite struct {
	suite.Suite
	gomockController *gomock.Controller
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

type UtilsTestSuite struct {
	suite.Suite
	mockClientDispatcher *mock_calls.MockClientDispatcher
	mockgasPricer        *mock_calls.MockGasPricer
}

func TestRunUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (s *UtilsTestSuite) SetupSuite() {
	gomockController := gomock.NewController(s.T())
	s.mockClientDispatcher = mock_calls.NewMockClientDispatcher(gomockController)
	s.mockgasPricer = mock_calls.NewMockGasPricer(gomockController)
}
func (s *UtilsTestSuite) TearDownSuite() {}
func (s *UtilsTestSuite) SetupTest()     {}
func (s *UtilsTestSuite) TearDownTest()  {}

func (s *UtilsTestSuite) TestToCallArg() {
	kp, err := secp256k1.GenerateKeypair()

	s.Nil(err)
	address := common.HexToAddress(kp.Address())

	msg := ethereum.CallMsg{
		From:     common.Address{},
		To:       &address,
		Value:    big.NewInt(1),
		Gas:      uint64(21000),
		GasPrice: big.NewInt(3000),
		Data:     []byte("test"),
	}
	got := calls.ToCallArg(msg)
	want := map[string]interface{}{
		"from":     msg.From,
		"to":       msg.To,
		"value":    (*hexutil.Big)(msg.Value),
		"gas":      hexutil.Uint64(msg.Gas),
		"gasPrice": (*hexutil.Big)(msg.GasPrice),
		"data":     hexutil.Bytes(msg.Data),
	}
	s.Equal(want, got)
}

func (s *UtilsTestSuite) TestToCallArgWithEmptyMessage() {
	msg := ethereum.CallMsg{}
	got := calls.ToCallArg(msg)
	want := map[string]interface{}{
		"from": common.HexToAddress(""),
		"to":   (*common.Address)(nil),
	}
	s.Equal(want, got)
}

func (s *UtilsTestSuite) TestTransactNonceUnlockCallWithErrorThrown() {
	s.mockClientDispatcher.EXPECT().LockNonce().Times(1)
	s.mockClientDispatcher.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockgasPricer.EXPECT().GasPrice().Return([]*big.Int{big.NewInt(10)}, nil)
	s.mockClientDispatcher.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Times(1).Return(common.Hash{}, errors.New("error"))
	s.mockClientDispatcher.EXPECT().UnlockNonce().Times(1)
	s.mockClientDispatcher.EXPECT().WaitAndReturnTxReceipt(gomock.Any()).Times(0)
	s.mockClientDispatcher.EXPECT().UnsafeIncreaseNonce().Times(0)

	toAddress := common.HexToAddress("0xtest1")
	gasLimit := uint64(250000)
	amount := big.NewInt(10)

	_, _ = calls.Transact(
		s.mockClientDispatcher,
		evmtransaction.NewTransaction,
		s.mockgasPricer,
		&toAddress,
		[]byte("test"),
		gasLimit,
		amount)
}

func (s *UtilsTestSuite) TestTransactNonceUnlockCallWithoutErrorsThrown() {
	s.mockClientDispatcher.EXPECT().LockNonce().Times(1)
	s.mockClientDispatcher.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockgasPricer.EXPECT().GasPrice().Return([]*big.Int{big.NewInt(10)}, nil)
	s.mockClientDispatcher.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Times(1).Return(common.Hash{}, nil)
	s.mockClientDispatcher.EXPECT().From().Times(1).Return(common.Address{})
	s.mockClientDispatcher.EXPECT().WaitAndReturnTxReceipt(gomock.Any()).Times(1).Return(nil, nil)
	s.mockClientDispatcher.EXPECT().UnsafeIncreaseNonce().Times(1)
	s.mockClientDispatcher.EXPECT().UnlockNonce().Times(1)

	toAddress := common.HexToAddress("0xtest1")
	gasLimit := uint64(250000)
	amount := big.NewInt(10)

	_, _ = calls.Transact(
		s.mockClientDispatcher,
		evmtransaction.NewTransaction,
		s.mockgasPricer,
		&toAddress,
		[]byte("test"),
		gasLimit,
		amount)
}
