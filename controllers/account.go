package controllers

import (
	"imageboard/config"
	"imageboard/database"
	"imageboard/models"
	"imageboard/session"
	"imageboard/utils/auth"
	"imageboard/utils/email"
	"imageboard/utils/shortcuts"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type LoginForm struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
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
	handleLoginError := func(errorMessage string, statusCode int) error {
		return TemplateErrorController(ctx, TemplateError{
			Template:     config.TEMPLATE_LOGIN,
			ErrorMessage: errorMessage,
			StatusCode:   statusCode,
		}, fiber.Map{
			"Username": form.Username,
		})
	}

	if err = ctx.BodyParser(&form); err != nil {
		return handleLoginError(config.ERR_INVALID_FORM_DATA, fiber.StatusBadRequest)
	}

	user, err := database.GetUserByUsername(form.Username)
	if err != nil {
		return handleLoginError(config.ERR_USER_NOT_FOUND, fiber.StatusUnauthorized)
	}

	if !user.CheckPassword(form.Password) {
		return handleLoginError(config.ERR_LOGIN_INVALID_CREDENTIALS, fiber.StatusUnauthorized)
	}

	if !user.IsActive() {
		return handleLoginError(config.ERR_ACCOUNT_DISABLED, fiber.StatusForbidden)
	}

	if !user.CanLogin() {
		return handleLoginError(config.ERR_ACCOUNT_UNABLE_TO_LOGIN, fiber.StatusForbidden)
	}

	sess, err := session.Store.Get(ctx)
	if err != nil {
		return handleLoginError(config.ERR_SESSION_FAILED_TO_CREATE, fiber.StatusInternalServerError)
	}
	sess.Set("user_id", user.ID)
	sess.Set("username", user.Username)
	if err := sess.Save(); err != nil {
		return handleLoginError(config.ERR_SESSION_FAILED_TO_SAVE, fiber.StatusInternalServerError)
	}

	user.UpdateLastUserLogin(database.DB)
	user.UpdateLastUserActivity(database.DB)

	return ctx.Redirect(auth.GetRedirectURL(ctx), fiber.StatusSeeOther)
}

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
		Username:             form.Username,
		Email:                form.Email,
		Password:             form.Password,
		PostsRequireApproval: true,
		Level:                config.UserLevelMember,
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

func ForgotPasswordPageController(ctx *fiber.Ctx) error {
	mode := ctx.Query("mode", "username")
	switch mode {
	case "username":
		ctx.Locals("Title", config.PT_FORGOT_USERNAME)
	case "password":
		ctx.Locals("Title", config.PT_FORGOT_PASSWORD)
	default:
		ctx.Locals("Title", config.PT_FORGOT_USERNAME)
		mode = "username"
	}

	if auth.IsAuthenticated(ctx) {
		return ctx.Redirect(auth.GetRedirectURL(ctx), fiber.StatusSeeOther)
	}

	return shortcuts.Render(ctx, config.TEMPLATE_FORGOT, fiber.Map{
		"Mode": mode,
	})
}

type ForgotPasswordInput struct {
	Email string `json:"email" form:"email"`
	Mode  string `json:"mode" form:"mode"`
}

