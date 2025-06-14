package middleware

import "github.com/gofiber/fiber/v2"

func JSON(context *fiber.Ctx) error {
	context.Accepts("application/json")
	return context.Next()
}
