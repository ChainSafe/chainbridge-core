package opentelemetry

import (
	"net/url"

	"github.com/ChainSafe/chainbridge-core/config/relayer"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	tracer "go.opentelemetry.io/otel/trace"
)

type telemetry struct {
	metrics *ChainbridgeMetrics
	tracer  tracer.Tracer
}

func NewOpenTelemetry(config relayer.RelayerConfig) (*telemetry, error) {
	collectorURL, err := url.Parse(config.OpenTelemetryCollectorURL)
	if err != nil {
		return &telemetry{}, err
	}

	options := []otlpmetrichttp.Option{
		otlpmetrichttp.WithURLPath(collectorURL.Path),
		otlpmetrichttp.WithEndpoint(collectorURL.Host),
	}
	if collectorURL.Scheme == "http" {
		options = append(options, otlpmetrichttp.WithInsecure())
	}

	metrics, err := initOpenTelemetryMetrics(options)
	if err != nil {
		return &telemetry{}, err
	}

	tracer, err := initOpenTelementryTracer(options)
	if err != nil {
		return &telemetry{}, err
	}

	return &telemetry{
		metrics: metrics,
		tracer:  tracer,
	}, nil
}

func (t *telemetry) TraceDepositEvent() {}