func ForgotPasswordPostController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", config.PT_FORGOT_PASSWORD)

	if auth.IsAuthenticated(ctx) {
		return ctx.Redirect(auth.GetRedirectURL(ctx), fiber.StatusSeeOther)
	}

	var input ForgotPasswordInput
	if err := ctx.BodyParser(&input); err != nil {
		return TemplateErrorController(ctx, TemplateError{
			Template:     config.TEMPLATE_FORGOT,
			ErrorMessage: config.ERR_INVALID_FORM_DATA,
			StatusCode:   fiber.StatusBadRequest,
		}, fiber.Map{
			"Mode": input.Mode,
		})
	}

	switch input.Mode {
	case "password":
		ctx.Locals("Title", config.PT_FORGOT_PASSWORD)
	case "username":
		ctx.Locals("Title", config.PT_FORGOT_USERNAME)
	default:
		ctx.Locals("Title", config.PT_FORGOT_USERNAME)
		input.Mode = "username"
	}

	users, err := database.GetUsersByEmail(input.Email)
	if err != nil || len(users) == 0 {
		return TemplateErrorController(ctx, TemplateError{
			Template:     config.TEMPLATE_FORGOT,
			ErrorMessage: config.ERR_NO_ACCOUNT_ASSOCIATED_WITH_EMAIL,
			StatusCode:   fiber.StatusNotFound,
		}, fiber.Map{
			"Mode": input.Mode,
		})
	}

	switch mode := input.Mode; mode {
	case "username":
		if err := email.SendForgotUsernameEmail(&users); err != nil {
			log.Printf("Failed to send forgot username email: %v", err)
			return TemplateErrorController(ctx, TemplateError{
				Template:     config.TEMPLATE_FORGOT,
				ErrorMessage: "Failed to send username email. Please try again later.",
				StatusCode:   fiber.StatusInternalServerError,
			}, fiber.Map{
				"Mode": input.Mode,
			})
		}
	case "password":
		// TODO
	default:
		return TemplateErrorController(ctx, TemplateError{
			Template:     config.TEMPLATE_FORGOT,
			ErrorMessage: config.ERR_INVALID_FORM_DATA,
			StatusCode:   fiber.StatusBadRequest,
		}, fiber.Map{
			"Mode": input.Mode,
		})
	}

	switch input.Mode {
	case "username":
		return shortcuts.Render(ctx, config.TEMPLATE_FORGOT, fiber.Map{
			"Success": config.SUCCESS_FORGOT_USERNAME_EMAIL_SENT,
			"Mode":    input.Mode,
		})
	case "password":
		// TODO
		return shortcuts.Render(ctx, config.TEMPLATE_FORGOT, fiber.Map{
			"Success": "If an account with that email exists, a password reset email has been sent.",
			"Mode":    input.Mode,
		})
	default:
		return shortcuts.Render(ctx, config.TEMPLATE_FORGOT, fiber.Map{
			"Success": config.SUCCESS_FORGOT_USERNAME_EMAIL_SENT,
			"Mode":    input.Mode,
		})
	}
}

func VerifyEmailController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", config.PT_VERIFY_EMAIL)
	if auth.IsAuthenticated(ctx) {
		return ctx.Redirect(auth.GetRedirectURL(ctx), fiber.StatusSeeOther)
	}
	token := ctx.Query("token")
	handleVerifyEmailError := func(errorMessage string, statusCode int) error {
		return TemplateErrorController(ctx, TemplateError{
			Template:     config.TEMPLATE_VERIFY_EMAIL,
			ErrorMessage: errorMessage,
			StatusCode:   statusCode,
		}, nil)
	}

	if token == "" {
		return handleVerifyEmailError(config.ERR_VERIFY_EMAIL_MISSING_TOKEN, fiber.StatusBadRequest)
	}

	emailToken, err := database.VerifyToken(token, config.EmailTokenTypeVerification)
	if err != nil {
		return handleVerifyEmailError(config.ERR_VERIFY_EMAIL_INVALID_OR_EXPIRED_TOKEN, fiber.StatusBadRequest)
	}

	user, err := database.GetUserByID(emailToken.UserID)
	if err != nil {
		if err.Error() == "record not found" {
			return handleVerifyEmailError(config.ERR_VERIFY_EMAIL_USER_NOT_FOUND, fiber.StatusBadRequest)
		}

		return handleVerifyEmailError(config.ERR_VERIFY_EMAIL_ACTIVATION_FAILED, fiber.StatusInternalServerError)
	}

	user.Activate()
	if err := database.DB.Save(user).Error; err != nil {
		if err.Error() == "record not found" {
			return handleVerifyEmailError(config.ERR_VERIFY_EMAIL_USER_NOT_FOUND, fiber.StatusBadRequest)
		}

		return handleVerifyEmailError(config.ERR_VERIFY_EMAIL_ACTIVATION_FAILED, fiber.StatusInternalServerError)
	}

	return shortcuts.Render(ctx, config.TEMPLATE_VERIFY_EMAIL, fiber.Map{
		"Success":  config.SUCCESS_VERIFY_EMAIL,
		"Username": user.Username,
	})

}
