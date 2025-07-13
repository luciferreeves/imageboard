package auth

import (
	"imageboard/models"

	"github.com/gofiber/fiber/v2"
)

func GetCurrentUser(ctx *fiber.Ctx) *models.User {
	if user, ok := ctx.Locals("User").(*models.User); ok {
		return user
	}
	return nil
}

func IsAuthenticated(ctx *fiber.Ctx) bool {
	return GetCurrentUser(ctx) != nil
}

func GetRedirectURL(ctx *fiber.Ctx) string {
	referer := ctx.Get("Referer")
	if referer != "" && referer != ctx.BaseURL()+"/login" && referer != ctx.BaseURL()+"/register" {
		return referer
	}
	return "/"
}
