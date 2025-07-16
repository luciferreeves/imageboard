package auth

import (
	"imageboard/config"
	"imageboard/models"
	"net/url"

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
	next := ctx.Query("next")
	if next == "" {
		next = ctx.FormValue("next")
	}
	if next != "" && isValidRedirectURL(next) {
		return next
	}
	return config.URL_HOME
}

func isValidRedirectURL(redirectURL string) bool {
	if redirectURL == "" {
		return false
	}

	if redirectURL == config.URL_LOGIN || redirectURL == config.URL_REGISTER || redirectURL == config.URL_LOGOUT {
		return false
	}

	if redirectURL[0] == '/' {
		return true
	}

	return false
}

func GetLoginURLWithRedirect(ctx *fiber.Ctx) string {
	currentPath := ctx.Path()
	if queryString := string(ctx.Request().URI().QueryString()); queryString != "" {
		currentPath += "?" + queryString
	}
	return config.URL_LOGIN + "?next=" + url.QueryEscape(currentPath)
}

func GetLogoutURLWithRedirect(ctx *fiber.Ctx) string {
	currentPath := ctx.Path()
	if queryString := string(ctx.Request().URI().QueryString()); queryString != "" {
		currentPath += "?" + queryString
	}
	return config.URL_LOGOUT + "?next=" + url.QueryEscape(currentPath)
}
