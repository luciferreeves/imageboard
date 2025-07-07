package controllers

import (
	"imageboard/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func HomeController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", "Home Page")
	return shortcuts.Render(ctx, "home", nil)
}
