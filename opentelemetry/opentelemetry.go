package opentelemetry

import (
	"context"
	"net/url"

	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
)

type Telemetry struct {
	metrics *ChainbridgeMetrics
}

// NewOpenTelemetry initializes OpenTelementry metrics
func NewOpenTelemetry(collectorRawURL string) (*Telemetry, error) {
	collectorURL, err := url.Parse(collectorRawURL)
	if err != nil {
		return &Telemetry{}, err
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
		return &Telemetry{}, err
	}

	return &Telemetry{
		metrics: metrics,
	}, nil
}

// TrackDepositMessage extracts metrics from deposit message and sends
// them to OpenTelemetry collector
func (t *Telemetry) TrackDepositMessage(m *message.Message) {
	t.metrics.DepositEventCount.Add(context.Background(), 1)
}
