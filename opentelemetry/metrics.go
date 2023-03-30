package opentelemetry

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/unit"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

type ChainbridgeMetrics struct {
	DepositEventCount   metric.Int64Counter
	ExecutionErrorCount metric.Int64Counter
	ExecutionLatency    metric.Int64Histogram
}

// NewChainbridgeMetrics creates an instance of ChainbridgeMetrics
// with provided OpenTelemetry meter
func NewChainbridgeMetrics(meter metric.Meter) *ChainbridgeMetrics {
	return &ChainbridgeMetrics{
		DepositEventCount: metric.Must(meter).NewInt64Counter(
			"chainbridge.DepositEventCount",
			metric.WithDescription("Number of deposit events across all chains"),
		),
		ExecutionErrorCount: metric.Must(meter).NewInt64Counter(
			"chainbridge.ExecutionErrorCount",
			metric.WithDescription("Number of executions that failed"),
		),
		ExecutionLatency: metric.Must(meter).NewInt64Histogram(
			"chainbridge.ExecutionLatency",
			metric.WithDescription("Execution time histogram between indexing event and executing it"),
			metric.WithUnit(unit.Milliseconds),
		),
	}
}

func InitOpenTelemetryMetrics(opts ...otlpmetrichttp.Option) (*ChainbridgeMetrics, error) {
	ctx := context.Background()

	client := otlpmetrichttp.NewClient(opts...)
	exp, err := otlpmetric.New(ctx, client)
	if err != nil {
		return nil, err
	}

	selector := simple.NewWithInexpensiveDistribution()
	proc := processor.NewFactory(selector, export.CumulativeExportKindSelector())
	cont := controller.New(proc, controller.WithExporter(exp))
	global.SetMeterProvider(cont)

	err = cont.Start(ctx)
	if err != nil {
		return nil, err
	}

	meter := cont.Meter("chainbridge")
	metrics := NewChainbridgeMetrics(meter)

	return metrics, nil
}
