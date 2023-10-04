package gas

import (
	"errors"
	"math/big"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/suite"

	"github.com/ChainSafe/sygma-core/mock"
)

type StaticGasPriceTestSuite struct {
	suite.Suite
	gasPricerMock *mock.MockGasPriceClient
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(StaticGasPriceTestSuite))
}

func (s *StaticGasPriceTestSuite) SetupSuite() {
	gomockController := gomock.NewController(s.T())
	s.gasPricerMock = mock.NewMockGasPriceClient(gomockController)
}

func (s *StaticGasPriceTestSuite) TestStaticGasPricerNoOpts() {
	twentyGwei := big.NewInt(20000000000)
	s.gasPricerMock.EXPECT().SuggestGasPrice(gomock.Any()).Return(twentyGwei, nil)
	gpd := NewStaticGasPriceDeterminant(s.gasPricerMock, nil)
	res, err := gpd.GasPrice(nil)
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
	res, err := gpd.GasPrice(nil)
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
	res, err := gpd.GasPrice(nil)
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
	res, err := gpd.GasPrice(nil)
	s.Nil(err)
	s.Equal(len(res), 1)
	s.Equal(res[0].Cmp(big.NewInt(1)), 0)
}

func (s *StaticGasPriceTestSuite) TestStaticGasPricerErrOnSuggest() {
	twentyGwei := big.NewInt(20000000000)
	gpd := NewStaticGasPriceDeterminant(s.gasPricerMock, nil)
	e := errors.New("err on suggest")
	s.gasPricerMock.EXPECT().SuggestGasPrice(gomock.Any()).Return(twentyGwei, e)
	res, err := gpd.GasPrice(nil)
	s.NotNil(err)
	s.Nil(res)
}
