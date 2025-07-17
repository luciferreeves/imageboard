package transformers

import (
	"image"
	"imageboard/config"
	"imageboard/models"
	"imageboard/utils/format"
	"imageboard/utils/validators"
	"strings"
)

func TransformImageToVariant(img image.Image, variant config.ImageSizeType) (models.ImageSize, image.Image, error) {
	variantSizeMap := map[config.ImageSizeType]int{
		config.ImageSizeTypeIcon:      64,
		config.ImageSizeTypeThumbnail: 256,
		config.ImageSizeTypeSmall:     512,
		config.ImageSizeTypeMedium:    1024,
		config.ImageSizeTypeLarge:     2048,
		config.ImageSizeTypeOriginal:  0, // Original size, no resizing
	}

	maxWidth := variantSizeMap[variant]
	if maxWidth > 0 {
		img = ResizeImage(img, maxWidth)
	}

	fileSize := format.GetImageFileSize(img)

	return models.ImageSize{
		SizeType: variant,
		Width:    img.Bounds().Dx(),
		Height:   img.Bounds().Dy(),
		FileSize: fileSize,
	}, img, nil
}

func ResizeImage(img image.Image, maxWidth int) image.Image {
	if maxWidth <= 0 || img.Bounds().Dx() <= maxWidth {
		return img
	}

	ratio := float64(maxWidth) / float64(img.Bounds().Dx())
	newWidth := int(float64(img.Bounds().Dx()) * ratio)
	newHeight := int(float64(img.Bounds().Dy()) * ratio)
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(float64(x) / ratio)
			srcY := int(float64(y) / ratio)
			if srcX < img.Bounds().Dx() && srcY < img.Bounds().Dy() {
				newImg.Set(x, y, img.At(srcX, srcY))
			}
		}
	}
	return newImg
}

func CreateUniqueFileName(sourceURLOrOriginalName, imageFormat string) string {
	fileName := sourceURLOrOriginalName
	if validators.IsValidURL(sourceURLOrOriginalName) {
		parts := strings.Split(sourceURLOrOriginalName, "/")
		fileName = parts[len(parts)-1]
	}

	currentTime := format.GetCurrentTimeAsTimestamp()
	fileNameWithoutExtension := format.RemoveExtension(fileName)
	fileName = GenerateTokenFromString(fileNameWithoutExtension + "_" + format.Int64ToString(currentTime))

	if len(fileName) > 32 {
		mid := len(fileName) / 2
		fileName = fileName[mid-16 : mid+16]
	}
	return fileName + "." + imageFormat
}

func ConvertStringRatingToType(rating string) (config.Rating, error) {
	switch rating {
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
