package controllers

import (
	"imageboard/database"
	"imageboard/session"
	"imageboard/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func getRedirectURL(ctx *fiber.Ctx) string {
	referer := ctx.Get("Referer")
	if referer != "" && referer != ctx.BaseURL()+URL_LOGIN && referer != ctx.BaseURL()+URL_REGISTER {
		return referer
	}
	return URL_HOME
}

func renderLoginError(ctx *fiber.Ctx, errorMsg string, statusCode int) error {
	return shortcuts.RenderWithStatus(ctx, TEMPLATE_LOGIN, fiber.Map{
		"Error":    errorMsg,
		"Username": ctx.FormValue("username"), // Preserve username in form
	}, statusCode)
}

func LoginPageController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", PT_LOGIN)
	sess, err := session.Store.Get(ctx)
	if err == nil {
		if userID, ok := sess.Get("user_id").(int); ok && userID != 0 {
			return ctx.Redirect(getRedirectURL(ctx), fiber.StatusSeeOther)
		}
	}

	return shortcuts.Render(ctx, TEMPLATE_LOGIN, nil)
}

func LoginPostController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", PT_LOGIN)
	type LoginForm struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}

	var form LoginForm
	var err error
	if err = ctx.BodyParser(&form); err != nil {
		return renderLoginError(ctx, ERR_INVALID_FORM_DATA, fiber.StatusBadRequest)
	}

	user, err := database.GetUserByUsername(form.Username)
	if err != nil {
		return renderLoginError(ctx, ERR_USER_NOT_FOUND, fiber.StatusUnauthorized)
	}

	if !user.CheckPassword(form.Password) {
		return renderLoginError(ctx, ERR_LOGIN_INVALID_CREDENTIALS, fiber.StatusUnauthorized)
	}

	if !user.IsActive() {
		return renderLoginError(ctx, ERR_ACCOUNT_DISABLED, fiber.StatusForbidden)
	}

	if !user.CanLogin() {
		return renderLoginError(ctx, ERR_ACCOUNT_UNABLE_TO_LOGIN, fiber.StatusForbidden)
	}

	sess, err := session.Store.Get(ctx)
	if err != nil {
		return renderLoginError(ctx, ERR_SESSION_FAILED_TO_CREATE, fiber.StatusInternalServerError)
	}
	sess.Set("user_id", user.ID)
	sess.Set("username", user.Username)
	if err := sess.Save(); err != nil {
		return renderLoginError(ctx, ERR_SESSION_FAILED_TO_SAVE, fiber.StatusInternalServerError)
	}

	user.UpdateLastUserLogin(database.DB)
	user.UpdateLastUserActivity(database.DB)

	return ctx.Redirect(getRedirectURL(ctx), fiber.StatusSeeOther)
}
