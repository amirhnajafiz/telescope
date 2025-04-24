package main

import (
	"fmt"
	"log"

	"github.com/amirhnajafiz/telescope/cmd"
	"github.com/amirhnajafiz/telescope/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// load configs
	cfg, err := config.LoadConfigs()
	if err != nil {
		log.Fatalln(err)
	}

	// build html engine
	engine := html.New("./public", ".html")

	// create a new fiber app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// create a new API instance
	apiInstance, err := cmd.RegisterAPI(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	// register the API endpoints
	apiInstance.Register(app)

	// start the server on a specified port
	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		log.Fatalln(err)
	}
}
