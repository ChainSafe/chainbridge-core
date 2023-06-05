package opentelemetry

import (
	"go.opentelemetry.io/otel/metric"
)

type ChainbridgeMetrics struct {
	DepositEventCount        metric.Int64Counter
	ExecutionErrorCount      metric.Int64Counter
	ExecutionLatency         metric.Int64Histogram
	ExecutionLatencyPerRoute metric.Int64Histogram
}

// NewChainbridgeMetrics creates an instance of ChainbridgeMetrics
// with provided OpenTelemetry meter
func NewChainbridgeMetrics(meter metric.Meter) (*ChainbridgeMetrics, error) {
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
	return &ChainbridgeMetrics{
		DepositEventCount:        depositEventCounter,
		ExecutionErrorCount:      executionErrorCount,
		ExecutionLatencyPerRoute: executionLatencyPerRoute,
		ExecutionLatency:         executionLatency,
	}, nil
}
