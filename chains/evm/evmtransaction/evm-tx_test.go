package evmtransaction

import (
	"math/big"
	"testing"

	mock_evmgaspricer "github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer/mock"
	"github.com/ChainSafe/chainbridge-core/e2e/evm"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ethereum/go-ethereum/common"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type EVMTxTestSuite struct {
	suite.Suite
	client *mock_evmgaspricer.MockLondonGasClient
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(EVMTxTestSuite))
}

func (s *EVMTxTestSuite) SetupSuite()    {}
func (s *EVMTxTestSuite) TearDownSuite() {}
func (s *EVMTxTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.client = mock_evmgaspricer.NewMockLondonGasClient(gomockController)
}
func (s *EVMTxTestSuite) TearDownTest() {}

func (s *EVMTxTestSuite) TestNewTransactionWithStaticGasPricer() {
	s.client.EXPECT().SuggestGasPrice(gomock.Any()).Return(big.NewInt(1000), nil)
	txFabric := NewTransaction
	gasPriceClient := evmgaspricer.NewStaticGasPriceDeterminant(s.client, nil)
	gp, err := gasPriceClient.GasPrice()
	s.Nil(err)
	tx, err := txFabric(1, &common.Address{}, big.NewInt(0), 10000, gp, []byte{})
	s.Nil(err)
	rawTx, err := tx.RawWithSignature(evm.AliceKp.PrivateKey(), big.NewInt(420))
	s.Nil(err)
	txt := types.Transaction{}
	err = txt.UnmarshalBinary(rawTx)
	s.Nil(err)
	s.Equal(types.LegacyTxType, int(txt.Type()))
}

func (s *EVMTxTestSuite) TestNewTransactionWithLondonGasPricer() {
	s.client.EXPECT().BaseFee().Return(big.NewInt(1000), nil)
	s.client.EXPECT().SuggestGasTipCap(gomock.Any()).Return(big.NewInt(1000), nil)
	txFabric := NewTransaction
	gasPriceClient := evmgaspricer.NewLondonGasPriceClient(s.client, nil)
	gp, err := gasPriceClient.GasPrice()
	s.Nil(err)
	tx, err := txFabric(1, &common.Address{}, big.NewInt(0), 10000, gp, []byte{})
	s.Nil(err)
	rawTx, err := tx.RawWithSignature(evm.AliceKp.PrivateKey(), big.NewInt(420))
	s.Nil(err)
	txt := types.Transaction{}
	err = txt.UnmarshalBinary(rawTx)
	s.Nil(err)
	s.Equal(types.DynamicFeeTxType, int(txt.Type()))
}
