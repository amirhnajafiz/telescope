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

		return c.Next()
	}
}
