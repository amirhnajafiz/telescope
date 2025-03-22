package main

import (
	"fmt"

	"github.com/amirhnajafiz/telescope/cmd"
	"github.com/amirhnajafiz/telescope/internal/config"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// load configs
	cfg, err := config.LoadConfigs()
	if err != nil {
		panic(err)
	}

	// create a new fiber app
	app := fiber.New()

	// create a new API instance
	apiInstance, err := cmd.RegisterAPI(cfg)
	if err != nil {
		panic(err)
	}

	// register the API endpoints
	apiInstance.Register(app)

	// start the server on a specified port
	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
