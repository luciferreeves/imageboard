package config

import (
	"log"

	"github.com/joho/godotenv"
)

var (
	AppSecret         string = "default_secret"
	Database          DatabaseConfig
	IsDevelopmentMode bool = true
	Image             ImageConfig
	S3                S3Config
	Server            ServerConfig
	Session           SessionConfig
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	AppSecret = getEnv("APP_SECRET", AppSecret)

	if err := Parse(&Database); err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}

	IsDevelopmentMode = getEnvBool("DEV_MODE", IsDevelopmentMode)

	if err := Parse(&Image); err != nil {
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
}
