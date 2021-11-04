package opentelemetry

import (
	"net/url"

	"github.com/ChainSafe/chainbridge-core/config/relayer"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
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

	metricOptions := []otlpmetrichttp.Option{
		otlpmetrichttp.WithURLPath(collectorURL.Path),
		otlpmetrichttp.WithEndpoint(collectorURL.Host),
	}
	tracerOptions := []otlptracehttp.Option{
		otlptracehttp.WithURLPath(collectorURL.Path),
		otlptracehttp.WithEndpoint(collectorURL.Host),
	}
	if collectorURL.Scheme == "http" {
		metricOptions = append(metricOptions, otlpmetrichttp.WithInsecure())
		tracerOptions = append(tracerOptions, otlptracehttp.WithInsecure())
	}

	metrics, err := initOpenTelemetryMetrics(metricOptions...)
	if err != nil {
		return &telemetry{}, err
	}

	tracer, err := initOpenTelementryTracer(tracerOptions...)
	if err != nil {
		return &telemetry{}, err
	}

	return &telemetry{
		metrics: metrics,
		tracer:  tracer,
	}, nil
}

func (t *telemetry) TraceDepositEvent(m *message.Message) {}
