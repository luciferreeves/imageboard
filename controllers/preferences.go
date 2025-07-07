package controllers

import (
	"imageboard/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func PreferencesController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", "Site Preferences")
	ctx.Locals("request", fiber.Map{"path": ctx.Path()})
	return shortcuts.Render(ctx, "preferences", nil)
}
