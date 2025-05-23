package cmd

import (
	"fmt"

	"github.com/amirhnajafiz/telescope/internal/api"
	"github.com/amirhnajafiz/telescope/internal/config"
	"github.com/amirhnajafiz/telescope/internal/controllers"
	"github.com/amirhnajafiz/telescope/internal/logr"
	"github.com/amirhnajafiz/telescope/internal/storage/cache"
	"github.com/amirhnajafiz/telescope/internal/storage/ipfs"
	"github.com/amirhnajafiz/telescope/internal/telemetry/metrics"
	"github.com/amirhnajafiz/telescope/internal/telemetry/tracing"

	"go.opentelemetry.io/otel/trace"
)

// RegisterAPI creates a new API instance
func RegisterAPI(cfg *config.Config) (*api.API, error) {
	// create a new logger instance
	logger := logr.NewZapLogger(cfg.Debug)

	// create a new otel tracer
	var (
		tr  trace.Tracer
		err error
	)
	if len(cfg.Jaeger) > 0 {
		tr, err = tracing.NewProductionTracer(cfg.Jaeger)
		logger.Info("Jaeger tracing enabled")
	} else {
		tr, err = tracing.NewDevelopmentTracer()
		logger.Info("Jaeger tracing disabled")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer: %w", err)
	}

	// create new metrics struct
	metricsInstance := metrics.NewMetrics()

	// check if metrics port is set
	if cfg.MetricsPort != 0 {
		metrics.NewServer(cfg.MetricsPort)
	}

	// create a new cache instance
	cacheInstance := cache.NewCache(cfg.CachePath)

	// create a new API instance
	return &api.API{
		Logr:        logger.Named("api"),
		Metrics:     metricsInstance,
		Tracer:      tr,
		IPFS:        ipfs.NewClient(cfg.IPFSGateway),
		Cache:       cacheInstance,
		ABRRewriter: controllers.NewAbrRewriter(cacheInstance, logger.Named("abr-rewriter"), tr),
		MPDBuilder:  controllers.NewMPDBuilder(logger.Named("mpd-builder"), tr),
	}, nil
}
