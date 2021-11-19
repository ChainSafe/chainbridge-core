package relayer

import (
	"fmt"
	"math/big"
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
}
func (s *RouteTestSuite) TearDownTest() {}

func (s *RouteTestSuite) TestLogsErrorIfDestinationDoesNotExist() {
	s.mockMetrics.EXPECT().TrackDepositMessage(gomock.Any())
	relayer := Relayer{
		metrics: s.mockMetrics,
	}

	relayer.route(&message.Message{})
}

// TestRouter tests relayers router
func (s *RouteTestSuite) TestAdjustDecimalsForERC20AmountMessageProcessor() {
	a, _ := big.NewInt(0).SetString("145556700000000000000", 10) // 145.5567 tokens
	msg := &message.Message{
		Destination: 2,
		Source:      1,
		Payload: []interface{}{
			a.Bytes(), // 145.5567 tokens
		},
	}
	err := message.AdjustDecimalsForERC20AmountMessageProcessor(map[uint8]uint64{1: 18, 2: 2})(msg)
	s.Nil(err)
	amount := new(big.Int).SetBytes(msg.Payload[0].([]byte))
	if amount.Cmp(big.NewInt(14555)) != 0 {
		s.Fail("wrong amount")
	}

}

func (s *RouteTestSuite) TestLogsErrorIfMessageProcessorReturnsError() {
	s.mockMetrics.EXPECT().TrackDepositMessage(gomock.Any())
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
	s.mockMetrics.EXPECT().TrackDepositMessage(gomock.Any())
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
	s.mockMetrics.EXPECT().TrackDepositMessage(gomock.Any())
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
