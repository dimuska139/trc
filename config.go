package trc

import sdktrace "go.opentelemetry.io/otel/sdk/trace"

type Config struct {
	Name         string
	SampleRatio  float64
	SpanExporter sdktrace.SpanExporter
}
