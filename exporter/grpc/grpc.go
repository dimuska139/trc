package grpc

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func NewExporter(
	ctx context.Context,
	collectorGRPChost string,
	transportCredentials credentials.TransportCredentials,
) (sdkTrace.SpanExporter, error) {
	if transportCredentials == nil {
		transportCredentials = insecure.NewCredentials()
	}

	conn, err := grpc.NewClient(collectorGRPChost, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithGRPCConn(conn),
	)
	if err != nil {
		return nil, fmt.Errorf("init gRPC exporter: %w", err)
	}

	return traceExporter, nil
}
