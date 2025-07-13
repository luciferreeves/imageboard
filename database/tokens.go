package database

import (
	"fmt"
	"imageboard/models"
	"imageboard/utils/validators"
	"time"
)

func GenerateEmailToken(userID int, tokenType models.EmailTokenType) (*models.EmailToken, error) {
	var existingToken models.EmailToken
	if err := DB.Where("user_id = ? AND type = ?", userID, tokenType).First(&existingToken).Error; err == nil {
		if err := DB.Delete(&existingToken).Error; err != nil {
			return nil, err
		}
	}

	tokenValue, err := validators.GenerateRandomToken()
	if err != nil {
		return nil, err
	}

	var expirationDuration time.Duration
	switch tokenType {
	case models.EmailTokenTypeVerification:
		expirationDuration = 24 * time.Hour
	case models.EmailTokenTypePasswordReset:
		expirationDuration = 1 * time.Hour
	case models.EmailTokenTypeChangeEmail:
		expirationDuration = 1 * time.Hour
	default:
		expirationDuration = 1 * time.Hour
	}

	token := &models.EmailToken{
		UserID:    uint(userID),
		Token:     tokenValue,
		Type:      tokenType,
		ExpiresAt: time.Now().Add(expirationDuration),
	}

	if err := DB.Create(token).Error; err != nil {
		return nil, err
	}

	return token, nil
}

func VerifyToken(userID int, token string, tokenType models.EmailTokenType) (*models.EmailToken, error) {
	var emailToken models.EmailToken
	if err := DB.Where("user_id = ? AND token = ? AND type = ?", userID, token, tokenType).First(&emailToken).Error; err != nil {
		return nil, err
	}

	if !emailToken.IsValid() {
		return nil, fmt.Errorf("token is invalid or expired")
	}

	emailToken.MarkAsUsed()
	if err := DB.Save(&emailToken).Error; err != nil {
		return nil, fmt.Errorf("failed to mark token as used: %w", err)
	}

	return &emailToken, nil
}
