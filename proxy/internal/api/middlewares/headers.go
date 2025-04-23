package middlewares

import (
	"github.com/amirhnajafiz/telescope/pkg/parser"

	"github.com/gofiber/fiber/v2"
)

// ExtractHeaders is a middleware that extracts headers from the request
func ExtractHeaders(ctx *fiber.Ctx) error {
	// extract headers from the request and add them to the context
	header := ctx.Get("X-Server-BW", "")

	// if the header is empty, set default values
	// otherwise, parse the header into a map
	var headerMap map[string]float64
	if header == "" {
		headerMap = map[string]float64{
			"cached":   0,
			"uncached": 0,
			"current":  0,
		}
	} else {
		headerMap = parser.ParseHeaderToMap(header)
	}

	// set the header map in the context
	ctx.Locals("headerMap", headerMap)

	return ctx.Next()
}
