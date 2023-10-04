package observability

import (
	"context"
	"math/big"
	"net/url"
	"sync"
	"time"

	"github.com/ChainSafe/sygma-core/types"
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

type RelayerMetrics struct {
	meter metric.Meter
	Opts  api.MeasurementOption

	DepositEventCount        metric.Int64Counter
	MessageEventTime         map[string]time.Time
	ExecutionErrorCount      metric.Int64Counter
	ExecutionLatency         metric.Int64Histogram
	ExecutionLatencyPerRoute metric.Int64Histogram
	BlockDelta               metric.Int64ObservableGauge
	BlockDeltaMap            map[uint8]*big.Int

	lock sync.Mutex
}

// NewRelayerMetrics initializes OpenTelemetry metrics
func NewRelayerMetrics(meter metric.Meter, attributes ...attribute.KeyValue) (*RelayerMetrics, error) {
	opts := api.WithAttributes(attributes...)
	depositEventCounter, err := meter.Int64Counter(
		"relayer.DepositEventCount",
		metric.WithDescription("Number of deposit events per domain"))
	if err != nil {
		return nil, err

	}
	executionErrorCount, err := meter.Int64Counter(
		"relayer.ExecutionErrorCount",
		metric.WithDescription("Number of executions that failed"))
	if err != nil {
		return nil, err
	}
	executionLatencyPerRoute, err := meter.Int64Histogram(
		"relayer.ExecutionLatencyPerRoute",
		metric.WithDescription("Execution time histogram between indexing event and executing it per route"))
	if err != nil {
		return nil, err
	}
	executionLatency, err := meter.Int64Histogram(
		"relayer.ExecutionLatency",
		metric.WithDescription("Execution time histogram between indexing even`t and executing it"),
		metric.WithUnit("ms"))
	if err != nil {
		return nil, err
	}

	blockDeltaMap := make(map[uint8]*big.Int)

	blockDeltaGauge, err := meter.Int64ObservableGauge(
		"relayer.BlockDelta",
		metric.WithInt64Callback(func(context context.Context, result metric.Int64Observer) error {
			for domainID, delta := range blockDeltaMap {
				result.Observe(delta.Int64(),
					opts,
					metric.WithAttributes(attribute.Int64("domainID", int64(domainID))),
				)
			}
			return nil
		}),
		metric.WithDescription("Difference between chain head and current indexed block per domain"),
	)
	return &RelayerMetrics{
		meter:                    meter,
		MessageEventTime:         make(map[string]time.Time),
		Opts:                     opts,
		DepositEventCount:        depositEventCounter,
		ExecutionErrorCount:      executionErrorCount,
		ExecutionLatencyPerRoute: executionLatencyPerRoute,
		ExecutionLatency:         executionLatency,
		BlockDelta:               blockDeltaGauge,
		BlockDeltaMap:            blockDeltaMap,
	}, err
}

// TrackDepositMessage extracts metrics from deposit message and sends
// them to OpenTelemetry collector
func (t *RelayerMetrics) TrackDepositMessage(m *types.Message) {
	t.DepositEventCount.Add(context.Background(), 1, t.Opts, api.WithAttributes(attribute.Int64("source", int64(m.Source))))

	t.lock.Lock()
	defer t.lock.Unlock()
	t.MessageEventTime[m.ID()] = time.Now()
}

func (t *RelayerMetrics) TrackExecutionError(m *types.Message) {
	t.ExecutionErrorCount.Add(context.Background(), 1, t.Opts, api.WithAttributes(attribute.Int64("destination", int64(m.Source))))

	t.lock.Lock()
	defer t.lock.Unlock()
	delete(t.MessageEventTime, m.ID())
}

func (t *RelayerMetrics) TrackSuccessfulExecutionLatency(m *types.Message) {
	executionLatency := time.Since(t.MessageEventTime[m.ID()]).Milliseconds() / 1000
	t.ExecutionLatency.Record(context.Background(), executionLatency)
	t.ExecutionLatencyPerRoute.Record(
		context.Background(),
		executionLatency,
		t.Opts,
		api.WithAttributes(attribute.Int64("source", int64(m.Source))),
		api.WithAttributes(attribute.Int64("destination", int64(m.Destination))),
	)

	t.lock.Lock()
	defer t.lock.Unlock()
	delete(t.MessageEventTime, m.ID())
}

func (t *RelayerMetrics) TrackSuccessfulExecution(m *types.Message) {
	executionLatency := time.Since(t.MessageEventTime[m.ID()]).Milliseconds() / 1000
	t.ExecutionLatency.Record(context.Background(), executionLatency)
	t.ExecutionLatencyPerRoute.Record(
		context.Background(),
		executionLatency,
		t.Opts,
		api.WithAttributes(attribute.Int64("source", int64(m.Source))),
		api.WithAttributes(attribute.Int64("destination", int64(m.Destination))),
	)

	t.lock.Lock()
	defer t.lock.Unlock()
	delete(t.MessageEventTime, m.ID())
}

func (t *RelayerMetrics) TrackBlockDelta(domainID uint8, head *big.Int, current *big.Int) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.BlockDeltaMap[domainID] = new(big.Int).Sub(head, current)
}
