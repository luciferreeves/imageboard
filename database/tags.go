package database

import (
	"imageboard/models"
)

func GetTotalTagsCount() (int64, error) {
	var count int64
	err := DB.Model(&models.Tag{}).Where("is_deleted = ?", false).Count(&count).Error
	return count, err
}

func GetPopularTags(limit int) ([]models.Tag, error) {
	var tags []models.Tag
	err := DB.Where("is_deleted = ?", false).Order("count DESC").Limit(limit).Find(&tags).Error
	return tags, err
}

func GetRecentTags(limit int) ([]models.Tag, error) {
	var tags []models.Tag
	err := DB.Where("is_deleted = ?", false).Order("created_at DESC").Limit(limit).Find(&tags).Error
	return tags, err
}
