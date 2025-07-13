package controllers

import (
	"imageboard/session"
	"imageboard/utils/auth"

	"github.com/gofiber/fiber/v2"
)

func LogoutController(ctx *fiber.Ctx) error {
	sess, err := session.Store.Get(ctx)
	if err != nil {
		return ctx.Redirect(auth.GetRedirectURL(ctx), fiber.StatusSeeOther)
	}

	if err := sess.Destroy(); err != nil {
		sess.Delete("user_id")
		sess.Delete("username")
		sess.Save()
	}

	return ctx.Redirect(auth.GetRedirectURL(ctx), fiber.StatusSeeOther)
}
