package transformers

import (
	"image"
	"imageboard/config"
	"imageboard/utils/format"
	"imageboard/utils/validators"
	"strings"

	"github.com/nfnt/resize"
)

func TransformImageToVariant(img image.Image, variant config.ImageSizeType) (int, int, int64, []byte, error) {
	variantSizeMap := map[config.ImageSizeType]int{
		config.ImageSizeTypeIcon:      64,
		config.ImageSizeTypeThumbnail: 256,
		config.ImageSizeTypeSmall:     512,
		config.ImageSizeTypeMedium:    1024,
		config.ImageSizeTypeLarge:     2048,
		config.ImageSizeTypeOriginal:  0,
	}

	maxWidth := variantSizeMap[variant]
	if maxWidth > 0 {
		img = ResizeImage(img, maxWidth)
	}

	fileSize, imageData, err := format.GetImageSizeAndData(img)
	if err != nil {
		return 0, 0, 0, nil, err
	}

	return img.Bounds().Dx(), img.Bounds().Dy(), fileSize, imageData, nil
}

func ResizeImage(img image.Image, maxWidth int) image.Image {
	if maxWidth <= 0 || img.Bounds().Dx() <= maxWidth {
		return img
	}

	ratio := float64(maxWidth) / float64(img.Bounds().Dx())
	newHeight := uint(float64(img.Bounds().Dy()) * ratio)

	return resize.Resize(uint(maxWidth), newHeight, img, resize.Lanczos3)
}

func CreateUniqueFileName(sourceURLOrOriginalName, imageFormat string) string {
	fileName := sourceURLOrOriginalName
	if validators.IsValidURL(sourceURLOrOriginalName) {
		parts := strings.Split(sourceURLOrOriginalName, "/")
		fileName = parts[len(parts)-1]
	}

	currentTime := format.GetCurrentTimeAsTimestamp()
	fileNameWithoutExtension := format.RemoveExtension(fileName)
	fileName = GenerateTokenFromString(fileNameWithoutExtension + "_" + format.Int64ToString(currentTime) + "_" + GenerateUUID())

	if len(fileName) > 32 {
		mid := len(fileName) / 2
		fileName = fileName[:mid-16] + format.Int64ToString(currentTime) + fileName[mid+16:]
	}
	return fileName + "." + imageFormat
}

func ConvertStringRatingToType(rating string) (config.Rating, error) {
	switch strings.ToLower(rating) {
	case "safe":
		return config.RatingSafe, nil
	case "questionable":
		return config.RatingQuestionable, nil
	case "sensitive":
		return config.RatingSensitive, nil
	case "explicit":
		return config.RatingExplicit, nil
	default:
		return config.RatingSafe, nil
	}
}

func ConvertStringToContentType(contentType string) (config.ImageContentType, error) {
	switch contentType {
	case "image/jpeg":
		return config.ImageContentTypeJPEG, nil
	case "image/png":
		return config.ImageContentTypePNG, nil
	case "image/gif":
		return config.ImageContentTypeGIF, nil
	case "image/webp":
		return config.ImageContentTypeWebP, nil
	default:
		return config.ImageContentTypeJPEG, nil
	}
}
