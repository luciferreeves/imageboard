package main

import (
	"fmt"
	"imageboard/config"
	"imageboard/middleware"
	"imageboard/processors"
	"imageboard/router"
	"imageboard/utils/handlers"
	"log"

	_ "imageboard/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/django/v3"
)

func main() {
	if config.Server.AppSecret == config.Defaults(&config.Server).AppSecret {
		log.Println("Warning: AppSecret is set to a default value which is not secure. Please set a strong random secret in your APP_SECRET environment variable or .env file.")
	}

	engine := django.New("./templates", ".django")
	engine.Reload(config.Server.IsDevMode)
	app := fiber.New(fiber.Config{
		Views:        engine,
		ErrorHandler: handlers.ServerErrorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(helmet.New(helmet.Config{
		CrossOriginEmbedderPolicy: "unsafe-none",
	}))
	app.Use(cors.New())

	processors.Initialize(app)
	middleware.Initialize(app)

	app.Static("/static", "./static")

	router.Initialize(app)

	log.Fatalf("Server failed to start: %v", app.Listen(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)))
}
