package processors

import (
	"imageboard/config"

	"github.com/gofiber/fiber/v2"
)

const defaultTitle = "default"

func MetaContextProcessor(ctx *fiber.Ctx) error {
	ctx.Locals("Title", defaultTitle)
	ctx.Locals("Appname", config.Server.AppName)
	return ctx.Next()
}
