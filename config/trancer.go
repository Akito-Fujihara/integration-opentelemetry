package config

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracer() (*trace.TracerProvider, error) {
	// 標準出力にトレースを出力するエクスポーターを作成
	// traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint()

	ctx := context.Background()

	collectorAddress := os.Getenv("OTEL_COLLECTOR_ADDRESS")
	if collectorAddress == "" {
		collectorAddress = "otel-collector:4317"
	}

	// gRPCを使用してOpenTelemetry Collectorと通信するクライアントを作成します
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(collectorAddress),
		otlptracegrpc.WithInsecure(),
	)

	// 作成したクライアントを使用してトレースデータをエクスポートするエクスポーターを作成します
	traceExporter, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		log.Fatalf("Failed to create trace exporter: %v", err)
	}

	if err != nil {
		return nil, err
	}

	// TracerProviderを設定
	tp := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("echo-server"),
			attribute.String("environment", "development"),
		)),
	)

	// OpenTelemetryのトレーサーをグローバルに設定
	otel.SetTracerProvider(tp)

	return tp, nil
}
