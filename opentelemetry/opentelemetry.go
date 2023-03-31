package opentelemetry

import (
	"context"
	"net/url"
	"time"

	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

func DefaultMeter(ctx context.Context, collectorRawURL string) (metric.Meter, error) {
	collectorURL, err := url.Parse(collectorRawURL)
	if err != nil {
		return metric.Meter{}, err
	}

	metricOptions := []otlpmetrichttp.Option{
		otlpmetrichttp.WithURLPath(collectorURL.Path),
		otlpmetrichttp.WithEndpoint(collectorURL.Host),
	}
	if collectorURL.Scheme == "http" {
		metricOptions = append(metricOptions, otlpmetrichttp.WithInsecure())
	}
	client := otlpmetrichttp.NewClient(metricOptions...)
	exp, err := otlpmetric.New(ctx, client)
	if err != nil {
		return metric.Meter{}, err
	}

	selector := simple.NewWithHistogramDistribution()
	proc := processor.NewFactory(selector, export.CumulativeExportKindSelector())
	cont := controller.New(proc, controller.WithExporter(exp))
	global.SetMeterProvider(cont)

	err = cont.Start(ctx)
	if err != nil {
		return metric.Meter{}, err
	}

	return cont.Meter("chainbridge"), nil
}

type OpenTelemetry struct {
	metrics          *ChainbridgeMetrics
	messageEventTime map[string]time.Time
}

// NewOpenTelemetry initializes OpenTelementry metrics
func NewOpenTelemetry(meter metric.Meter) *OpenTelemetry {
	metrics := NewChainbridgeMetrics(meter)
	return &OpenTelemetry{
		metrics:          metrics,
		messageEventTime: make(map[string]time.Time),
	}
}

// TrackDepositMessage extracts metrics from deposit message and sends
// them to OpenTelemetry collector
func (t *OpenTelemetry) TrackDepositMessage(m *message.Message) {
	t.metrics.DepositEventCount.Add(context.Background(), 1, attribute.Int64("source", int64(m.Source)))
	t.messageEventTime[m.ID()] = time.Now()
}

func (t *OpenTelemetry) TrackExecutionError(m *message.Message) {
	t.metrics.ExecutionErrorCount.Add(context.Background(), 1, attribute.Int64("destination", int64(m.Source)))
	delete(t.messageEventTime, m.ID())
}

func (t *OpenTelemetry) TrackSuccessfulExecution(m *message.Message) {
	executionLatency := time.Since(t.messageEventTime[m.ID()])
	t.metrics.ExecutionLatency.Record(context.Background(), executionLatency.Milliseconds())
	t.metrics.ExecutionLatencyPerRoute.Record(
		context.Background(),
		executionLatency.Milliseconds(),
		attribute.Int64("source", int64(m.Source)),
		attribute.Int64("destination", int64(m.Destination)))
	delete(t.messageEventTime, m.ID())
}
