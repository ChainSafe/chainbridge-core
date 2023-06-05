package opentelemetry

import (
	"context"
	"math/big"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/unit"
)

type ChainbridgeMetrics struct {
	DepositEventCount        metric.Int64Counter
	ExecutionErrorCount      metric.Int64Counter
	ExecutionLatency         metric.Int64Histogram
	ExecutionLatencyPerRoute metric.Int64Histogram
	BlockDelta               metric.Int64GaugeObserver

	BlockDeltaMap map[uint8]*big.Int
}

// NewChainbridgeMetrics creates an instance of ChainbridgeMetrics
// with provided OpenTelemetry meter
func NewChainbridgeMetrics(meter metric.Meter, genericLabels ...attribute.KeyValue) *ChainbridgeMetrics {
	blockDeltaMap := make(map[uint8]*big.Int)
	return &ChainbridgeMetrics{
		DepositEventCount: metric.Must(meter).NewInt64Counter(
			"chainbridge.DepositEventCount",
			metric.WithDescription("Number of deposit events per domain"),
		),
		ExecutionErrorCount: metric.Must(meter).NewInt64Counter(
			"chainbridge.ExecutionErrorCount",
			metric.WithDescription("Number of executions that failed"),
		),
		ExecutionLatencyPerRoute: metric.Must(meter).NewInt64Histogram(
			"chainbridge.ExecutionLatencyPerRoute",
			metric.WithDescription("Execution time histogram between indexing event and executing it per route"),
		),
		ExecutionLatency: metric.Must(meter).NewInt64Histogram(
			"chainbridge.ExecutionLatency",
			metric.WithDescription("Execution time histogram between indexing event and executing it"),
			metric.WithUnit(unit.Milliseconds),
		),
		BlockDelta: metric.Must(meter).NewInt64GaugeObserver(
			"chainbridge.BlockDelta",
			func(ctx context.Context, result metric.Int64ObserverResult) {
				for domainID, delta := range blockDeltaMap {
					labels := append(genericLabels, attribute.Int64("domainID", int64(domainID)))
					result.Observe(delta.Int64(), labels...)
				}
			},
			metric.WithDescription("Difference between chain head and current indexed block per domain"),
		),
		BlockDeltaMap: blockDeltaMap,
	}
}
