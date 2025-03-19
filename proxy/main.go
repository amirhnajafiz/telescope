package main

import (
	"github.com/amirhnajafiz/telescope/internal/api"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// create a new fiber app
	app := fiber.New()

	// create a new API instance
	apiInstance := api.API{}

	// register the API endpoints
	apiInstance.Register(app)

	// start the server on port 3000
	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
