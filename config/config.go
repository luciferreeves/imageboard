package config

import (
	"log"

	"github.com/joho/godotenv"
)

var (
	Database DatabaseConfig
	Upload   UploadConfig
	S3       S3Config
	Server   ServerConfig
	Session  SessionConfig
	SMTP     SMTPConfig
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	if err := Parse(&Database); err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}

	if err := Parse(&Upload); err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}

	if err := Parse(&S3); err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}

	if err := Parse(&Server); err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}

	if err := Parse(&Session); err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}

	if err := Parse(&SMTP); err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}
}
