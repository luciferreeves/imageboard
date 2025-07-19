package database

import (
	"imageboard/config"
	"imageboard/models"
)

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func ListAllUsers() ([]models.User, error) {
	var users []models.User
	err := DB.Where("is_deleted = ?", false).Order("LOWER(username) ASC").Find(&users).Error
	return users, err
}

func ListAllApprovers() ([]models.User, error) {
	var users []models.User
	err := DB.Where("is_deleted = ? AND level >= ?", false, config.UserLevelJanitor).Order("LOWER(username) ASC").Find(&users).Error
	return users, err
}

func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *models.User) error {
	return DB.Create(user).Error
}
