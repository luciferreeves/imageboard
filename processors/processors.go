package processors

import "github.com/gofiber/fiber/v2"

func Initialize(app *fiber.App) {
	app.Use(MetaContextProcessor)
}
