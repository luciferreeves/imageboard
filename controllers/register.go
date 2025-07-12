package controllers

import (
	"imageboard/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func RegisterPageController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", "Register")
	return shortcuts.Render(ctx, "register", nil)
}
