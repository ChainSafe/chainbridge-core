package evmgaspricer

import (
	"math/big"
	"testing"

	mock_evmgaspricer "github.com/ChainSafe/sygma-core/chains/evm/calls/evmgaspricer/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type LondonGasPriceTestSuite struct {
	suite.Suite
	gasPricerMock *mock_evmgaspricer.MockLondonGasClient
}

func TestRuLondonTestSuite(t *testing.T) {
	suite.Run(t, new(LondonGasPriceTestSuite))
}

func (s *LondonGasPriceTestSuite) SetupSuite()    {}
func (s *LondonGasPriceTestSuite) TearDownSuite() {}
func (s *LondonGasPriceTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.gasPricerMock = mock_evmgaspricer.NewMockLondonGasClient(gomockController)
}
func (s *LondonGasPriceTestSuite) TearDownTest() {}

func (s *LondonGasPriceTestSuite) TestLondonGasPricerNoOpts() {
	twentyGwei := big.NewInt(20000000000)
	twoGwei := big.NewInt(2000000000)
	gpd := NewLondonGasPriceClient(s.gasPricerMock, nil)
	s.gasPricerMock.EXPECT().BaseFee().Return(twentyGwei, nil)
	s.gasPricerMock.EXPECT().SuggestGasTipCap(gomock.Any()).Return(twoGwei, nil)

	res, err := gpd.GasPrice(nil)
	s.Nil(err)
	s.Equal(len(res), 2)
	s.Equal(res[0].Cmp(twoGwei), 0)
	s.Equal(0, res[1].Cmp(big.NewInt(42000000000))) // Base fee 20Gwei * 2 + maxTipCap = 42Gwei
}

func (s *LondonGasPriceTestSuite) TestLondonGasPricerWithUpperLimit() {
	twentyGwei := big.NewInt(20000000000)
	thirtyGwei := big.NewInt(30000000000)
	twoGwei := big.NewInt(2000000000)
	gpd := NewLondonGasPriceClient(s.gasPricerMock, &GasPricerOpts{UpperLimitFeePerGas: thirtyGwei})
	s.gasPricerMock.EXPECT().BaseFee().Return(twentyGwei, nil)
	s.gasPricerMock.EXPECT().SuggestGasTipCap(gomock.Any()).Return(twoGwei, nil)

	res, err := gpd.GasPrice(nil)
	s.Nil(err)
	s.Equal(len(res), 2)
	s.Equal(res[0].Cmp(big.NewInt(10000000000)), 0) // 10 gwei. Bcs MaxFee is UpperLimit and BaseFee is 20Gwei so PriorityFee is 30-20=10Gwei
	s.Equal(0, res[1].Cmp(thirtyGwei))              // Equals to UpperLimit
}

func (s *LondonGasPriceTestSuite) TestLondonGasPricerWithUpperLimitLowerBaseFee() {
	twentyGwei := big.NewInt(20000000000)
	thirtyGwei := big.NewInt(30000000000)
	gpd := NewLondonGasPriceClient(s.gasPricerMock, &GasPricerOpts{UpperLimitFeePerGas: twentyGwei})
	s.gasPricerMock.EXPECT().BaseFee().Return(thirtyGwei, nil)
	//s.gasPricerMock.EXPECT().SuggestGasTipCap(gomock.Any()).Return(twoGwei, nil) // Code is not get to the point where this call happens

	res, err := gpd.GasPrice(nil)
	s.Nil(err)
	s.Equal(len(res), 2)
	s.Equal(res[0].Cmp(big.NewInt(TwoAndTheHalfGwei)), 0) // Lowest MaxPriorityFee
	s.Equal(0, res[1].Cmp(big.NewInt(32500000000)))       // Equals to BaseFee  + MaxPriorityFee (22,5 gwei)
}
