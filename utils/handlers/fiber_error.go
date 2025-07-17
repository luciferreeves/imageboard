package handlers

import "github.com/gofiber/fiber/v2"

func ServerErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "Internal Server Error"
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		msg = e.Message
	} else if err != nil {
		msg = err.Error()
	}
	return ctx.Status(code).SendString(msg)
}
