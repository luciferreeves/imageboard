package controllers

import (
	"imageboard/config"
	"imageboard/session"

	"github.com/gofiber/fiber/v2"
)

func LogoutController(ctx *fiber.Ctx) error {
	sess, err := session.Store.Get(ctx)
	if err != nil {
		return ctx.Redirect(config.URL_HOME, fiber.StatusSeeOther)
	}

	if err := sess.Destroy(); err != nil {
		sess.Delete("user_id")
		sess.Delete("username")
		sess.Save()
	}

	next := ctx.Query("next")
	if next != "" {
		return ctx.Redirect(next, fiber.StatusSeeOther)
	}

	return ctx.Redirect(config.URL_HOME, fiber.StatusSeeOther)
}
