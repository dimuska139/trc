package trc

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/dimuska139/trc/exporter/dummy"
)

var tracer trace.Tracer

func NewProvider(ctx context.Context, config Config, attrs ...attribute.KeyValue) (*sdktrace.TracerProvider, error) {
	tracer = otel.Tracer(config.Name)

	res, err := resource.New(ctx,
		resource.WithAttributes(
			attrs...,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	if config.SpanExporter == nil {
		config.SpanExporter = dummy.NewExporter()
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(
			sdktrace.ParentBased(sdktrace.TraceIDRatioBased(config.SampleRatio)),
		),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(
			sdktrace.NewBatchSpanProcessor(config.SpanExporter),
		),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider, nil
}
