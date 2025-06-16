package models

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Body      string    `gorm:"not null;type:text" json:"body"`
	UserID    uint      `gorm:"not null;index" json:"-"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	ImageID   uint      `gorm:"not null;index" json:"-"`
	Image     Image     `gorm:"foreignKey:ImageID" json:"image"`
	ParentID  *uint     `gorm:"index" json:"-"`
	Parent    *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies   []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
	IsDeleted bool      `gorm:"not null;default:false" json:"is_deleted"`
	IsSticky  bool      `gorm:"not null;default:false" json:"is_sticky"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	c.Body = strings.TrimSpace(c.Body)
	if c.Body == "" {
		return fmt.Errorf("comment body cannot be empty")
	}

	if len(c.Body) > 10000 {
		return fmt.Errorf("comment body must not exceed 10,000 characters")
	}

	return nil
}

func (c *Comment) BeforeUpdate(tx *gorm.DB) error {
	c.Body = strings.TrimSpace(c.Body)

	if c.Body == "" {
		return fmt.Errorf("comment body cannot be empty")
	}

	if len(c.Body) > 10000 {
		return fmt.Errorf("comment body cannot exceed 10000 characters")
	}

	return nil
}

func (c *Comment) AfterCreate(tx *gorm.DB) error {
	return tx.Model(&Image{}).Where("id = ?", c.ImageID).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error
}

func (c *Comment) CanEdit(user *User) bool {
	if user == nil || !user.IsActive() {
		return false
	}

	if c.UserID == user.ID {
		return true
	}

	return user.CanEditPosts()
}

func (c *Comment) CanDelete(user *User) bool {
	if user == nil || !user.IsActive() {
		return false
	}

	if c.UserID == user.ID {
		return true
	}

	return user.CanDeletePosts()
}

func (c *Comment) DeleteComment(tx *gorm.DB) error {
	if c.IsDeleted {
		return fmt.Errorf("comment is already deleted")
	}

	c.IsDeleted = true
	if err := tx.Save(c).Error; err != nil {
		return err
	}

	return tx.Model(&Image{}).Where("id = ?", c.ImageID).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error
}

func (c *Comment) GetReplies(tx *gorm.DB) ([]Comment, error) {
	var replies []Comment
	err := tx.Where("parent_id = ? AND is_deleted = ?", c.ID, false).
		Preload("User").
		Order("created_at ASC").
		Find(&replies).Error
	return replies, err
}
