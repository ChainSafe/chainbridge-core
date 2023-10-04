package transaction_test

import (
	"math/big"
	"testing"

	"github.com/ChainSafe/sygma-core/chains/evm/transactor/gas"
	"github.com/ChainSafe/sygma-core/chains/evm/transactor/transaction"
	"github.com/ChainSafe/sygma-core/crypto/keystore"
	"github.com/ChainSafe/sygma-core/mock"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var aliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]

type EVMTxTestSuite struct {
	suite.Suite
	client *mock.MockLondonGasClient
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(EVMTxTestSuite))
}

func (s *EVMTxTestSuite) SetupSuite()    {}
func (s *EVMTxTestSuite) TearDownSuite() {}
func (s *EVMTxTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.client = mock.NewMockLondonGasClient(gomockController)
}
func (s *EVMTxTestSuite) TearDownTest() {}

func (s *EVMTxTestSuite) TestNewTransactionWithStaticGasPricer() {
	s.client.EXPECT().SuggestGasPrice(gomock.Any()).Return(big.NewInt(1000), nil)
	txFabric := transaction.NewTransaction
	gasPriceClient := gas.NewStaticGasPriceDeterminant(s.client, nil)
	gp, err := gasPriceClient.GasPrice(nil)
	s.Nil(err)
	tx, err := txFabric(1, &common.Address{}, big.NewInt(0), 10000, gp, []byte{})
	s.Nil(err)
	rawTx, err := tx.RawWithSignature(aliceKp, big.NewInt(420))
	s.Nil(err)
	txt := types.Transaction{}
	err = txt.UnmarshalBinary(rawTx)
	s.Nil(err)
	s.Equal(types.LegacyTxType, int(txt.Type()))
}

func (s *EVMTxTestSuite) TestNewTransactionWithLondonGasPricer() {
	s.client.EXPECT().BaseFee().Return(big.NewInt(1000), nil)
	s.client.EXPECT().SuggestGasTipCap(gomock.Any()).Return(big.NewInt(1000), nil)
	txFabric := transaction.NewTransaction
	gasPriceClient := gas.NewLondonGasPriceClient(s.client, nil)
	gp, err := gasPriceClient.GasPrice(nil)
	s.Nil(err)
	tx, err := txFabric(1, &common.Address{}, big.NewInt(0), 10000, gp, []byte{})
	s.Nil(err)
	rawTx, err := tx.RawWithSignature(aliceKp, big.NewInt(420))
	s.Nil(err)
	txt := types.Transaction{}
	err = txt.UnmarshalBinary(rawTx)
	s.Nil(err)
	s.Equal(types.DynamicFeeTxType, int(txt.Type()))
}
