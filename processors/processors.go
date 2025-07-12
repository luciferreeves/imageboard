package processors

import "github.com/gofiber/fiber/v2"

func Initialize(app *fiber.App) {
	app.Use(RequestContextProcessor)
	app.Use(MetaContextProcessor)
	app.Use(SidebarContextProcessor)
	app.Use(PreferencesContextProcessor)
}
