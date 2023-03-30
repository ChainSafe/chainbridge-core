package opentelemetry

import (
	"context"
	"net/url"
	"time"

	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
)

type OpenTelemetry struct {
	metrics          *ChainbridgeMetrics
	messageEventTime map[string]time.Time
}

// NewOpenTelemetry initializes OpenTelementry metrics
func NewOpenTelemetry(collectorRawURL string) (*OpenTelemetry, error) {
	collectorURL, err := url.Parse(collectorRawURL)
	if err != nil {
		return &OpenTelemetry{}, err
	}

	metricOptions := []otlpmetrichttp.Option{
		otlpmetrichttp.WithURLPath(collectorURL.Path),
		otlpmetrichttp.WithEndpoint(collectorURL.Host),
	}
	if collectorURL.Scheme == "http" {
		metricOptions = append(metricOptions, otlpmetrichttp.WithInsecure())
	}

	metrics, err := InitOpenTelemetryMetrics(metricOptions...)
	if err != nil {
		return &OpenTelemetry{}, err
	}

	return &OpenTelemetry{
		metrics:          metrics,
		messageEventTime: make(map[string]time.Time),
	}, nil
}

// TrackDepositMessage extracts metrics from deposit message and sends
// them to OpenTelemetry collector
func (t *OpenTelemetry) TrackDepositMessage(m *message.Message) {
	t.metrics.DepositEventCount.Add(context.Background(), 1)
	t.messageEventTime[m.ID()] = time.Now()
}

func (t *OpenTelemetry) TrackExecutionError(m *message.Message) {
	t.metrics.DepositErrorRate.Add(context.Background(), 1)
	delete(t.messageEventTime, m.ID())
}

func (t *OpenTelemetry) TrackSuccessfulExecution(m *message.Message) {
	executionLatency := time.Since(t.messageEventTime[m.ID()])
	t.metrics.ExecutionLatency.Record(context.Background(), executionLatency.Milliseconds())
	delete(t.messageEventTime, m.ID())
}

// ConsoleTelemetry is telemetry that logs metrics and should be used
// when metrics sending to OpenTelemetry should be disabled
type ConsoleTelemetry struct{}

func (t *ConsoleTelemetry) TrackDepositMessage(m *message.Message) {
	log.Info().Msgf("Deposit message: %+v", m)
}

func (t *ConsoleTelemetry) TrackExecutionError(m *message.Message) {
	log.Info().Msgf("Execution error: %+v", m)
}

func (t *ConsoleTelemetry) TrackSuccessfulExecution(m *message.Message) {
	log.Error().Msgf("Successful execution: %+v", m)
}
