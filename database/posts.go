package database

import "imageboard/models"

func GetPosts(limit int) ([]models.Image, error) {
	var posts []models.Image
	if err := DB.Preload("Sizes").Preload("Uploader").Preload("Tags").Limit(limit).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
