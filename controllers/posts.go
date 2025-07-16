package controllers

import (
	"imageboard/config"
	"imageboard/database"
	"imageboard/utils/auth"
	"imageboard/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func PostsPageController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", config.PT_POST_LIST)
	preferences := ctx.Locals("Preferences")
	prefs, ok := preferences.(config.SitePreferences)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid preferences type")
	}

	posts, err := database.GetPosts(prefs.PostsPerPage)

	return shortcuts.Render(ctx, config.TEMPLATE_POST_LIST, fiber.Map{
		"Posts": posts,
		"Error": err,
	})
}

func PostsUploadPageController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", config.PT_POST_NEW)
	if !auth.IsAuthenticated(ctx) {
		loginURL := auth.GetLoginURLWithRedirect(ctx)
		ctx.Set("Location", loginURL)
		ctx.Status(fiber.StatusFound)
		return nil
	}

	return shortcuts.Render(ctx, config.TEMPLATE_POST_NEW, nil)
}
