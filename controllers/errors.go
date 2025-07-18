package controllers

import (
	"errors"
	"imageboard/config"
	"imageboard/utils/shortcuts"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type TemplateError struct {
	Template     string
	ErrorMessage string
	StatusCode   int
}

func TemplateErrorController(ctx *fiber.Ctx, err TemplateError, bind fiber.Map) error {
	bind["Error"] = err.ErrorMessage
	return shortcuts.RenderWithStatus(ctx, err.Template, bind, err.StatusCode)
}

func GenericErrorController(ctx *fiber.Ctx, title string, err error, statusCode int) error {
	ctx.Locals("Title", title)

	if strings.HasSuffix(ctx.Path(), ".json") {
		return ctx.Status(statusCode).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if len(ctx.Path()) > 1 && strings.Contains(ctx.Path()[1:], ".") && !strings.HasSuffix(ctx.Path(), ".html") {
		return ctx.SendStatus(statusCode)
	}

	return shortcuts.RenderWithStatus(ctx, config.TEMPLATE_ERROR, fiber.Map{
		"Title": title,
		"Error": err.Error(),
	}, statusCode)
}

func NotFoundController(ctx *fiber.Ctx) error {
	error := errors.New("The page you are looking for does not exist.")
	return GenericErrorController(ctx, "Page Not Found", error, fiber.StatusNotFound)
}

func InternalServerErrorController(ctx *fiber.Ctx, err error) error {
	return GenericErrorController(ctx, "Internal Server Error", err, fiber.StatusInternalServerError)
}

func BadRequestController(ctx *fiber.Ctx, err error) error {
	return GenericErrorController(ctx, "Bad Request", err, fiber.StatusBadRequest)
}

func UnauthorizedController(ctx *fiber.Ctx, err error) error {
	return GenericErrorController(ctx, "Unauthorized", err, fiber.StatusUnauthorized)
}

func ForbiddenController(ctx *fiber.Ctx, err error) error {
	return GenericErrorController(ctx, "Forbidden", err, fiber.StatusForbidden)
}
