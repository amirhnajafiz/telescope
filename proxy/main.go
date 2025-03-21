package main

import (
	"fmt"

	"github.com/amirhnajafiz/telescope/internal/api"
	"github.com/amirhnajafiz/telescope/internal/config"
	"github.com/amirhnajafiz/telescope/internal/logr"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// load configs
	cfg, err := config.LoadConfigs()
	if err != nil {
		panic(err)
	}

	// create a new logger instance
	logger, err := logr.NewZapLogger(cfg.Debug)
	if err != nil {
		panic(err)
	}

	// create a new fiber app
	app := fiber.New()

	// create a new API instance
	apiInstance := api.API{
		Logr: logger.Named("api"),
	}

	// register the API endpoints
	apiInstance.Register(app)

	// start the server on port 3000
	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
