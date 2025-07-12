package controllers

import (
	"imageboard/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func PreferencesPageController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", "Site Preferences")
	return shortcuts.Render(ctx, "preferences", nil)
}
