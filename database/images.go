package database

import (
	"imageboard/models"
	"imageboard/utils/format"
	"time"
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
