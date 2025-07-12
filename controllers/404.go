package controllers

import (
	"imageboard/utils/shortcuts"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NotFoundController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", "Page Not Found")

	path := ctx.Path()

	if strings.HasSuffix(path, ".json") {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Not Found",
		})
	}

	if len(path) > 1 && strings.Contains(path[1:], ".") && !strings.HasSuffix(path, ".html") {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return shortcuts.Render(ctx, "404", nil)
}
