package opentelemetry

import (
	"context"
	"net/url"
	"time"

	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	api "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initResource() *sdkresource.Resource {
	res, _ := sdkresource.New(context.Background(),
		sdkresource.WithProcess(),
		sdkresource.WithTelemetrySDK(),
		sdkresource.WithHost(),
		sdkresource.WithAttributes(
			semconv.ServiceName("relayer"),
		),
	)
	return res
}

func InitMetricProvider(ctx context.Context, agentURL string) (*sdkmetric.MeterProvider, error) {
	collectorURL, err := url.Parse(agentURL)
	if err != nil {
		return nil, err
	}

	metricOptions := []otlpmetrichttp.Option{
		otlpmetrichttp.WithURLPath(collectorURL.Path),
		otlpmetrichttp.WithEndpoint(collectorURL.Host),
	}
	if collectorURL.Scheme == "http" {
		metricOptions = append(metricOptions, otlpmetrichttp.WithInsecure())
	}

	metricHTTPExporter, err := otlpmetrichttp.New(ctx, metricOptions...)
	if err != nil {
		return nil, err
	}

	httpMetricReader := sdkmetric.NewPeriodicReader(metricHTTPExporter)

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(httpMetricReader),
		sdkmetric.WithResource(initResource()),
	)

	return meterProvider, nil
}

type OpenTelemetry struct {
	metrics          *ChainbridgeMetrics
	messageEventTime map[string]time.Time
	Opts             api.MeasurementOption
}

// NewOpenTelemetry initializes OpenTelementry metrics
func NewOpenTelemetry(meter metric.Meter, env, relayerID string) (*OpenTelemetry, error) {
	opts := api.WithAttributes(attribute.String("relayerid", relayerID), attribute.String("env", env))
	metrics, err := NewChainbridgeMetrics(meter)
	if err != nil {
		return nil, err
	}
	return &OpenTelemetry{
		metrics:          metrics,
		messageEventTime: make(map[string]time.Time),
		Opts:             opts,
	}, err
}

// TrackDepositMessage extracts metrics from deposit message and sends
// them to OpenTelemetry collector
func (t *OpenTelemetry) TrackDepositMessage(m *message.Message) {
	t.metrics.DepositEventCount.Add(context.Background(), 1, t.Opts, api.WithAttributes(attribute.Int64("source", int64(m.Source))))
	t.messageEventTime[m.ID()] = time.Now()
}

func (t *OpenTelemetry) TrackExecutionError(m *message.Message) {
	t.metrics.ExecutionErrorCount.Add(context.Background(), 1, t.Opts, api.WithAttributes(attribute.Int64("destination", int64(m.Source))))
	delete(t.messageEventTime, m.ID())
}

func (t *OpenTelemetry) TrackSuccessfulExecutionLatency(m *message.Message) {
	executionLatency := time.Since(t.messageEventTime[m.ID()]).Milliseconds() / 1000
	t.metrics.ExecutionLatency.Record(context.Background(), executionLatency)
	t.metrics.ExecutionLatencyPerRoute.Record(
		context.Background(),
		executionLatency,
		t.Opts,
		api.WithAttributes(attribute.Int64("source", int64(m.Source))),
		api.WithAttributes(attribute.Int64("destination", int64(m.Destination))))
	delete(t.messageEventTime, m.ID())
}
