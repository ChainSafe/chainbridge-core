package observability

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	traceapi "go.opentelemetry.io/otel/trace"
)

// CreateSpanAndLoggerFromContext creates span and logger from context with provided name.
// Logger explicitly extended with dd.trace_id attribute for DataDog logs and traces connections
func CreateSpanAndLoggerFromContext(ctx context.Context, tracerName, spanName string, kv ...attribute.KeyValue) (context.Context, traceapi.Span, zerolog.Logger) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, spanName)
	l := zerolog.Ctx(ctx).With().Str("dd.trace_id", span.SpanContext().TraceID().String()).Logger()
	setOTLPAttrsToLogger(&l, kv)
	span.SetAttributes(kv...)

	return ctx, span, l
}

func CreateSpanAndLoggerFromExternalTraceID(ctx context.Context, traceID, tracerName, spanName string, kv ...attribute.KeyValue) (context.Context, traceapi.Span, zerolog.Logger) {
	tID, err := traceapi.TraceIDFromHex(traceID)
	if err == nil {
		ctx = traceapi.ContextWithSpanContext(ctx, traceapi.NewSpanContext(traceapi.SpanContextConfig{TraceID: tID, Remote: true}))
	}
	return CreateSpanAndLoggerFromContext(ctx, tracerName, spanName, kv...)
}

func SetSpanAndLoggerAttrs(logger *zerolog.Logger, span traceapi.Span, kv ...attribute.KeyValue) {
	setOTLPAttrsToLogger(logger, kv)
	span.SetAttributes(kv...)
}

// LogAndEvent CreatesTrace Span event with attributes and provided name and also creates log with the same attributes
// and name to keep necessary informing level for developers that are not using Trace onitoring tool
func LogAndEvent(logger *zerolog.Event, span traceapi.Span, msg string, kv ...attribute.KeyValue) {
	span.AddEvent(msg, traceapi.WithAttributes(kv...))
	if logger != nil {
		logger.Msg(msg)
		addOTLPAttrsToLogEvent(logger, kv)
	}
}

// LogAndRecordError CreatesTrace Span event with attributes and provided name and also creates log with the same attributes
// and name to keep necessary informing level for developers that are not using Trace onitoring tool
// use logger == nil for not logging an error
func LogAndRecordError(logger *zerolog.Logger, span traceapi.Span, err error, msg string, kv ...attribute.KeyValue) error {
	err = fmt.Errorf("%s with err: %e", msg, err)
	span.RecordError(err, traceapi.WithAttributes(kv...))
	if logger != nil {
		setOTLPAttrsToLogger(logger, kv)
		logger.Err(err).Msg(msg)
	}
	return err
}

// LogAndRecordErrorWithStatus Records error to traces logs it and set error status to span
// Should be used when span will be ended afterwards. Corresponding span will be marked as errored
// use logger == nil for not logging an error
// Returns error wrapped with message
func LogAndRecordErrorWithStatus(logger *zerolog.Logger, span traceapi.Span, err error, msg string, kv ...attribute.KeyValue) error {
	err = fmt.Errorf("%s with err: %e", msg, err)
	span.RecordError(err, traceapi.WithAttributes(kv...))
	span.SetStatus(codes.Error, err.Error())
	if logger != nil {
		setOTLPAttrsToLogger(logger, kv)
		logger.Err(err).Msg(msg)
	}
	return err
}

func SetAttrsToSpanAnLogger(logger *zerolog.Logger, span traceapi.Span, kv ...attribute.KeyValue) {
	span.SetAttributes(kv...)
	setOTLPAttrsToLogger(logger, kv)
}

func setOTLPAttrsToLogger(logger *zerolog.Logger, attrs []attribute.KeyValue) {
	for _, attr := range attrs {
		switch attr.Value.Type() {
		case attribute.STRING:
			logger.With().Str(string(attr.Key), attr.Value.AsString())
		case attribute.BOOL:
			logger.With().Bool(string(attr.Key), attr.Value.AsBool())
		case attribute.INT64:
			logger.With().Int64(string(attr.Key), attr.Value.AsInt64())
		default:
			logger.With().Str(string(attr.Key), fmt.Sprintf("%+v", attr.Value))
		}
	}
}

func addOTLPAttrsToLogEvent(logger *zerolog.Event, attrs []attribute.KeyValue) {
	for _, attr := range attrs {
		switch attr.Value.Type() {
		case attribute.STRING:
			logger.Str(string(attr.Key), attr.Value.AsString())
		case attribute.BOOL:
			logger.Bool(string(attr.Key), attr.Value.AsBool())
		case attribute.INT64:
			logger.Int64(string(attr.Key), attr.Value.AsInt64())
		default:
			logger.Str(string(attr.Key), fmt.Sprintf("%+v", attr.Value))
		}
	}
}
