package tracer

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	tr "go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"time"
)

type Config struct {
	ServiceName        string
	Host               string
	Port               string
	Environment        string
	TraceRatioFraction float64

	// OTELExporterOTLPEndpoint example: http://jaeger:4317
	OTELExporterOTLPEndpoint string
}

func New(cfg *Config) (func(ctx context.Context), error) {
	if err := os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", cfg.OTELExporterOTLPEndpoint); err != nil {
		return func(ctx context.Context) {}, errors.Wrap(err, "fail to set env var otel exporter")
	}

	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)),
	)

	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, err
	}
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.TraceRatioFraction)),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			attribute.String("environment", cfg.Environment),
		)),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(traceProvider)

	// Return func for graceful shutdown tracer
	return func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := traceProvider.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}, nil
}

func StartTrace(ctx context.Context, spanName string) (context.Context, tr.Span) {
	tp := otel.GetTracerProvider()
	t := tp.Tracer("")
	return t.Start(ctx, spanName)
}
