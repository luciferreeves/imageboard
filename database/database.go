package database

import (
	"fmt"
	"imageboard/config"
	"imageboard/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	err error
)

func init() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.Username,
		config.Database.DatabaseName,
		config.Database.SSLMode,
	)

	if config.Database.Password != "" {
		dsn += fmt.Sprintf(" password=%s", config.Database.Password)
	}

	logLevel := logger.Silent
	if config.Server.IsDevMode {
		logLevel = logger.Info
	}

	dialector := postgres.Open(dsn)

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if config.Server.IsDevMode && config.Database.WipeAndResetDatabase {
		if err := wipeAndResetDatabase(); err != nil {
			log.Fatalf("failed to wipe and reset database: %v", err)
		}
		log.Println("Database wiped and reset successfully")
	}

	if err := autoMigrate(); err != nil {
		log.Fatalf("failed to auto migrate database: %v", err)
	}

	log.Println("Database connection established successfully")
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Image{},
		&models.ImageSize{},
		&models.Tag{},
		&models.Comment{},
	)
}

func wipeAndResetDatabase() error {
	if err := DB.Exec("DROP SCHEMA public CASCADE").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE SCHEMA public").Error; err != nil {
		return err
	}
	return nil
}
