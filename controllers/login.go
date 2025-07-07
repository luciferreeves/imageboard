package controllers

import (
	"imageboard/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func LoginController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", "Login")
	return shortcuts.Render(ctx, "login", nil)
}
