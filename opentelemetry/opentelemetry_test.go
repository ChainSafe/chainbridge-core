package opentelemetry

import (
	"context"
	"testing"

	mock_opentelementry "github.com/ChainSafe/chainbridge-core/opentelemetry/mock"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/metric"
)

type TraceDepositEventTestSuite struct {
	suite.Suite
	mockTracer    *mock_opentelementry.MockTracer
	mockSpan      *mock_opentelementry.MockSpan
	mockMeterImpl *mock_opentelementry.MockMeterImpl
}

func TestRunTraceDepositEventTestSuite(t *testing.T) {
	suite.Run(t, new(TraceDepositEventTestSuite))
}

func (s *TraceDepositEventTestSuite) SetupSuite()    {}
func (s *TraceDepositEventTestSuite) TearDownSuite() {}
func (s *TraceDepositEventTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.mockTracer = mock_opentelementry.NewMockTracer(gomockController)
	s.mockSpan = mock_opentelementry.NewMockSpan(gomockController)
	s.mockMeterImpl = mock_opentelementry.NewMockMeterImpl(gomockController)
}
func (s *TraceDepositEventTestSuite) TearDownTest() {}

func (s *TraceDepositEventTestSuite) TestDoesNothingIfTelemetryNil() {
	telemetry := Telemetry{}

	ctx := telemetry.TraceDepositEvent(context.Background(), &message.Message{})

	s.Equal(context.Background(), ctx)
}

func (s *TraceDepositEventTestSuite) TestTracesMessageIfTracerExists() {
	s.mockSpan.EXPECT().End().Return()
	s.mockTracer.EXPECT().Start(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(context.Background(), s.mockSpan)
	telemetry := Telemetry{
		tracer: s.mockTracer,
	}

	ctx := telemetry.TraceDepositEvent(context.Background(), &message.Message{})

	s.Equal(context.Background(), ctx)
}

func (s *TraceDepositEventTestSuite) TestIncreasesDepositEventCountIfMetricsExist() {
	metrics := newChainbridgeMetrics(metric.Meter{})
	telemetry := Telemetry{
		metrics: metrics,
	}

	ctx := telemetry.TraceDepositEvent(context.Background(), &message.Message{})

	s.Equal(context.Background(), ctx)
}
