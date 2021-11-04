package relayer

import (
	"context"
	"fmt"
	"testing"

	"github.com/ChainSafe/chainbridge-core/relayer/message"
	mock_relayer "github.com/ChainSafe/chainbridge-core/relayer/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type RouteTestSuite struct {
	suite.Suite
	mockTracer       *mock_relayer.MockTracer
	mockRelayedChain *mock_relayer.MockRelayedChain
}

func TestRunRouteTestSuite(t *testing.T) {
	suite.Run(t, new(RouteTestSuite))
}

func (s *RouteTestSuite) SetupSuite()    {}
func (s *RouteTestSuite) TearDownSuite() {}
func (s *RouteTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.mockTracer = mock_relayer.NewMockTracer(gomockController)
	s.mockRelayedChain = mock_relayer.NewMockRelayedChain(gomockController)
}
func (s *RouteTestSuite) TearDownTest() {}

func (s *RouteTestSuite) TestLogsErrorIfDestinationDoesNotExist() {
	s.mockTracer.EXPECT().TraceDepositEvent(gomock.Any(), gomock.Any()).Return(context.Background())
	relayer := Relayer{}

	relayer.route(&message.Message{}, s.mockTracer)
}

func (s *RouteTestSuite) TestLogsErrorIfMessageProcessorReturnsError() {
	s.mockTracer.EXPECT().TraceDepositEvent(gomock.Any(), gomock.Any()).Return(context.Background())
	s.mockRelayedChain.EXPECT().DomainID().Return(uint8(1))
	relayer := Relayer{
		messageProcessors: []message.MessageProcessor{
			func(m *message.Message) error { return fmt.Errorf("error") },
		},
	}
	relayer.addRelayedChain(s.mockRelayedChain)

	relayer.route(&message.Message{
		Destination: 1,
	}, s.mockTracer)
}

func (s *RouteTestSuite) TestLogsErrorIfWriteReturnsError() {
	s.mockTracer.EXPECT().TraceDepositEvent(gomock.Any(), gomock.Any()).Return(context.Background())
	s.mockRelayedChain.EXPECT().DomainID().Return(uint8(1))
	s.mockRelayedChain.EXPECT().Write(gomock.Any()).Return(fmt.Errorf("Error"))
	relayer := Relayer{
		messageProcessors: []message.MessageProcessor{
			func(m *message.Message) error { return nil },
		},
	}
	relayer.addRelayedChain(s.mockRelayedChain)

	relayer.route(&message.Message{
		Destination: 1,
	}, s.mockTracer)
}

func (s *RouteTestSuite) TestWritesToDestChainIfMessageValid() {
	s.mockTracer.EXPECT().TraceDepositEvent(gomock.Any(), gomock.Any()).Return(context.Background())
	s.mockRelayedChain.EXPECT().DomainID().Return(uint8(1))
	s.mockRelayedChain.EXPECT().Write(gomock.Any()).Return(nil)
	relayer := Relayer{
		messageProcessors: []message.MessageProcessor{
			func(m *message.Message) error { return nil },
		},
	}
	relayer.addRelayedChain(s.mockRelayedChain)

	relayer.route(&message.Message{
		Destination: 1,
	}, s.mockTracer)
}
