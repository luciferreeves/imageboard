package controllers

import (
	"imageboard/config"
	"imageboard/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func HomePageController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", config.PT_HOME)
	queryRatings := map[string]bool{
		"safe":         true,
		"questionable": true,
		"sensitive":    true,
	}
	return shortcuts.Render(ctx, config.TEMPLATE_HOME, fiber.Map{
		"QueryRatings": queryRatings,
	})
}
