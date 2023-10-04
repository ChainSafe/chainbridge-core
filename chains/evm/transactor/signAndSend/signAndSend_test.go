package signAndSend_test

import (
	"math/big"
	"testing"

	"github.com/ChainSafe/sygma-core/chains/evm/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/transactor/signAndSend"
	"github.com/ChainSafe/sygma-core/chains/evm/transactor/transaction"
	"github.com/ChainSafe/sygma-core/mock"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type TransactorTestSuite struct {
	suite.Suite
	gomockController *gomock.Controller
	mockClient       *mock.MockClient
	mockTransactor   *mock.MockTransactor
	mockGasPricer    *mock.MockGasPricer
}

func TestSignAndSendTestSuite(t *testing.T) {
	suite.Run(t, new(TransactorTestSuite))
}

func (s *TransactorTestSuite) SetupSuite()    {}
func (s *TransactorTestSuite) TearDownSuite() {}
func (s *TransactorTestSuite) SetupTest() {
	s.gomockController = gomock.NewController(s.T())
	s.mockClient = mock.NewMockClient(s.gomockController)
	s.mockTransactor = mock.NewMockTransactor(s.gomockController)
	s.mockGasPricer = mock.NewMockGasPricer(s.gomockController)
}

func (s *TransactorTestSuite) TestTransactor_SignAndSend_Success() {
	var byteData = []byte{47, 47, 241, 93, 159, 45, 240, 254, 210, 199, 118, 72, 222, 88, 96, 164, 204, 80, 140, 208, 129, 140, 133, 184, 184, 161, 171, 76, 238, 239, 141, 152, 28, 137, 86, 166, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 60, 48, 181, 109, 237, 4, 127, 230, 34, 95, 112, 4, 234, 75, 225, 174, 112, 201, 2, 106}

	s.mockClient.EXPECT().LockNonce()
	s.mockClient.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockGasPricer.EXPECT().GasPrice(gomock.Any()).Return([]*big.Int{big.NewInt(1)}, nil)
	s.mockClient.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Return(common.Hash{1, 2, 3, 4, 5}, nil)
	s.mockClient.EXPECT().WaitAndReturnTxReceipt(gomock.Any()).Return(&types.Receipt{}, nil)
	s.mockClient.EXPECT().UnsafeIncreaseNonce().Return(nil)
	s.mockClient.EXPECT().UnlockNonce()

	txFabric := transaction.NewTransaction
	var trans = signAndSend.NewSignAndSendTransactor(
		txFabric,
		s.mockGasPricer,
		s.mockClient,
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
