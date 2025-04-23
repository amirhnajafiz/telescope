package middlewares

import "github.com/gofiber/fiber/v2"

// ExtractHeaders is a middleware that extracts specific headers from the request
func ExtractHeaders(ctx *fiber.Ctx) error {
	// extract headers from the request and add them to the context
	ctx.Locals("XSC", ctx.Get("X-Server-Cached", "0"))
	ctx.Locals("XSU", ctx.Get("X-Server-Uncached", "0"))
	ctx.Locals("XSCB", ctx.Get("X-Server-Current-Bandwidth", "0"))
	ctx.Locals("XSA", ctx.Get("X-Server-Alpha", "0"))

	return ctx.Next()
}
