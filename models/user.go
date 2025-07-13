package models

import (
	"fmt"
	"imageboard/config"
	"imageboard/utils/validators"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username             string     `gorm:"uniqueIndex;not null;size:255" json:"username"`
	Email                string     `gorm:"not null;size:255" json:"email"`
	Password             string     `gorm:"not null;size:255" json:"-"`
	Level                UserLevel  `gorm:"not null;default:0" json:"level"`
	EmailVerified        bool       `gorm:"not null;default:false" json:"email_verified"`
	Bio                  string     `gorm:"default:'';size:500" json:"bio"`
	AvatarURL            string     `gorm:"default:'';size:255" json:"avatar_url"`
	WebsiteURL           string     `gorm:"default:'';size:255" json:"website_url"`
	Location             string     `gorm:"default:'';size:255" json:"location"`
	Timezone             string     `gorm:"default:'UTC';size:50" json:"timezone"`
	AccountDisabled      bool       `gorm:"not null;default:false" json:"-"`
	AccountBanned        bool       `gorm:"not null;default:false" json:"-"`
	PostsRequireApproval bool       `gorm:"not null;default:false" json:"-"`
	IsDeleted            bool       `gorm:"not null;default:false" json:"-"`
	LastLoginAt          *time.Time `gorm:"default:null" json:"last_login_at"`
	LastActivityAt       *time.Time `gorm:"default:null" json:"last_activity_at"`
	Images               []Image    `gorm:"foreignKey:UploaderID" json:"images,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))

	if u.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if u.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if len(u.Username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}

	if len(u.Username) > 72 {
		return fmt.Errorf("username must not exceed 72 characters")
	}

	if !validators.IsValidUsername(u.Username) {
		return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
	}

	if validators.IsReservedUsername(u.Username) {
		return fmt.Errorf("username '%s' is reserved and cannot be used", u.Username)
	}

	if !validators.IsValidEmail(u.Email) {
		return fmt.Errorf("invalid email format")
	}

	// Check if username is already taken
	var existingUser User
	if err := tx.Where("username = ?", u.Username).First(&existingUser).Error; err == nil {
		return fmt.Errorf("username '%s' is already taken", u.Username)
	}

	var userCount int64
	if err := tx.Model(&User{}).Where("is_deleted = ?", false).Count(&userCount).Error; err != nil {
		return fmt.Errorf("failed to count existing users: %v", err)
	}

	if userCount == 0 {
		u.Level = UserLevelSuperAdmin // First user becomes Super Admin
	}

	if len(u.Password) < config.Server.MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", config.Server.MinPasswordLength)
	}
	if len(u.Password) > 255 {
		return fmt.Errorf("password must not exceed 255 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))

	return nil
}

func (u *User) SetPassword(password string) error {
	if len(password) < config.Server.MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", config.Server.MinPasswordLength)
	}
	if len(password) > 255 {
		return fmt.Errorf("password must not exceed 255 characters")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	if u.IsDeleted {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) IsActive() bool {
	return !u.IsDeleted && !u.AccountDisabled && !u.AccountBanned
}

func (u *User) CanLogin() bool {
	return u.IsActive() && (u.EmailVerified || u.IsAdmin())
}

func (u *User) IsAdmin() bool {
	return u.Level >= UserLevelAdmin
}

func (u *User) IsModerator() bool {
	return u.IsActive() && u.Level >= UserLevelModerator
}

func (u *User) IsJanitor() bool {
	return u.IsActive() && u.Level >= UserLevelJanitor
}

func (u *User) IsContributor() bool {
	return u.IsActive() && u.Level >= UserLevelContributor
}

func (u *User) IsMember() bool {
	return u.IsActive() && u.Level >= UserLevelMember
}

func (u *User) CanUpload() bool {
	return u.IsActive() && (u.EmailVerified || u.IsAdmin())
}

func (u *User) CanComment() bool {
	return u.IsActive() && (u.EmailVerified || u.IsAdmin())
}

func (u *User) CanMessage() bool {
	return u.IsActive() && (u.EmailVerified || u.IsAdmin())
}

func (u *User) CanCreateTags() bool {
	return u.IsContributor()
}

func (u *User) CanEditTags() bool {
	return u.IsJanitor()
}

func (u *User) CanEditPosts() bool {
	return u.IsJanitor()
}

func (u *User) CanDeletePosts() bool {
	return u.IsModerator()
}

func (u *User) CanApprovePosts() bool {
	return u.IsJanitor()
}

func (u *User) CanEditUser(targetUser *User) bool {
	if u.ID == targetUser.ID {
		return true
	}

	if targetUser.IsDeleted {
		return false
	}

	return (u.IsAdmin() || u.IsModerator()) && targetUser.Level < u.Level
}

func (u *User) CanPromoteUser(targetUser *User, newLevel UserLevel) bool {
	if u.ID == targetUser.ID || targetUser.IsDeleted {
		return false
	}

	if u.Level <= UserLevelContributor {
		return false
	}

	return newLevel > UserLevelMember && newLevel <= u.Level && newLevel <= UserLevelAdmin
}

func (u *User) CanDemoteUser(targetUser *User, newLevel UserLevel) bool {
	if u.ID == targetUser.ID || targetUser.IsDeleted {
		return false
	}

	if u.Level <= UserLevelContributor {
		return false
	}

	return newLevel >= UserLevelMember && newLevel < u.Level && newLevel <= UserLevelAdmin
}

func (u *User) CanDisableUser(targetUser *User) bool {
	if u.ID == targetUser.ID || targetUser.IsDeleted {
		return false
	}

	return (u.IsJanitor() || u.IsModerator() || u.IsAdmin()) && targetUser.Level < u.Level
}

func (u *User) CanBanUser(targetUser *User) bool {
	if u.ID == targetUser.ID || targetUser.IsDeleted {
		return false
	}

	return (u.IsModerator() || u.IsAdmin()) && targetUser.Level < u.Level
}

func (u *User) CanDeleteUser(targetUser *User) bool {
	if targetUser.IsDeleted {
		return false
	}

	if u.ID == targetUser.ID {
		return true // Users can delete their own account
	}

	if u.Level <= UserLevelContributor {
		return false
	}

	return (u.IsAdmin() || u.IsModerator()) && targetUser.Level < u.Level
}

func (u *User) CanMakeUserPostsRequireApproval(targetUser *User) bool {
	if targetUser.IsDeleted {
		return false
	}

	return (u.IsJanitor() || u.IsModerator() || u.IsAdmin()) && targetUser.Level < u.Level
}

func (u *User) GetDailyPostLimit() int {
	switch u.Level {
	case UserLevelMember:
		return 10
	case UserLevelContributor:
		return 25
	default:
		return -1 // No limit for Janitors, Moderators, and Admins
	}
}

func (u *User) GetDailyRemainingUploadLimit(tx *gorm.DB) (int64, error) {
	totalAllowed := u.GetDailyPostLimit()
	if totalAllowed == -1 {
		return -1, nil
	}

	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	var count int64
	err := tx.Model(&Image{}).Where("uploader_id = ? AND created_at >= ? AND created_at < ? AND is_deleted = ?",
		u.ID, today, tomorrow, false).Count(&count).Error

	return int64(totalAllowed) - count, err
}

func (u *User) CanUploadToday(tx *gorm.DB) (bool, error) {
	remaining, err := u.GetDailyRemainingUploadLimit(tx)
	if err != nil {
		return false, err
	}
	return remaining != 0, nil
}

func (u *User) UpdateLastUserActivity(tx *gorm.DB) error {
	now := time.Now()
	u.LastActivityAt = &now
	return tx.Model(u).Update("last_activity_at", now).Error
}

func (u *User) UpdateLastUserLogin(tx *gorm.DB) error {
	now := time.Now()
	u.LastLoginAt = &now
	u.LastActivityAt = &now
	return tx.Model(u).Updates(map[string]interface{}{
		"last_login_at":    now,
		"last_activity_at": now,
	}).Error
}

func (u *User) DeleteUser(tx *gorm.DB) error {
	if u.IsDeleted {
		return fmt.Errorf("user is already deleted")
	}

	u.IsDeleted = true
	return tx.Save(u).Error
}
