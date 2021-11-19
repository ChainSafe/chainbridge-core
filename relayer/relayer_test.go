package relayer

import (
	"fmt"
	"testing"

	"github.com/ChainSafe/chainbridge-core/relayer/message"
	mock_relayer "github.com/ChainSafe/chainbridge-core/relayer/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type RouteTestSuite struct {
	suite.Suite
	mockRelayedChain *mock_relayer.MockRelayedChain
	mockMetrics      *mock_relayer.MockMetrics
}

func TestRunRouteTestSuite(t *testing.T) {
	suite.Run(t, new(RouteTestSuite))
}

func (s *RouteTestSuite) SetupSuite()    {}
func (s *RouteTestSuite) TearDownSuite() {}
func (s *RouteTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.mockRelayedChain = mock_relayer.NewMockRelayedChain(gomockController)
	s.mockMetrics = mock_relayer.NewMockMetrics(gomockController)
	s.mockMetrics.EXPECT().TrackDepositMessage(gomock.Any())
}
func (s *RouteTestSuite) TearDownTest() {}

func (s *RouteTestSuite) TestLogsErrorIfDestinationDoesNotExist() {
	relayer := Relayer{
		metrics: s.mockMetrics,
	}

	relayer.route(&message.Message{})
}

func (s *RouteTestSuite) TestLogsErrorIfMessageProcessorReturnsError() {
	s.mockRelayedChain.EXPECT().DomainID().Return(uint8(1))
	relayer := NewRelayer(
		[]RelayedChain{},
		s.mockMetrics,
		func(m *message.Message) error { return fmt.Errorf("error") },
	)
	relayer.addRelayedChain(s.mockRelayedChain)

	relayer.route(&message.Message{
		Destination: 1,
	})
}

func (s *RouteTestSuite) TestLogsErrorIfWriteReturnsError() {
	s.mockRelayedChain.EXPECT().DomainID().Return(uint8(1))
	s.mockRelayedChain.EXPECT().Write(gomock.Any()).Return(fmt.Errorf("Error"))
	relayer := NewRelayer(
		[]RelayedChain{},
		s.mockMetrics,
		func(m *message.Message) error { return nil },
	)
	relayer.addRelayedChain(s.mockRelayedChain)

	relayer.route(&message.Message{
		Destination: 1,
	})
}

func (s *RouteTestSuite) TestWritesToDestChainIfMessageValid() {
	s.mockRelayedChain.EXPECT().DomainID().Return(uint8(1))
	s.mockRelayedChain.EXPECT().Write(gomock.Any()).Return(nil)
	relayer := NewRelayer(
		[]RelayedChain{},
		s.mockMetrics,
		func(m *message.Message) error { return nil },
	)
	relayer.addRelayedChain(s.mockRelayedChain)

	relayer.route(&message.Message{
		Destination: 1,
	})
}
