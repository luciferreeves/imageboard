package models

import (
	"fmt"
	"imageboard/config"
	"imageboard/utils/validators"
	"strings"

	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name        string         `gorm:"not null;uniqueIndex;size:100" json:"name"`
	Type        config.TagType `gorm:"not null;default:'general';size:20" json:"type"`
	Description string         `gorm:"default:'';type:text" json:"description"`
	Count       int            `gorm:"not null;default:0" json:"count"`
	IsDeleted   bool           `gorm:"not null;default:false" json:"is_deleted"`
	ParentID    *uint          `gorm:"index" json:"-"`
	Parent      *Tag           `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Tag          `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Images      []Image        `gorm:"many2many:image_tags;joinForeignKey:tag_id;joinReferences:image_id;constraint:OnDelete:CASCADE" json:"images,omitempty"`
	Wiki        *TagWiki       `gorm:"foreignKey:TagID" json:"wiki,omitempty"`
}

type TagWiki struct {
	gorm.Model
	TagID       uint   `gorm:"not null;uniqueIndex" json:"-"`
	Tag         Tag    `gorm:"foreignKey:TagID" json:"tag,omitempty"`
	Content     string `gorm:"type:text" json:"content"`
	EditorID    uint   `gorm:"not null" json:"-"`
	Editor      User   `gorm:"foreignKey:EditorID" json:"editor,omitempty"`
	IsProtected bool   `gorm:"not null;default:false" json:"is_protected"`
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
