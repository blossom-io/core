package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// initTracer creates and registers trace provider instance.
// func New() (*sdktrace.TracerProvider, error) {
// 	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to initialize stdouttrace exporter: %w", err)
// 	}
// 	bsp := sdktrace.NewBatchSpanProcessor(exp)
// 	tp := sdktrace.NewTracerProvider(
// 		sdktrace.WithSampler(sdktrace.AlwaysSample()),
// 		sdktrace.WithSpanProcessor(bsp),
// 	)

// 	otel.SetTracerProvider(tp)

// 	return tp, nil
// }

func New() (*tracesdk.TracerProvider, error) {
	exp, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("core"),
			attribute.String("environment", "dev"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}
