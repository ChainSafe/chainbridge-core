package opentelemetry

import (
	"context"
	"encoding/hex"
	"net/url"

	"github.com/ChainSafe/chainbridge-core/config/relayer"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/trace"
	tracer "go.opentelemetry.io/otel/trace"
)

type Telemetry struct {
	metrics *ChainbridgeMetrics
	tracer  tracer.Tracer
}

func NewOpenTelemetry(config relayer.RelayerConfig) (*Telemetry, error) {
	collectorURL, err := url.Parse(config.OpenTelemetryCollectorURL)
	if err != nil {
		return &Telemetry{}, err
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
		return &Telemetry{}, err
	}

	tracer, err := initOpenTelementryTracer(tracerOptions...)
	if err != nil {
		return &Telemetry{}, err
	}

	return &Telemetry{
		metrics: metrics,
		tracer:  tracer,
	}, nil
}

func (t *Telemetry) TraceDepositEvent(ctx context.Context, m *message.Message) context.Context {
	if t.tracer != nil {
		var span tracer.Span
		ctx, span = t.tracer.Start(
			context.Background(),
			"deposit-event",
			trace.WithAttributes(
				attribute.String("Type", string(m.Type)),
				attribute.Int("Source", int(m.Source)),
				attribute.Int("Destination", int(m.Destination)),
				attribute.Int("DepositNonce", int(m.DepositNonce)),
				attribute.String("ResourceId", hex.EncodeToString(m.ResourceId[:])),
			),
		)
		defer span.End()
	}

	if t.metrics != nil {
		t.metrics.DepositEventCount.Add(ctx, 1)
	}

	return ctx
}
