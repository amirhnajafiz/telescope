package main

import (
	"fmt"

	"github.com/amirhnajafiz/telescope/internal/api"
	"github.com/amirhnajafiz/telescope/internal/config"
	"github.com/amirhnajafiz/telescope/internal/logr"
	"github.com/amirhnajafiz/telescope/internal/telemetry/metrics"
	"github.com/amirhnajafiz/telescope/internal/telemetry/tracing"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	// load configs
	cfg, err := config.LoadConfigs()
	if err != nil {
		panic(err)
	}

	// create a new logger instance
	logger, err := logr.NewZapLogger(cfg.Debug)
	if err != nil {
		panic(err)
	}

	// create a new otel tracer
	var tr trace.Tracer
	if len(cfg.Jaeger) > 0 {
		tr, err = tracing.NewProductionTracer(cfg.Jaeger)
	} else {
		tr, err = tracing.NewDevelopmentTracer()
	}
	if err != nil {
		panic(err)
	}

	// create new metrics struct
	metricsInstance := metrics.NewMetrics()

	// check if metrics port is set
	if cfg.MetricsPort != 0 {
		metrics.NewServer(cfg.MetricsPort)
	}

	// create a new fiber app
	app := fiber.New()

	// create a new API instance
	apiInstance := api.API{
		Logr:    logger.Named("api"),
		Metrics: metricsInstance,
		Tracer:  tr,
	}

	// register the API endpoints
	apiInstance.Register(app)

	// start the server on port 3000
	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
