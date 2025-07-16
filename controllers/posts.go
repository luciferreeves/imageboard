package controllers

import (
	"imageboard/config"
	"imageboard/database"
	"imageboard/utils/auth"
	"imageboard/utils/format"
	"imageboard/utils/shortcuts"
	"imageboard/utils/validators"
	"io"
	"net/http"
	"strings"

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

	allowedTypes := []string{}
	for t := range strings.SplitSeq(config.Upload.AllowedTypes, ",") {
		if idx := strings.Index(t, "/"); idx != -1 && idx+1 < len(t) {
			subtype := t[idx+1:]
			if subtype != "" {
				allowedTypes = append(allowedTypes, "."+subtype)
			}
		}
	}

	return shortcuts.Render(ctx, config.TEMPLATE_POST_NEW, fiber.Map{
		"AllowedTypes": allowedTypes,
		"MaxSize":      format.FileSize(int64(config.Upload.MaxSize)),
	})
}

func PostsUploadImageLinkProxyController(ctx *fiber.Ctx) error {
	maxSize := int64(config.Upload.MaxSize)
	if !auth.IsAuthenticated(ctx) {
		return fiber.NewError(fiber.StatusForbidden, "Forbidden")
	}

	url := ctx.Query("url")
	if url == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing url parameter")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid URL")
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")

	referer := validators.GetRefererForURL(url)
	if referer != "" {
		req.Header.Set("Referer", referer)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, "Failed to fetch image")
	}
	if resp.StatusCode != 200 {
		return fiber.NewError(fiber.StatusBadGateway, "Failed to fetch image")
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return fiber.NewError(fiber.StatusBadRequest, "URL does not point to an image")
	}

	ctx.Set("Content-Type", contentType)
	ctx.Set("Cache-Control", "no-store")
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read image data")
	}
	if int64(len(buf)) > maxSize {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "Image exceeds maximum allowed size of "+format.FileSize(maxSize))
	}
	return ctx.Send(buf)
}
