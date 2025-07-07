package controllers

import (
	"imageboard/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func LoginController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", "Login")
	ctx.Locals("request", fiber.Map{"path": ctx.Path()})
	return shortcuts.Render(ctx, "login", nil)
}
