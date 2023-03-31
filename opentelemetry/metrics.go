package opentelemetry

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/unit"
)

type ChainbridgeMetrics struct {
	DepositEventCount        metric.Int64Counter
	ExecutionErrorCount      metric.Int64Counter
	ExecutionLatency         metric.Int64Histogram
	ExecutionLatencyPerRoute metric.Int64Histogram
}

// NewChainbridgeMetrics creates an instance of ChainbridgeMetrics
// with provided OpenTelemetry meter
func NewChainbridgeMetrics(meter metric.Meter) *ChainbridgeMetrics {
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
			metric.WithUnit(unit.Milliseconds),
		),
		ExecutionLatency: metric.Must(meter).NewInt64Histogram(
			"chainbridge.ExecutionLatency",
			metric.WithDescription("Execution time histogram between indexing event and executing it"),
			metric.WithUnit(unit.Milliseconds),
		),
	}
}
