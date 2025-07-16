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
	preferences, ok := ctx.Locals("Preferences").(config.SitePreferences)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid preferences type")
	}

	request, ok := ctx.Locals("Request").(config.Request)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid request type")
	}

	queryTags := ""
	queryRatings := map[string]bool{}
	for _, param := range request.Query {
		switch param.Key {
		case "tags":
			queryTags = param.Value
		case "rating":
			queryRatings[param.Value] = true
		}
	}

	if len(queryRatings) == 0 {
		for _, rating := range []string{"safe", "questionable", "sensitive"} {
			queryRatings[rating] = true
		}
	}

	posts, err := database.GetPosts(preferences.PostsPerPage)

	return shortcuts.Render(ctx, config.TEMPLATE_POST_LIST, fiber.Map{
		"Posts":        posts,
		"Error":        err,
		"QueryTags":    queryTags,
		"QueryRatings": queryRatings,
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
