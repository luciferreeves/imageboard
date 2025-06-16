package database

import (
	"fmt"
	"imageboard/config"
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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Database.Host,
		config.Database.Username,
		config.Database.Password,
		config.Database.DatabaseName,
		config.Database.Port,
		config.Database.SSLMode,
	)

	logLevel := logger.Silent
	if config.Server.IsDevMode {
		logLevel = logger.Info
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := autoMigrate(); err != nil {
		log.Fatalf("failed to auto migrate database: %v", err)
	}

	log.Println("Database connection established successfully")
}

func autoMigrate() error {
	return DB.AutoMigrate()
}
