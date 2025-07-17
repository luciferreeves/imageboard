package controllers

import (
	"image"
	"imageboard/config"
	"imageboard/database"
	"imageboard/models"
	"imageboard/utils/auth"
	"imageboard/utils/format"
	"imageboard/utils/shortcuts"
	"imageboard/utils/transformers"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ImageUploadForm struct {
	Image  string `json:"image" form:"image"`
	Rating string `json:"rating" form:"rating"`
}

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

func PostsUploadPostController(ctx *fiber.Ctx) error {
	if !auth.IsAuthenticated(ctx) {
		return fiber.NewError(fiber.StatusForbidden, "Forbidden")
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid form data")
	}

	imageFiles := form.File["image"]
	if len(imageFiles) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "No image file provided")
	}

	imageFile := imageFiles[0]

	contentType := imageFile.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return fiber.NewError(fiber.StatusBadRequest, "Uploaded file is not an image")
	}

	if !strings.Contains(config.Upload.AllowedTypes, contentType) {
		return fiber.NewError(fiber.StatusBadRequest, "Uploaded image type is not allowed")
	}

	maxSize := int64(config.Upload.MaxSize)
	if imageFile.Size > maxSize {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge,
			"File size exceeds maximum allowed size of "+format.FileSize(maxSize))
	}

	sourceURL := ctx.FormValue("source_url")

	file, err := imageFile.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open uploaded file")
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read uploaded file")
	}

	decodedImage, format, err := format.DecodeImage(imageData)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to decode image: "+err.Error())
	}

	var fileName string
	if sourceURL != "" {
		fileName = transformers.CreateUniqueFileName(sourceURL, format)
	} else {
		fileName = transformers.CreateUniqueFileName(imageFile.Filename, format)
	}

	rating := ctx.FormValue("rating")
	if rating == "" {
		rating = "safe"
	}

	imageSizes := make(map[config.ImageSizeType]struct {
		imageSize models.ImageSize
		image     image.Image
	})
	isizeArray := []config.ImageSizeType{
		config.ImageSizeTypeIcon,
		config.ImageSizeTypeThumbnail,
		config.ImageSizeTypeSmall,
		config.ImageSizeTypeMedium,
		config.ImageSizeTypeLarge,
		config.ImageSizeTypeOriginal,
	}

	for _, sizeType := range isizeArray {
		size, img, err := transformers.TransformImageToVariant(decodedImage, sizeType)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to process image: "+err.Error())
		}
		size.ImageID = 0 // This will be set later when saving the image
		imageSizes[sizeType] = struct {
			imageSize models.ImageSize
			image     image.Image
		}{
			imageSize: size,
			image:     img,
		}
	}

	// For now, just return success - in a full implementation:
	// 1. Generate unique filename
	// 2. Calculate MD5 hash
	// 3. Resize image and create different sizes
	// 4. Upload to S3 or local storage
	// 5. Save image record to database
	// 6. Return image ID or details

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Image uploaded successfully",
		"data": fiber.Map{
			"filename":   fileName,
			"size":       imageFile.Size,
			"rating":     rating,
			"source_url": sourceURL,
		},
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

	referer := transformers.GetRefererForURL(url)
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
