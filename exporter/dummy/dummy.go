package dummy

import (
	"context"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

type Exporter struct {
}

func (exp *Exporter) ExportSpans(_ context.Context, _ []sdkTrace.ReadOnlySpan) error {
	return nil
}

func (exp *Exporter) Shutdown(_ context.Context) error {
	return nil
}

func NewExporter() sdkTrace.SpanExporter {
	return &Exporter{}
}
