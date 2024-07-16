# Tracing

This is a simple library to trace the execution of your GO services

## Installation

```bash
go get github.com/dimuska139/trc
```

## Usage

### Initialize tracing provider
```go
import (
	"github.com/dimuska139/trc"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	grpcExporter "github.com/dimuska139/trc/exporter/grpc"
)
...
spanExporter, err := grpcExporter.NewExporter(ctx, "localhost:4318", insecure.NewCredentials())
if err != nil {
	logger.Fatalf("Can't initialize gRPC span exporter: %s", err)
}

tracingProvider, err := trc.NewProvider(ctx, trc.Config{
    Name: "mytracer",
	SampleRatio: "1",
	SpanExporter: spanExporter,
},
    semconv.ServiceName("myservice"),
    semconv.ServiceVersion("1.0.0"),
    semconv.ServiceInstanceID("stage"),
    semconv.DeploymentEnvironment("dev"),
)
if err != nil {
    logger.Fatalf("Can't initialize tracing: %s", err)
}
...

```
Don't forget to shutdown tracing provider when you don't need it anymore
```go
if stopErr := tracingProvider.Shutdown(ctx); stopErr != nil {
    logger.Fatalf("Can't shutdown tracing provider")
}
```

### Add tracing middleware to you server

See https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation

### Adding tracing to your functions

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr>
<td>

```go
func count(ctx context.Context) (_ int, err error){
    ctx, span := trc.StartSpan(ctx)
    defer trc.EndSpan(span, err)
    ...
    // Do some work
    ...
}
```
</td>
<td>

```go
func count(ctx context.Context) (_ int, err error) {
    ctx, span := trc.StartSpan(ctx)
    defer func() { trc.EndSpan(span, err) }()
    ...
    // Do some work
    ...
}
```
</td>
</tr>
<tr>
<td>

```plain
`defer` call should be wrapped
```
</td>
<td>
</td>
</tr>
<tr>
<td>

```go
func count(ctx context.Context) error {
    ctx, span := trc.StartSpan(ctx)
    defer func() { trc.EndSpan(span, err) }()
    ...
    // Do some work
    ...
}
```
</td>
<td>
</td>
</tr>
<tr>
<td>

```plain
`err` should be declared as named return value
```
</td>
<td>
</td>
</tr>
</tbody></table>

You can use this `docker-compose.yml` configuration for local development
```yaml
version: "3.1"

services:
  # Jaeger
  jaeger-all-in-one:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"
      - "14268"
      - "14250"

  zipkin-all-in-one:
    image: openzipkin/zipkin:latest
    environment:
      - JAVA_OPTS=-Xms1024m -Xmx1024m -XX:+ExitOnOutOfMemoryError
    ports:
      - "9411:9411"

  otel-collector:
    image: otel/opentelemetry-collector:0.88.0
    command: [ "--config=/etc/otel-collector-config.yaml", "${OTELCOL_ARGS}" ]
    volumes:
      - ./config/otel-collector-config.yaml:/etc/otel-collector-config.yaml
    ports:
      #- "1888:1888"   # pprof extension
      #- "8888:8888"   # Prometheus metrics exposed by the collector
      #- "8889:8889"   # Prometheus exporter metrics
      #- "13133:13133" # health_check extension
      - "4318:4317"   # OTLP gRPC receiver
      #- "55679:55679" # zpages extension
    depends_on:
      - jaeger-all-in-one
      - zipkin-all-in-one
```
`otel-collector-config.yaml`
```yaml
receivers:
  otlp:
    protocols:
      grpc:

exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"
    const_labels:
      label1: value1

  debug:

  zipkin:
    endpoint: "http://zipkin-all-in-one:9411/api/v2/spans"
    format: proto

  otlp:
    endpoint: jaeger-all-in-one:4317
    tls:
      insecure: true

processors:
  batch:

extensions:
  health_check:
  pprof:
    endpoint: :1888
  zpages:
    endpoint: :55679

service:
  extensions: [pprof, zpages, health_check]
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [debug, zipkin, otlp]
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [debug, prometheus]
```