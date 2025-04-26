package middlewares

import (
	"strconv"

	"github.com/amirhnajafiz/telescope/internal/telemetry/metrics"

	"github.com/gofiber/fiber/v2"
)

// ExtractHeadersToPrometheus is a middleware that extracts headers from the request
// and sets them to Prometheus metrics
func ExtractHeadersToPrometheus(m *metrics.Metrics) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// extract headers
		clientID := c.Get("X-Client-ID", "default")
		bandwidth, _ := strconv.ParseFloat(c.Get("X-Bandwidth", "0"), 64)
		stRate, _ := strconv.ParseFloat(c.Get("X-Stall-Rate", "0"), 64)
		quality, _ := strconv.Atoi(c.Get("X-Segment-Quality", "0"))

		// set Prometheus metrics
		m.ClientBandwidth.WithLabelValues(clientID).Set(bandwidth)
		m.ClientStallRate.WithLabelValues(clientID).Add(stRate)
		m.ClientQuality.WithLabelValues(clientID).Set(float64(quality))

		return c.Next()
	}
}
