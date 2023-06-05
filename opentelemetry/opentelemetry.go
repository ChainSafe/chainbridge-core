package opentelemetry

import (
	"context"
	"math/big"
	"net/url"
	"time"

	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
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

	selector := simple.NewWithHistogramDistribution(histogram.WithExplicitBoundaries([]float64{15, 60, 300, 900, 2700, 5400}))
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
	meter            metric.Meter
	messageEventTime map[string]time.Time
	genericLabels    []attribute.KeyValue
}

// NewOpenTelemetry initializes OpenTelementry metrics
func NewOpenTelemetry(meter metric.Meter, labels ...attribute.KeyValue) *OpenTelemetry {
	metrics := NewChainbridgeMetrics(meter, labels...)
	return &OpenTelemetry{
		metrics:          metrics,
		meter:            meter,
		genericLabels:    labels,
		messageEventTime: make(map[string]time.Time),
	}
}

// TrackDepositMessage extracts metrics from deposit message and sends
// them to OpenTelemetry collector
func (t *OpenTelemetry) TrackDepositMessage(m *message.Message) {
	labels := append(t.genericLabels, attribute.Int64("source", int64(m.Source)))
	t.metrics.DepositEventCount.Add(context.Background(), 1, labels...)
	t.messageEventTime[m.ID()] = time.Now()
}

func (t *OpenTelemetry) TrackExecutionError(m *message.Message) {
	labels := append(t.genericLabels, attribute.Int64("destination", int64(m.Source)))
	t.metrics.ExecutionErrorCount.Add(context.Background(), 1, labels...)
	delete(t.messageEventTime, m.ID())
}

func (t *OpenTelemetry) TrackSuccessfulExecution(m *message.Message) {
	labels := append(t.genericLabels, attribute.Int64("source", int64(m.Source)))
	labels = append(labels, attribute.Int64("destination", int64(m.Destination)))
	executionLatency := time.Since(t.messageEventTime[m.ID()]).Milliseconds() / 1000
	t.metrics.ExecutionLatency.Record(context.Background(), executionLatency)
	t.metrics.ExecutionLatencyPerRoute.Record(
		context.Background(),
		executionLatency,
		labels...,
	)
	delete(t.messageEventTime, m.ID())
}

func (t *OpenTelemetry) TrackBlockDelta(domainID uint8, head *big.Int, current *big.Int) {
	t.metrics.BlockDeltaMap[domainID] = new(big.Int).Sub(head, current)
	t.meter.RecordBatch(context.Background(), []attribute.KeyValue{})
}
