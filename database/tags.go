package database

import (
	"fmt"
	"imageboard/config"
	"imageboard/models"
	"strings"
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

func SearchTags(query string, limit int, offset int, tagType *config.TagType) ([]models.Tag, error) {
	var tags []models.Tag
	searchPattern := "%" + strings.TrimSpace(strings.ToLower(query)) + "%"

	dbQuery := DB.Where("name LIKE ? AND is_deleted = ?", searchPattern, false)
	if tagType != nil && strings.ToLower(string(*tagType)) != "" {
		dbQuery = dbQuery.Where("type = ?", strings.ToLower(string(*tagType)))
	}
	dbQuery = dbQuery.Order("count DESC, name ASC").Limit(limit).Offset(offset)

	err := dbQuery.Find(&tags).Error
	return tags, err
}

func SearchTagsExcluding(query string, imageID uint, limit int, tagType *config.TagType) ([]models.Tag, error) {
	var tags []models.Tag
	searchPattern := "%" + strings.TrimSpace(strings.ToLower(query)) + "%"

	dbQuery := DB.Where("name LIKE ? AND is_deleted = ? AND id NOT IN (?)",
		searchPattern, false,
		DB.Table("image_tags").Select("tag_id").Where("image_id = ?", imageID))

	if tagType != nil && strings.ToLower(string(*tagType)) != "" {
		dbQuery = dbQuery.Where("type = ?", strings.ToLower(string(*tagType)))
	}

	err := dbQuery.Order("count DESC, name ASC").Limit(limit).Find(&tags).Error

	return tags, err
}

func FindOrCreateTag(name string, tagType config.TagType) (*models.Tag, error) {
	name = strings.TrimSpace(strings.ToLower(name))

	// First check for active tag with exact name and type match
	var tag models.Tag
	if err := DB.Where("name = ? AND type = ? AND is_deleted = ?", name, tagType, false).First(&tag).Error; err == nil {
		return &tag, nil
	}

	// Check if a tag with the same name but different type exists
	var existingTag models.Tag
	if err := DB.Where("name = ? AND is_deleted = ?", name, false).First(&existingTag).Error; err == nil {
		if existingTag.Type != tagType {
			return nil, fmt.Errorf("tag '%s' already exists as %s type", name, existingTag.Type)
		}
	}

	// Check for deleted tag with same name and type and restore it
	if err := DB.Where("name = ? AND type = ? AND is_deleted = ?", name, tagType, true).First(&tag).Error; err == nil {
		tag.IsDeleted = false
		if err := DB.Save(&tag).Error; err != nil {
			return nil, fmt.Errorf("failed to restore tag: %v", err)
		}
		return &tag, nil
	}

	// Check if a deleted tag with same name but different type exists
	var deletedTag models.Tag
	if err := DB.Where("name = ? AND is_deleted = ?", name, true).First(&deletedTag).Error; err == nil {
		if deletedTag.Type != tagType {
			return nil, fmt.Errorf("tag '%s' previously existed as %s type", name, deletedTag.Type)
		}
	}

	// Create new tag
	tag = models.Tag{
		Name: name,
		Type: tagType,
	}

	if err := DB.Create(&tag).Error; err != nil {
		return nil, err
	}

	return &tag, nil
}

func AddTagToImage(imageID uint, tagID uint) error {
	// First get the tag to validate it exists and is not deleted
	var tag models.Tag
	if err := DB.Where("id = ? AND is_deleted = ?", tagID, false).First(&tag).Error; err != nil {
		return fmt.Errorf("tag not found or is deleted")
	}

	// Check if the association already exists
	var count int64
	err := DB.Table("image_tags").Where("image_id = ? AND tag_id = ?", imageID, tagID).Count(&count).Error
	if err != nil {
		return err
	}

	// If it doesn't exist, create it
	if count == 0 {
		err := DB.Exec("INSERT INTO image_tags (image_id, tag_id) VALUES (?, ?)", imageID, tagID).Error
		if err != nil {
			return err
		}
		// Increment tag count by 1
		return DB.Model(&models.Tag{}).Where("id = ?", tagID).Update("count", DB.Raw("count + 1")).Error
	}

	return nil // Already exists
}

func RemoveTagFromImage(imageID uint, tagID uint) error {
	err := DB.Exec("DELETE FROM image_tags WHERE image_id = ? AND tag_id = ?", imageID, tagID).Error
	if err != nil {
		return err
	}
	// Decrement tag count by 1
	return DB.Model(&models.Tag{}).Where("id = ?", tagID).Update("count", DB.Raw("count - 1")).Error
}

func GetImageTags(imageID uint) (map[string][]models.Tag, error) {
	var tags []models.Tag
	err := DB.Joins("JOIN image_tags ON image_tags.tag_id = tags.id").
		Where("image_tags.image_id = ? AND tags.is_deleted = ?", imageID, false).
		Preload("Parent").Preload("Children").Find(&tags).Error

	if err != nil {
		return nil, err
	}

	result := map[string][]models.Tag{
		"general":   {},
		"artist":    {},
		"character": {},
		"copyright": {},
		"meta":      {},
	}

	for _, tag := range tags {
		switch tag.Type {
		case config.TagTypeGeneral:
			result["general"] = append(result["general"], tag)
		case config.TagTypeArtist:
			result["artist"] = append(result["artist"], tag)
		case config.TagTypeCharacter:
			result["character"] = append(result["character"], tag)
		case config.TagTypeCopyright:
			result["copyright"] = append(result["copyright"], tag)
		case config.TagTypeMeta:
			result["meta"] = append(result["meta"], tag)
		}
	}

	return result, nil
}

func GetTagWithAncestors(tagID uint) (*models.Tag, []models.Tag, error) {
	var tag models.Tag
	if err := DB.Preload("Parent").Preload("Children").First(&tag, tagID).Error; err != nil {
		return nil, nil, err
	}

	var ancestors []models.Tag
	current := &tag
	for current.Parent != nil {
		ancestors = append(ancestors, *current.Parent)
		current = current.Parent
	}

	return &tag, ancestors, nil
}

func GetTagWithDescendants(tagID uint) (*models.Tag, []models.Tag, error) {
	var tag models.Tag
	if err := DB.Preload("Children").First(&tag, tagID).Error; err != nil {
		return nil, nil, err
	}

	var descendants []models.Tag
	var getChildren func(t *models.Tag)
	getChildren = func(t *models.Tag) {
		for _, child := range t.Children {
			descendants = append(descendants, child)
			childWithChildren := models.Tag{}
			DB.Preload("Children").First(&childWithChildren, child.ID)
			getChildren(&childWithChildren)
		}
	}

	getChildren(&tag)
	return &tag, descendants, nil
}
