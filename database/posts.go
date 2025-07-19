package database

import (
	"imageboard/config"
	"imageboard/models"
)

func GetPosts(limit int, ratings []config.Rating, tags []string) ([]models.Image, error) {
	var posts []models.Image
	query := DB.Preload("Sizes").Preload("Uploader").Limit(limit).Order("created_at DESC")

	if len(ratings) > 0 {
		query = query.Where("rating IN ?", ratings)
	}

	if len(tags) > 0 {
		query = query.Joins("JOIN image_tags ON images.id = image_tags.image_id").
			Joins("JOIN tags ON image_tags.tag_id = tags.id").
			Where("tags.name IN ?", tags).
			Group("images.id").
			Preload("Tags")
	} else {
		query = query.Preload("Tags")
	}

	if err := query.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func GetPostByID(postID uint) (*models.Image, error) {
	var post models.Image
	if err := DB.Preload("Sizes").Preload("Uploader").Preload("Approver").Preload("Tags").First(&post, postID).Error; err != nil {
		return nil, err
	}
	return &post, nil
}
