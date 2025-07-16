package controllers

import (
	"imageboard/config"
	"imageboard/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func HomePageController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", config.PT_HOME)
	return shortcuts.Render(ctx, config.TEMPLATE_HOME, nil)
}
