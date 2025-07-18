package controllers

import (
	"imageboard/config"
	"imageboard/database"
	"imageboard/models"
	"imageboard/utils/auth"
	"imageboard/utils/email"
	"imageboard/utils/shortcuts"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type RegisterForm struct {
	Username        string `json:"username" form:"username"`
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
}

func RegisterPageController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", config.PT_REGISTER)

	if auth.IsAuthenticated(ctx) {
		return ctx.Redirect(auth.GetRedirectURL(ctx), fiber.StatusSeeOther)
	}

	return shortcuts.Render(ctx, config.TEMPLATE_REGISTER, nil)
}

func RegisterPostController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", config.PT_REGISTER)

	if auth.IsAuthenticated(ctx) {
		return ctx.Redirect(auth.GetRedirectURL(ctx), fiber.StatusSeeOther)
	}

	var form RegisterForm
	handleRegisterError := func(errorMessage string, statusCode int) error {
		return TemplateErrorController(ctx, TemplateError{
			Template:     config.TEMPLATE_REGISTER,
			ErrorMessage: errorMessage,
			StatusCode:   statusCode,
		}, fiber.Map{
			"Username": form.Username,
			"Email":    form.Email,
		})
	}

	if err := ctx.BodyParser(&form); err != nil {
		return handleRegisterError(config.ERR_INVALID_FORM_DATA, fiber.StatusBadRequest)
	}

	if form.Password != form.ConfirmPassword {
		return handleRegisterError(config.ERR_PASSWORD_MISMATCH, fiber.StatusBadRequest)
	}

	user := &models.User{
		Username: form.Username,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := database.CreateUser(user); err != nil {
		var statusCode int
		if strings.Contains(err.Error(), "username") {
			statusCode = fiber.StatusConflict
		} else if strings.Contains(err.Error(), "email") {
			statusCode = fiber.StatusBadRequest
		} else {
			statusCode = fiber.StatusInternalServerError
		}

		return handleRegisterError(config.ERR_REGISTER_FAILED_TO_CREATE_USER+err.Error(), statusCode)
	}

	if err := email.SendVerificationEmail(user); err != nil {
		log.Printf("Failed to send verification email: %v", err)
		return handleRegisterError(config.ERR_REGISTER_USER_CREATED_EMAIL_FAILED, fiber.StatusInternalServerError)
	}

	return shortcuts.Render(ctx, config.TEMPLATE_REGISTER, fiber.Map{
		"Success": config.SUCCESS_USER_REGISTERED,
	})
}
