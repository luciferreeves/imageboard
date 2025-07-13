package controllers

import (
	"imageboard/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func HomePageController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", PT_HOME)
	return shortcuts.Render(ctx, TEMPLATE_HOME, nil)
}
