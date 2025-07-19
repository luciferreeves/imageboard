package database

import (
	"imageboard/config"
	"imageboard/models"
	"imageboard/utils/format"
	"imageboard/utils/transformers"
	"time"

	"gorm.io/gorm"
)

func GetTotalPostsCount() (int64, error) {
	var count int64
	err := DB.Model(&models.Image{}).Count(&count).Error
	return count, err
}

func GetTodayPostsCount() (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := DB.Model(&models.Image{}).Where("created_at >= ?", today).Count(&count).Error
	return count, err
}

func GetTotalStorageSize() (string, error) {
	var imageSizes []models.ImageSize
	if err := DB.Select("file_size").Find(&imageSizes).Error; err != nil {
		return "0 B", err
	}

	var totalSize int64
	for _, size := range imageSizes {
		totalSize += size.FileSize
	}

	return format.FileSize(totalSize), nil
}

func CreateImageWithTx(tx *gorm.DB, fileName, contentType, md5Hash, sourceURL, rating string, uploaderID uint, requiresApproval bool) (*models.Image, error) {
	ratingEnum, err := transformers.ConvertStringRatingToType(rating)
	if err != nil {
		return nil, err
	}

	contentTypeEnum, err := transformers.ConvertStringToContentType(contentType)
	if err != nil {
		return nil, err
	}

	image := models.Image{
		FileName:    fileName,
		ContentType: contentTypeEnum,
		MD5Hash:     md5Hash,
		SourceURL:   sourceURL,
		Rating:      ratingEnum,
		UploaderID:  uploaderID,
		IsApproved:  !requiresApproval,
	}

	if err := tx.Create(&image).Error; err != nil {
		return nil, err
	}

	return &image, nil
}

func CreateImageSizeWithTx(tx *gorm.DB, imageID uint, sizeType config.ImageSizeType, width, height int, fileSize int64) (*models.ImageSize, error) {
	imageSize := models.ImageSize{
		ImageID:  imageID,
		SizeType: sizeType,
		Width:    width,
		Height:   height,
		FileSize: fileSize,
	}

	if err := tx.Create(&imageSize).Error; err != nil {
		return nil, err
	}

	return &imageSize, nil
}

func UpdateImage(imageID uint, updates map[string]interface{}) error {
	return DB.Model(&models.Image{}).Where("id = ?", imageID).Updates(updates).Error
}
