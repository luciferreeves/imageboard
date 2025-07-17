package controllers

import (
	"imageboard/config"
	"imageboard/database"
	"imageboard/utils/auth"
	"imageboard/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func renderVerifyEmailError(ctx *fiber.Ctx, errorMsg string, statusCode int) error {
	return shortcuts.RenderWithStatus(ctx, config.TEMPLATE_VERIFY_EMAIL, fiber.Map{
		"Error": errorMsg,
	}, statusCode)
}

func VerifyEmailController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", config.PT_VERIFY_EMAIL)
	if auth.IsAuthenticated(ctx) {
		return ctx.Redirect(auth.GetRedirectURL(ctx), fiber.StatusSeeOther)
	}
	token := ctx.Query("token")
	if token == "" {
		return renderVerifyEmailError(ctx, config.ERR_VERIFY_EMAIL_MISSING_TOKEN, fiber.StatusBadRequest)
	}

	emailToken, err := database.VerifyToken(token, config.EmailTokenTypeVerification)
	if err != nil {
		return renderVerifyEmailError(ctx, config.ERR_VERIFY_EMAIL_INVALID_OR_EXPIRED_TOKEN, fiber.StatusBadRequest)
	}

	user, err := database.GetUserByID(emailToken.UserID)
	if err != nil {
		return renderVerifyEmailError(ctx, config.ERR_VERIFY_EMAIL_USER_NOT_FOUND, fiber.StatusInternalServerError)
	}

	user.Activate()
	if err := database.DB.Save(user).Error; err != nil {
		return renderVerifyEmailError(ctx, config.ERR_VERIFY_EMAIL_ACTIVATION_FAILED, fiber.StatusInternalServerError)
	}

	return shortcuts.Render(ctx, config.TEMPLATE_VERIFY_EMAIL, fiber.Map{
		"Success":  config.SUCCESS_VERIFY_EMAIL,
		"Username": user.Username,
	})

}
