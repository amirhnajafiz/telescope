package api

import (
	"github.com/amirhnajafiz/telescope/internal/abr"
	"github.com/amirhnajafiz/telescope/internal/cache"
	"github.com/amirhnajafiz/telescope/internal/ipfs"
	"github.com/amirhnajafiz/telescope/internal/telemetry/metrics"
	"github.com/amirhnajafiz/telescope/internal/throughput"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// API struct holds endpoint functions of the proxy server
type API struct {
	Logr      *zap.Logger
	Metrics   *metrics.Metrics
	Tracer    trace.Tracer
	IPFS      ipfs.Client
	ABR       abr.CacheBasedPolicy
	Cache     *cache.SegmentCache
	Estimator *throughput.Estimator
}

// Register method takes a fiber.App instance and defines all the endpoints
func (a *API) Register(app *fiber.App) {
	// enable CORS for all endpoints
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// define the health check endpoint
	app.Get("/healthz", a.healthCheck)

	// enable logging for all endpoints
	app.Use(logger.New(logger.Config{
		Format: "${status} - ${method} ${path} ${latency}",
		Done: func(c *fiber.Ctx, logString []byte) {
			a.Logr.Info(string(logString))
		},
	}))

	// create API groups
	contents := app.Group("/api/contents")

	// define the contents endpoints
	contents.Get("/", a.listContents)
	contents.Put("/", a.newContent)
	contents.Get("/:cid", a.getContent)
	contents.Get("/:cid/stream", a.streamContent)
}
