package evmgaspricer

import (
	"errors"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"

	mock_evmgaspricer "github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer/mock"
	"github.com/stretchr/testify/suite"
)

type StaticGasPriceTestSuite struct {
	suite.Suite
	gasPricerMock *mock_evmgaspricer.MockGasPriceClient
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(StaticGasPriceTestSuite))
}

func (s *StaticGasPriceTestSuite) SetupSuite()    {}
func (s *StaticGasPriceTestSuite) TearDownSuite() {}
func (s *StaticGasPriceTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.gasPricerMock = mock_evmgaspricer.NewMockGasPriceClient(gomockController)
}
func (s *StaticGasPriceTestSuite) TearDownTest() {}

func (s *StaticGasPriceTestSuite) TestStaticGasPricerNoOpts() {
	twentyGwei := big.NewInt(20000000000)
	gpd := NewStaticGasPriceDeterminant(s.gasPricerMock, nil)
	s.gasPricerMock.EXPECT().SuggestGasPrice(gomock.Any()).Return(twentyGwei, nil)
	res, err := gpd.GasPrice()
	s.Nil(err)
	s.Equal(len(res), 1)
	s.Equal(res[0].Cmp(twentyGwei), 0)
}

func (s *StaticGasPriceTestSuite) TestStaticGasPricerFactorSet() {
	twentyGwei := big.NewInt(20000000000)
	gpd := NewStaticGasPriceDeterminant(s.gasPricerMock, &GasPricerOpts{
		GasPriceFactor: big.NewFloat(2.5),
	})
	s.gasPricerMock.EXPECT().SuggestGasPrice(gomock.Any()).Return(twentyGwei, nil)
	res, err := gpd.GasPrice()
	s.Nil(err)
	s.Equal(len(res), 1)
	s.Equal(res[0].Cmp(big.NewInt(50000000000)), 0)
}

func (s *StaticGasPriceTestSuite) TestStaticGasPricerUpperLimitSet() {
	twentyGwei := big.NewInt(20000000000)
	gpd := NewStaticGasPriceDeterminant(s.gasPricerMock, &GasPricerOpts{
		UpperLimitFeePerGas: big.NewInt(1),
	})
	s.gasPricerMock.EXPECT().SuggestGasPrice(gomock.Any()).Return(twentyGwei, nil)
	res, err := gpd.GasPrice()
	s.Nil(err)
	s.Equal(len(res), 1)
	s.Equal(res[0].Cmp(big.NewInt(1)), 0)
}

func (s *StaticGasPriceTestSuite) TestStaticGasPricerUpperLimitAndFactorSet() {
	twentyGwei := big.NewInt(20000000000)
	gpd := NewStaticGasPriceDeterminant(s.gasPricerMock, &GasPricerOpts{
		UpperLimitFeePerGas: big.NewInt(1),
		GasPriceFactor:      big.NewFloat(2.5),
	})
	s.gasPricerMock.EXPECT().SuggestGasPrice(gomock.Any()).Return(twentyGwei, nil)
	res, err := gpd.GasPrice()
	s.Nil(err)
	s.Equal(len(res), 1)
	s.Equal(res[0].Cmp(big.NewInt(1)), 0)
}

func (s *StaticGasPriceTestSuite) TestStaticGasPricerErrOnSuggest() {
	twentyGwei := big.NewInt(20000000000)
	gpd := NewStaticGasPriceDeterminant(s.gasPricerMock, nil)
	e := errors.New("err on suggest")
	s.gasPricerMock.EXPECT().SuggestGasPrice(gomock.Any()).Return(twentyGwei, e)
	res, err := gpd.GasPrice()
	s.NotNil(err)
	s.Nil(res)
}
