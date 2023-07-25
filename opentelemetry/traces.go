package opentelemetry

import (
	"context"
	"net/url"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracesProvider(ctx context.Context, res *sdkresource.Resource, agentURL string) (*tracesdk.TracerProvider, error) {
	collectorURL, err := url.Parse(agentURL)
	if err != nil {
		return nil, err
	}

	traceOptions := []otlptracehttp.Option{
		otlptracehttp.WithURLPath(collectorURL.Path),
		otlptracehttp.WithEndpoint(collectorURL.Host),
	}
	if collectorURL.Scheme == "http" {
		traceOptions = append(traceOptions, otlptracehttp.WithInsecure())
	}

	traceHTTP := otlptracehttp.NewClient(traceOptions...)

	traceExp, err := otlptrace.New(ctx, traceHTTP)
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(traceExp),
		tracesdk.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}
