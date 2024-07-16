package trc

import (
	"context"
	"go.opentelemetry.io/otel"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func StartSpan(ctx context.Context, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)

	var funcName string
	if ok && details != nil {
		funcName = details.Name()
	}

	// For testing purposes
	if tracer == nil {
		return otel.Tracer("").Start(ctx, funcName, trace.WithAttributes(attrs...), trace.WithTimestamp(time.Now()))
	}

	return tracer.Start(ctx, funcName, trace.WithAttributes(attrs...), trace.WithTimestamp(time.Now()))
}

func EndSpan(span trace.Span, err error, attrs ...attribute.KeyValue) {
	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	span.End()
}
