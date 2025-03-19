package api

import "github.com/gofiber/fiber/v2"

// healthCheck endpoint returns the health status of the server
func (a *API) healthCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status":  "ok",
		"message": "The server is healthy!",
	})
}
