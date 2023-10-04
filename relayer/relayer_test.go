package relayer

import (
	"fmt"
	"testing"

	"github.com/ChainSafe/sygma-core/mock"
	"github.com/ChainSafe/sygma-core/types"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type RouteTestSuite struct {
	suite.Suite
	mockRelayedChain *mock.MockRelayedChain
	mockMetrics      *mock.MockDepositMeter
}

func TestRunRouteTestSuite(t *testing.T) {
	suite.Run(t, new(RouteTestSuite))
}

func (s *RouteTestSuite) SetupSuite()    {}
func (s *RouteTestSuite) TearDownSuite() {}
func (s *RouteTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.mockRelayedChain = mock.NewMockRelayedChain(gomockController)
	s.mockMetrics = mock.NewMockDepositMeter(gomockController)
}
func (s *RouteTestSuite) TearDownTest() {}

func (s *RouteTestSuite) TestLogsErrorIfDestinationDoesNotExist() {
	relayer := Relayer{
		metrics: s.mockMetrics,
	}

	relayer.route([]*types.Message{
		{},
	})
}

func (s *RouteTestSuite) TestWriteFail() {
	s.mockMetrics.EXPECT().TrackExecutionError(gomock.Any())
	s.mockRelayedChain.EXPECT().DomainID().Return(uint8(1)).Times(3)
	s.mockRelayedChain.EXPECT().Write(gomock.Any()).Return(fmt.Errorf("error"))
	relayer := NewRelayer(
		[]RelayedChain{},
		s.mockMetrics,
	)
	relayer.addRelayedChain(s.mockRelayedChain)

	relayer.route([]*types.Message{
		{Destination: 1},
	})
}

func (s *RouteTestSuite) TestWritesToDestChainIfMessageValid() {
	s.mockMetrics.EXPECT().TrackSuccessfulExecutionLatency(gomock.Any())
	s.mockRelayedChain.EXPECT().DomainID().Return(uint8(1)).Times(2)
	s.mockRelayedChain.EXPECT().Write(gomock.Any())
	relayer := NewRelayer(
		[]RelayedChain{},
		s.mockMetrics,
	)
	relayer.addRelayedChain(s.mockRelayedChain)

	relayer.route([]*types.Message{
		{Destination: 1},
	})
}
