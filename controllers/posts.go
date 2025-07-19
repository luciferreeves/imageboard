package controllers

import (
	"errors"
	"imageboard/config"
	"imageboard/database"
	"imageboard/utils/auth"
	"imageboard/utils/format"
	"imageboard/utils/handlers"
	"imageboard/utils/minio"
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

	queryTags, queryTagsList := handlers.ExtractQueryTags(request.Query)
	queryRatings, queryRatingsMap := handlers.ExtractRatingsAndMap(request.Query)

	posts, err := database.GetPosts(preferences.PostsPerPage, queryRatings, queryTagsList)
	if err != nil {
		return InternalServerErrorController(ctx, err)
	}

	return shortcuts.Render(ctx, config.TEMPLATE_POST_LIST, fiber.Map{
		"Posts":        posts,
		"QueryTags":    queryTags,
		"QueryRatings": queryRatingsMap,
		"CDNURL":       format.GetCDNURL(),
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

	decodedImage, imageFormat, err := format.DecodeImage(imageData)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to decode image: "+err.Error())
	}

	var fileName string
	if sourceURL != "" {
		fileName = transformers.CreateUniqueFileName(sourceURL, imageFormat)
	} else {
		fileName = transformers.CreateUniqueFileName(imageFile.Filename, imageFormat)
	}

	rating := ctx.FormValue("rating")
	if rating == "" {
		rating = "safe"
	}

	isizeArray := []config.ImageSizeType{
		config.ImageSizeTypeIcon,
		config.ImageSizeTypeThumbnail,
		config.ImageSizeTypeSmall,
		config.ImageSizeTypeMedium,
		config.ImageSizeTypeLarge,
		config.ImageSizeTypeOriginal,
	}

	currentUser := auth.GetCurrentUser(ctx)
	if currentUser == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	md5Hash := transformers.GenerateMD5Hash(imageData)

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to start transaction: "+tx.Error.Error())
	}

	dbImage, err := database.CreateImageWithTx(tx,
		fileName,
		contentType,
		md5Hash,
		sourceURL,
		rating,
		currentUser.ID,
		currentUser.PostsRequireApproval,
	)
	if err != nil {
		tx.Rollback()
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create image record: "+err.Error())
	}

	type imageSizeData struct {
		sizeType config.ImageSizeType
		width    int
		height   int
		fileSize int64
	}
	var imageSizes []imageSizeData

	for _, sizeType := range isizeArray {
		width, height, fileSize, processedImageData, err := transformers.TransformImageToVariant(decodedImage, sizeType)
		if err != nil {
			tx.Rollback()
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to process image: "+err.Error())
		}

		err = minio.UploadImage(processedImageData, sizeType, fileName, contentType)
		if err != nil {
			tx.Rollback()
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload image: "+err.Error())
		}

		imageSizes = append(imageSizes, imageSizeData{
			sizeType: sizeType,
			width:    width,
			height:   height,
			fileSize: fileSize,
		})
	}

	for _, sizeData := range imageSizes {
		_, err = database.CreateImageSizeWithTx(tx, dbImage.ID, sizeData.sizeType, sizeData.width, sizeData.height, sizeData.fileSize)
		if err != nil {
			tx.Rollback()
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to create image size record: "+err.Error())
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction: "+err.Error())
	}

	return ctx.SendStatus(fiber.StatusOK)
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

func PostsSinglePostPageController(ctx *fiber.Ctx) error {
	ctx.Locals("Title", config.PT_POST_SINGLE)

	postID := ctx.Params("id")
	if postID == "" {
		return NotFoundController(ctx)
	}

	uintPostID, err := format.StringToUint(postID)
	if err != nil {
		return NotFoundController(ctx)
	}

	post, err := database.GetPostByID(uintPostID)
	if err != nil {
		if err.Error() == "record not found" {
			return NotFoundController(ctx)
		}
		return InternalServerErrorController(ctx, err)
	}

	currentUser := auth.GetCurrentUser(ctx)
	isUserFavourited := false
	if currentUser != nil {
		isUserFavourited = post.IsUserFavourited(database.DB, currentUser)
	}

	ctx.Locals("Title", config.PT_POST_SINGLE+" #"+format.Int64ToString(int64(post.ID)))
	return shortcuts.Render(ctx, config.TEMPLATE_POST_SINGLE, fiber.Map{
		"Post":             post,
		"CDNURL":           format.GetCDNURL(),
		"IsUserFavourited": isUserFavourited,
	})
}

func PostsSinglePostFavouritePostController(ctx *fiber.Ctx) error {
	if !auth.IsAuthenticated(ctx) {
		return ctx.Redirect(auth.GetLoginURLWithNextField(ctx), fiber.StatusFound)
	}

	postID := ctx.Params("id")
	if postID == "" {
		return NotFoundController(ctx)
	}

	uintPostID, err := format.StringToUint(postID)
	if err != nil {
		return NotFoundController(ctx)
	}

	post, err := database.GetPostByID(uintPostID)
	if err != nil {
		if err.Error() == "record not found" {
			return NotFoundController(ctx)
		}
		return InternalServerErrorController(ctx, err)
	}

	currentUser := auth.GetCurrentUser(ctx)
	if currentUser == nil {
		return UnauthorizedController(ctx, errors.New("User not found"))
	}

	if err := post.ToggleFavourite(database.DB, currentUser); err != nil {
		return InternalServerErrorController(ctx, err)
	}

	return ctx.Redirect(auth.GetRedirectURL(ctx), fiber.StatusSeeOther)
}

func PostsSinglePostEditPageController(ctx *fiber.Ctx) error {
	if !auth.IsAuthenticated(ctx) {
		return ctx.Redirect(auth.GetLoginURLWithNextField(ctx), fiber.StatusFound)
	}

	postID := ctx.Params("id")
	if postID == "" {
		return NotFoundController(ctx)
	}

	uintPostID, err := format.StringToUint(postID)
	if err != nil {
		return NotFoundController(ctx)
	}

	post, err := database.GetPostByID(uintPostID)
	if err != nil {
		if err.Error() == "record not found" {
			return NotFoundController(ctx)
		}
		return InternalServerErrorController(ctx, err)
	}

	currentUser := auth.GetCurrentUser(ctx)
	if post.Uploader.Username != currentUser.Username && !currentUser.CanEditPosts() {
		return ForbiddenController(ctx, errors.New("You do not have permission to edit this post"))
	}

	users, err := database.ListAllUsers()
	if err != nil {
		return InternalServerErrorController(ctx, err)
	}
	approvers, err := database.ListAllApprovers()
	if err != nil {
		return InternalServerErrorController(ctx, err)
	}

	postTags, err := database.GetImageTags(post.ID)
	if err != nil {
		return InternalServerErrorController(ctx, err)
	}

	ctx.Locals("Title", config.PT_POST_EDIT+" #"+format.Int64ToString(int64(post.ID)))
	return shortcuts.Render(ctx, config.TEMPLATE_POST_EDIT, fiber.Map{
		"Post":      post,
		"CDNURL":    format.GetCDNURL(),
		"Users":     users,
		"Approvers": approvers,
		"PostTags":  postTags,
	})
}

func PostsSinglePostEditPostController(ctx *fiber.Ctx) error {
	if !auth.IsAuthenticated(ctx) {
		return ctx.Redirect(auth.GetLoginURLWithNextField(ctx), fiber.StatusFound)
	}

	postID := ctx.Params("id")
	if postID == "" {
		return NotFoundController(ctx)
	}

	uintPostID, err := format.StringToUint(postID)
	if err != nil {
		return NotFoundController(ctx)
	}

	post, err := database.GetPostByID(uintPostID)
	if err != nil {
		if err.Error() == "record not found" {
			return NotFoundController(ctx)
		}
		return InternalServerErrorController(ctx, err)
	}

	currentUser := auth.GetCurrentUser(ctx)
	if post.Uploader.Username != currentUser.Username && !currentUser.CanEditPosts() {
		return ForbiddenController(ctx, errors.New("You do not have permission to edit this post"))
	}

	title := ctx.FormValue("title")
	description := ctx.FormValue("description")
	sourceURL := ctx.FormValue("source_url")
	rating := ctx.FormValue("rating")

	updates := make(map[string]interface{})

	if title != post.Title {
		updates["title"] = title
	}

	if description != post.Description {
		updates["description"] = description
	}

	if sourceURL != post.SourceURL {
		updates["source_url"] = sourceURL
	}

	if rating != "" && rating != string(post.Rating) {
		ratingEnum, err := transformers.ConvertStringRatingToType(rating)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid rating value")
		}
		updates["rating"] = ratingEnum
	}

	if currentUser.CanApprovePosts() {
		isApproved := ctx.FormValue("is_approved") == "1"
		if isApproved != post.IsApproved {
			updates["is_approved"] = isApproved
		}
	}

	if currentUser.CanDeletePosts() {
		isDeleted := ctx.FormValue("is_deleted") == "1"
		if isDeleted != post.IsDeleted {
			updates["is_deleted"] = isDeleted
		}
	}

	if currentUser.IsAdmin() {
		uploaderID := ctx.FormValue("uploader")
		if uploaderID != "" {
			uintUploaderID, err := format.StringToUint(uploaderID)
			if err == nil && uintUploaderID != post.UploaderID {
				updates["uploader_id"] = uintUploaderID
			}
		}

		approverID := ctx.FormValue("approver")
		if approverID != "" {
			if approverID == "0" {
				if post.ApproverID != nil {
					updates["approver_id"] = nil
				}
			} else {
				uintApproverID, err := format.StringToUint(approverID)
				if err == nil && (post.ApproverID == nil || *post.ApproverID != uintApproverID) {
					updates["approver_id"] = uintApproverID
				}
			}
		}
	}

	if len(updates) > 0 {
		if err := database.UpdateImage(post.ID, updates); err != nil {
			return InternalServerErrorController(ctx, err)
		}
	}

	nextURL := ctx.FormValue("next")
	if nextURL == "" {
		nextURL = "/posts/" + format.Int64ToString(int64(post.ID))
	}
	return ctx.Redirect(nextURL, fiber.StatusSeeOther)
}
