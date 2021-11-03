package opentelemetry

import (
	"github.com/ChainSafe/chainbridge-core/config/relayer"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	tracer "go.opentelemetry.io/otel/trace"
)

type telemetry struct {
	metrics *ChainbridgeMetrics
	tracer  tracer.Tracer
}

func NewOpenTelemetry(config relayer.RelayerConfig) (*Telemetry, error) {
	metrics, err := initOpenTelemetryMetrics(otlpmetrichttp.WithEndpoint("localhost:4318"), otlpmetrichttp.WithInsecure())
	if err != nil {
		return &Telemetry{}, err
	}

	tracer, err := initOpenTelementryTracer(otlptracehttp.WithEndpoint("localhost:4318"), otlptracehttp.WithInsecure())
	if err != nil {
		return &Telemetry{}, err
	}

	return &Telemetry{
		metrics: metrics,
		tracer:  tracer,
	}, nil
}
