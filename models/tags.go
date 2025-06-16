package models

import (
	"fmt"
	"imageboard/utils/validators"
	"strings"

	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name        string  `gorm:"not null;uniqueIndex;size:100" json:"name"`
	Type        TagType `gorm:"not null;default:'general';size:20" json:"type"`
	Description string  `gorm:"default:'';type:text" json:"description"`
	Count       int     `gorm:"not null;default:0" json:"count"`
	IsDeleted   bool    `gorm:"not null;default:false" json:"is_deleted"`
	ParentID    *uint   `gorm:"index" json:"-"`
	Parent      *Tag    `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Tag   `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Images      []Image `gorm:"many2many:image_tags" json:"images,omitempty"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	t.Name = strings.TrimSpace(strings.ToLower(t.Name))
	t.Description = strings.TrimSpace(t.Description)

	if t.Name == "" {
		return fmt.Errorf("tag name cannot be empty")
	}

	if len(t.Name) < 2 || len(t.Name) > 100 {
		return fmt.Errorf("tag name must be between 2 and 100 characters")
	}

	if !validators.IsValidTagName(t.Name) {
		return fmt.Errorf("tag name can only contain letters, numbers, and underscores")
	}

	var existingTag Tag
	if err := tx.Where("name = ?", t.Name).First(&existingTag).Error; err == nil {
		return fmt.Errorf("tag name '%s' is already taken", t.Name)
	}

	return nil
}

func (t *Tag) BeforeUpdate(tx *gorm.DB) error {
	t.Name = strings.TrimSpace(strings.ToLower(t.Name))
	t.Description = strings.TrimSpace(t.Description)
	return nil
}

func (t *Tag) GetFullPath() string {
	if t.Parent == nil {
		return t.Name
	}
	return t.Parent.GetFullPath() + ":" + t.Name
}

func SearchTags(tx *gorm.DB, query string, limit int) ([]Tag, error) {
	var tags []Tag
	searchPattern := "%" + strings.TrimSpace(strings.ToLower(query)) + "%"

	err := tx.Where("name LIKE ? AND is_deleted = ?", searchPattern, false).
		Order("count DESC, name ASC").Limit(limit).Find(&tags).Error

	return tags, err
}

func SearchTagsExcluding(tx *gorm.DB, query string, imageID uint, limit int) ([]Tag, error) {
	var tags []Tag
	searchPattern := "%" + strings.TrimSpace(strings.ToLower(query)) + "%"

	err := tx.Where("name LIKE ? AND is_deleted = ? AND id NOT IN (?)",
		searchPattern, false,
		tx.Table("image_tags").Select("tag_id").Where("image_id = ?", imageID)).
		Order("count DESC, name ASC").Limit(limit).Find(&tags).Error

	return tags, err
}

func FindOrCreateTag(tx *gorm.DB, name string, tagType TagType) (*Tag, error) {
	name = strings.TrimSpace(strings.ToLower(name))

	// First check for active tag
	var tag Tag
	if err := tx.Where("name = ? AND is_deleted = ?", name, false).First(&tag).Error; err == nil {
		return &tag, nil
	}

	// Check for deleted tag and restore it
	if err := tx.Where("name = ? AND is_deleted = ?", name, true).First(&tag).Error; err == nil {
		tag.IsDeleted = false
		tag.Type = tagType // Update type in case it changed
		if err := tx.Save(&tag).Error; err != nil {
			return nil, fmt.Errorf("failed to restore tag: %v", err)
		}
		return &tag, nil
	}

	// Create new tag
	tag = Tag{
		Name: name,
		Type: tagType,
	}

	if err := tx.Create(&tag).Error; err != nil {
		return nil, err
	}

	return &tag, nil
}

func (t *Tag) DeleteTag(tx *gorm.DB) error {
	if t.IsDeleted {
		return fmt.Errorf("tag is already deleted")
	}

	if err := tx.Model(t).Association("Images").Clear(); err != nil {
		return fmt.Errorf("failed to clear image associations: %v", err)
	}

	t.IsDeleted = true
	t.Count = 0
	return tx.Save(t).Error
}
