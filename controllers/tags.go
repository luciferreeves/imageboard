package controllers

import (
	"imageboard/config"
	"imageboard/database"
	"imageboard/models"
	"imageboard/utils/auth"
	"imageboard/utils/format"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func TagsSearchJSONController(ctx *fiber.Ctx) error {
	tagName := ctx.Query("name")
	tagType := config.TagType(ctx.Query("type"))
	limit := ctx.QueryInt("limit", 20)
	offset := ctx.QueryInt("offset", 0)
	if tagName == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tag name is required",
		})
	}

	tags, err := database.SearchTags(tagName, limit, offset, &tagType)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tags",
		})
	}

	if len(tags) == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No tags found",
		})
	}

	return ctx.JSON(tags)
}

func FindOrCreateTagJSONController(ctx *fiber.Ctx) error {
	var request struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	if !auth.IsAuthenticated(ctx) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	currentUser := auth.GetCurrentUser(ctx)
	if !currentUser.CanCreateTags() {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to create tags",
		})
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	tag, err := database.FindOrCreateTag(request.Name, config.TagType(request.Type))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find or create tag",
		})
	}

	return ctx.JSON(tag)
}

func TagsSearchForImageJSONController(ctx *fiber.Ctx) error {
	if !auth.IsAuthenticated(ctx) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	query := ctx.Query("q")
	imageID := ctx.Query("image_id")
	tagType := ctx.Query("type")

	if query == "" {
		return ctx.JSON([]interface{}{})
	}

	var uintImageID uint
	if imageID != "" {
		if id, err := strconv.ParseUint(imageID, 10, 32); err == nil {
			uintImageID = uint(id)
		}
	}

	var tagTypeEnum *config.TagType
	if tagType != "" {
		t := config.TagType(tagType)
		tagTypeEnum = &t
	}

	var tags []models.Tag
	var err error

	if uintImageID > 0 {
		tags, err = database.SearchTagsExcluding(query, uintImageID, 10, tagTypeEnum)
	} else {
		tags, err = database.SearchTags(query, 10, 0, tagTypeEnum)
	}

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Search failed",
		})
	}

	currentUser := auth.GetCurrentUser(ctx)
	canCreate := currentUser != nil && currentUser.CanCreateTags()

	result := fiber.Map{
		"tags":       tags,
		"can_create": canCreate,
		"query":      query,
	}

	return ctx.JSON(result)
}

func TagsAddToImageJSONController(ctx *fiber.Ctx) error {
	if !auth.IsAuthenticated(ctx) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	imageID := ctx.Query("image_id")
	if imageID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image ID required",
		})
	}

	uintImageID, err := format.StringToUint(imageID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid image ID",
		})
	}

	var request struct {
		TagName string `json:"tag_name"`
		TagType string `json:"tag_type"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if request.TagName == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tag name required",
		})
	}

	if request.TagType == "" {
		request.TagType = "general"
	}

	tagTypeEnum := config.TagType(request.TagType)
	currentUser := auth.GetCurrentUser(ctx)

	tag, err := database.FindOrCreateTag(request.TagName, tagTypeEnum)
	if err != nil {
		// Check if this is a cross-category error
		if strings.Contains(err.Error(), "already exists as") || strings.Contains(err.Error(), "previously existed as") {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		// For other errors, check permissions
		if !currentUser.CanCreateTags() {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Cannot create new tags",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create tag",
		})
	}

	var tagsToAdd []uint
	if tag.ParentID != nil {
		_, ancestors, err := database.GetTagWithAncestors(tag.ID)
		if err == nil {
			for _, ancestor := range ancestors {
				tagsToAdd = append(tagsToAdd, ancestor.ID)
			}
		}
	}
	tagsToAdd = append(tagsToAdd, tag.ID)

	for _, tagID := range tagsToAdd {
		database.AddTagToImage(uintImageID, tagID)
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"tag":     tag,
	})
}

func TagsRemoveFromImageJSONController(ctx *fiber.Ctx) error {
	if !auth.IsAuthenticated(ctx) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	imageID := ctx.Query("image_id")
	if imageID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image ID required",
		})
	}

	uintImageID, err := format.StringToUint(imageID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid image ID",
		})
	}

	var request struct {
		TagID uint `json:"tag_id"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	currentUser := auth.GetCurrentUser(ctx)
	post, err := database.GetPostByID(uintImageID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	if post.Uploader.Username != currentUser.Username && !currentUser.CanEditPosts() {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot edit this post",
		})
	}

	_, descendants, err := database.GetTagWithDescendants(request.TagID)
	if err == nil {
		for _, descendant := range descendants {
			database.RemoveTagFromImage(uintImageID, descendant.ID)
		}
	}

	err = database.RemoveTagFromImage(uintImageID, request.TagID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove tag",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
	})
}
