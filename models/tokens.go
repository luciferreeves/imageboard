package models

import (
	"imageboard/config"
	"time"

	"gorm.io/gorm"
)

type EmailToken struct {
	gorm.Model
	UserID    uint                  `gorm:"not null;index" json:"user_id"`
	Token     string                `gorm:"uniqueIndex;not null;size:64" json:"token"`
	Type      config.EmailTokenType `gorm:"not null;size:20" json:"type"`
	ExpiresAt time.Time             `gorm:"not null" json:"expires_at"`
	UsedAt    *time.Time            `gorm:"default:null" json:"used_at"`
	User      User                  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (et *EmailToken) IsExpired() bool {
	return time.Now().After(et.ExpiresAt)
}

func (et *EmailToken) IsUsed() bool {
	return et.UsedAt != nil
}

func (et *EmailToken) IsValid() bool {
	return !et.IsExpired() && !et.IsUsed()
}

func (et *EmailToken) MarkAsUsed() {
	now := time.Now()
	et.UsedAt = &now
}
