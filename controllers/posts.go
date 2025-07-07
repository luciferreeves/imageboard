package controllers

import (
	"imageboard/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func PostsController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", "Posts")

	searchQuery := ctx.Query("tags", "")

	customdata := struct {
		SearchQuery string
		Posts       []interface{}
	}{
		SearchQuery: searchQuery,
		Posts:       []interface{}{},
	}
	return shortcuts.Render(ctx, "posts", customdata)
}
