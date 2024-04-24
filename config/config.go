package config

import (
	"os"
	"strconv"
)

type Config struct {
	APIPort            int
	JWTSecret          string
	DBURL              string
	S3APIKey           string
	S3APISecret        string
	S3APIEndpoint      string
	S3Region           string
	GoogleClientID     string
	GoogleClientSecret string
	BaseHostURL        string
	PostmarkAPIKey     string
	EmailFrom          string
}

func LoadConfig() (*Config, error) {
	port, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		port = 3000 // default port if not specified
	}

	return &Config{
		APIPort:            port,
		JWTSecret:          os.Getenv("JWT_SECRET"),
		DBURL:              os.Getenv("DATABASE_URL"),
		S3APIKey:           os.Getenv("S3_API_KEY"),
		S3APISecret:        os.Getenv("S3_API_SECRET"),
		S3APIEndpoint:      os.Getenv("S3_API_ENDPOINT"),
		S3Region:           os.Getenv("S3_REGION"),
		GoogleClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		BaseHostURL:        os.Getenv("BASE_HOST_URL"),
		PostmarkAPIKey:     os.Getenv("POSTMARK_API_KEY"),
		EmailFrom:          os.Getenv("EMAIL_FROM"),
	}, nil
}
