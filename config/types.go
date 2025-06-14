package config

import "time"

type ServerConfig struct {
	Host    string `env:"SERVER_HOST" default:"localhost"`
	Port    int    `env:"SERVER_PORT" default:"8080"`
	AppName string `env:"APP_NAME" default:"ImageBoard"`
}

type DatabaseConfig struct {
	Host         string `env:"DB_HOST" default:"localhost"`
	Port         int    `env:"DB_PORT" default:"5432"`
	Username     string `env:"DB_USERNAME" default:"postgres"`
	Password     string `env:"DB_PASSWORD" default:""`
	DatabaseName string `env:"DB_NAME" default:"imageboard"`
	SSLMode      string `env:"DB_SSLMODE" default:"disable"`
}

type SessionConfig struct {
	Expiration     time.Duration `env:"SESSION_EXPIRATION" default:"24h"`
	CookieName     string        `env:"SESSION_COOKIE_NAME" default:"session_id"`
	CookieDomain   string        `env:"SESSION_COOKIE_DOMAIN" default:""`
	CookiePath     string        `env:"SESSION_COOKIE_PATH" default:"/"`
	CookieSecure   bool          `env:"SESSION_COOKIE_SECURE" default:"false"`
	CookieSameSite string        `env:"SESSION_COOKIE_SAMESITE" default:"Lax"`
}

type ImageConfig struct {
	MaxSize      int    `env:"IMAGE_MAX_SIZE" default:"10485760"`
	AllowedTypes string `env:"IMAGE_ALLOWED_TYPES" default:"image/jpeg,image/png,image/gif,image/webp"`
}

type S3Config struct {
	Endpoint        string `env:"S3_ENDPOINT" default:"localhost:9000"`
	AccessKey       string `env:"S3_ACCESS_KEY" default:"minioadmin"`
	SecretAccessKey string `env:"S3_SECRET_KEY" default:"minioadmin"`
	BucketName      string `env:"S3_BUCKET_NAME" default:"shifoo"`
	FolderPath      string `env:"S3_FOLDER_PATH" default:"imageboard"`
	Region          string `env:"S3_REGION" default:"us-east-1"`
	UseSSL          bool   `env:"S3_USE_SSL" default:"false"`
	PublicURL       string `env:"S3_PUBLIC_URL" default:""`
}

type SMTPConfig struct {
	Host     string `env:"SMTP_HOST" default:""`
	Port     int    `env:"SMTP_PORT" default:"587"`
	Username string `env:"SMTP_USERNAME" default:""`
	Password string `env:"SMTP_PASSWORD" default:""`
	From     string `env:"EMAIL_FROM" default:""`
}
