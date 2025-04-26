package middlewares

import (
	"github.com/gofiber/fiber/v2"

	"go.uber.org/zap"
)

// Prelog is a middleware that logs the request method and path before processing the request
func Prelog(logr *zap.Logger) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logr.Debug("Prelog",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)

		logr.Debug("Headers",
			zap.String("seg-q", c.Get("X-Segment-Quality")),
			zap.String("bw", c.Get("X-Bandwidth")),
			zap.String("rt", c.Get("X-Stall-Rate")),
		)

		return c.Next()
	}
}
