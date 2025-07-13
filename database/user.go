package database

import "imageboard/models"

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
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
