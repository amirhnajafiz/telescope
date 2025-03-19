package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// API struct holds endpoint functions of the proxy server
type API struct{}

// Register method takes a fiber.App instance and defines all the endpoints
func (a *API) Register(app *fiber.App) {
	// enable CORS for all routes
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
}
