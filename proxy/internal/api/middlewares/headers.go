package middlewares

import (
	"github.com/amirhnajafiz/telescope/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// ExtractHeaders is a middleware that extracts headers from the request
func ExtractHeaders(ctx *fiber.Ctx) error {
	// extract headers from the request and add them to the context
	xsc, _ := utils.StrToFloat64(ctx.Get("X-Server-Cached", "0"))
	xsu, _ := utils.StrToFloat64(ctx.Get("X-Server-Uncached", "0"))
	xsb, _ := utils.StrToFloat64(ctx.Get("X-Server-Current-Bandwidth", "0"))

	ctx.Locals("XSC", xsc)
	ctx.Locals("XSU", xsu)
	ctx.Locals("XSCB", xsb)

	return ctx.Next()
}
