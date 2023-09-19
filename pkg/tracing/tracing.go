package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func newResource(sn, sv string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(sn),
		semconv.ServiceVersion(sv),
	)
}

func InstallExportPipeline(ctx context.Context,
	serviceName, serviceVersion string) (
	func(context.Context) error, error,
) {
	if serviceName == "" || serviceVersion == "" {
		return nil, fmt.Errorf("serviceName nor serviceVersion can be empty")
	}

	client := otlptracehttp.NewClient(otlptracehttp.WithInsecure())
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource(serviceName, serviceVersion)),
	)
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider.Shutdown, nil
}
