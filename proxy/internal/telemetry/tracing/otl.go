package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	otrace "go.opentelemetry.io/otel/trace"
)

// NewTracer creates a new TracerProvider with the specified endpoint
func NewTracer(endpoint string) (otrace.Tracer, error) {
	var tp *trace.TracerProvider

	if endpoint != "" {
		// build agent
		agent, err := NewAgent(endpoint)
		if err != nil {
			return nil, err
		}

		// create trace provider
		tp = trace.NewTracerProvider(
			trace.WithBatcher(agent),
			trace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("telescope-proxy-tracer"),
			)),
		)

		// set the global trace provider
		otel.SetTracerProvider(tp)
	} else {
		tp = trace.NewTracerProvider()
	}

	// create tracer
	tracer := tp.Tracer("telescope-proxy")

	return tracer, nil
}
