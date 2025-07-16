package controllers

import (
	"imageboard/config"
	"imageboard/database"
	"imageboard/session"
	"imageboard/utils/auth"
	"imageboard/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

type LoginForm struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func renderLoginError(ctx *fiber.Ctx, errorMsg string, statusCode int) error {
	return shortcuts.RenderWithStatus(ctx, config.TEMPLATE_LOGIN, fiber.Map{
		"Error":    errorMsg,
		"Username": ctx.FormValue("username"), // Preserve username in form
	}, statusCode)
}

func LoginPageController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", config.PT_LOGIN)

	if auth.IsAuthenticated(ctx) {
		return ctx.Redirect(auth.GetRedirectURL(ctx), fiber.StatusSeeOther)
	}

	next := ctx.Query("next")
	return shortcuts.Render(ctx, config.TEMPLATE_LOGIN, fiber.Map{
		"Next": next,
	})
}

func LoginPostController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", config.PT_LOGIN)

	var form LoginForm
	var err error
	if err = ctx.BodyParser(&form); err != nil {
		return renderLoginError(ctx, config.ERR_INVALID_FORM_DATA, fiber.StatusBadRequest)
	}

	user, err := database.GetUserByUsername(form.Username)
	if err != nil {
		return renderLoginError(ctx, config.ERR_USER_NOT_FOUND, fiber.StatusUnauthorized)
	}

	if !user.CheckPassword(form.Password) {
		return renderLoginError(ctx, config.ERR_LOGIN_INVALID_CREDENTIALS, fiber.StatusUnauthorized)
	}

	if !user.IsActive() {
		return renderLoginError(ctx, config.ERR_ACCOUNT_DISABLED, fiber.StatusForbidden)
	}

	if !user.CanLogin() {
		return renderLoginError(ctx, config.ERR_ACCOUNT_UNABLE_TO_LOGIN, fiber.StatusForbidden)
	}

	sess, err := session.Store.Get(ctx)
	if err != nil {
		return renderLoginError(ctx, config.ERR_SESSION_FAILED_TO_CREATE, fiber.StatusInternalServerError)
	}
	sess.Set("user_id", user.ID)
	sess.Set("username", user.Username)
	if err := sess.Save(); err != nil {
		return renderLoginError(ctx, config.ERR_SESSION_FAILED_TO_SAVE, fiber.StatusInternalServerError)
	}

	user.UpdateLastUserLogin(database.DB)
	user.UpdateLastUserActivity(database.DB)

	return ctx.Redirect(auth.GetRedirectURL(ctx), fiber.StatusSeeOther)
}
