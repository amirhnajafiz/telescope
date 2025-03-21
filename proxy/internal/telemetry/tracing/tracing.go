package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	otrace "go.opentelemetry.io/otel/trace"
)

const (
	AgentName  = "telescope-proxy-tracer"
	TracerName = "telescope-proxy"
)

// NewProductionTracer creates a tracer for production use
func NewProductionTracer(endpoint string) (otrace.Tracer, error) {
	// build agent
	agent, err := NewAgent(endpoint)
	if err != nil {
		return nil, err
	}

	// create trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(agent),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(AgentName),
		)),
	)

	// set the global trace provider
	otel.SetTracerProvider(tp)

	// create tracer
	tracer := tp.Tracer(TracerName)

	return tracer, nil
}

// NewDevelopmentTracer creates a tracer for development purposes
func NewDevelopmentTracer() (otrace.Tracer, error) {
	// create trace provider
	tp := trace.NewTracerProvider()

	// create tracer
	tracer := tp.Tracer(TracerName)

	return tracer, nil
}
