package database

import (
	"imageboard/models"
)

func GetTotalCommentsCount() (int64, error) {
	var count int64
	err := DB.Model(&models.Comment{}).Count(&count).Error
	return count, err
}
