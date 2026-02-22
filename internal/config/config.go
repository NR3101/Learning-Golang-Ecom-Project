package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Aws      AwsConfig
	Upload   UploadConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret                string
	ExpiresIn             time.Duration
	RefreshTokenExpiresIn time.Duration
}

type AwsConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	S3Bucket        string
	S3Endpoint      string
	EventQueueName  string
}

type UploadConfig struct {
	Path           string
	MaxFileSize    int64
	UploadProvider string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	jwtExpiresIn, _ := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))                            // 24 hours
	jwtRefreshTokenExpiresIn, _ := time.ParseDuration(getEnv("JWT_REFRESH_TOKEN_EXPIRES_IN", "168h")) // 7 days
	maxUploadSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)               // 10 MB

	config := &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "release"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "ecomdb"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:                getEnv("JWT_SECRET", "mysecretkey"),
			ExpiresIn:             jwtExpiresIn,             // 24 hours
			RefreshTokenExpiresIn: jwtRefreshTokenExpiresIn, // 7 days
		},
		Aws: AwsConfig{
			Region:          getEnv("AWS_REGION", "us-east-1"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", "test"),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", "test"),
			S3Bucket:        getEnv("AWS_S3_BUCKET_NAME", "ecom-uploads"),
			S3Endpoint:      getEnv("AWS_S3_ENDPOINT", "http://localhost:9000"), // For local testing
			EventQueueName:  getEnv("AWS_EVENT_QUEUE_NAME", "ecom-events"),
		},
		Upload: UploadConfig{
			Path:           getEnv("UPLOAD_PATH", "./uploads"),
			MaxFileSize:    maxUploadSize,                      // 10 MB
			UploadProvider: getEnv("UPLOAD_PROVIDER", "local"), // "local" or "s3"
		},
	}
	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
