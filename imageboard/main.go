package main

import (
	"imageboard/config"
	"imageboard/middleware"
	"imageboard/processors"
	"imageboard/router"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func main() {
	if config.AppSecret == "default_secret" {
		log.Println("Warning: AppSecret is set to a default value which is not secure. Please set a strong random secret in your APP_SECRET environment variable or .env file.")
	}

	engine := html.New("./templates", ".html")
	engine.Reload(config.IsDevelopmentMode)
	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			log.Printf("Error: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		},
		BodyLimit: config.Image.MaxSize,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(helmet.New())
	app.Use(cors.New())

	processors.Initialize(app)

	app.Use(middleware.JSON)
	app.Static("/", "./static")

	router.Initialize(app)

	log.Fatalf("Server failed to start: %v", app.Listen(config.Server.Host+":"+config.Server.Port))
}
