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
	logger, err := logr.NewZapLogger(cfg.Debug)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// create a new otel tracer
	var tr trace.Tracer
	if len(cfg.Jaeger) > 0 {
		tr, err = tracing.NewProductionTracer(cfg.Jaeger)
	} else {
		tr, err = tracing.NewDevelopmentTracer()
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

	// create a new IPFS client instance
	ipfsClient := ipfs.NewClient(cfg.IPFSGateway)

	// create a new cache instance
	cacheInstance := cache.NewCache(cfg.CachePath)

	// create a new abr-rewriter controller
	abrRewriter := controllers.AbrRewriter{
		Estimator: controllers.NewEstimator(),
		Cache:     cacheInstance,
	}

	// create a new API instance
	return &api.API{
		Logr:        logger.Named("api"),
		Metrics:     metricsInstance,
		Tracer:      tr,
		IPFS:        ipfsClient,
		Cache:       cacheInstance,
		ABRRewriter: abrRewriter,
	}, nil
}
