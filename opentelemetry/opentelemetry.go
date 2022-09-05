package opentelemetry

import (
	"context"
	"net/url"

	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
)

type OpenTelemetry struct {
	metrics *ChainbridgeMetrics
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

	metrics, err := initOpenTelemetryMetrics(metricOptions...)
	if err != nil {
		return &OpenTelemetry{}, err
	}

	return &OpenTelemetry{
		metrics: metrics,
	}, nil
}

// TrackDepositMessage extracts metrics from deposit message and sends
// them to OpenTelemetry collector
func (t *OpenTelemetry) TrackDepositMessage(m *message.Message) {
	t.metrics.DepositEventCount.Add(context.Background(), 1)
}

// ConsoleTelemetry is telemetry that logs metrics and should be used
// when metrics sending to OpenTelemetry should be disabled
type ConsoleTelemetry struct{}

func (t *ConsoleTelemetry) TrackDepositMessage(m *message.Message) {
	log.Info().Msgf("Deposit message: %+v", m)
}
