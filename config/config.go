package config

import (
	"os"
)

type Config struct {
	Port           string
	SendGridAPIKey string
	FromEmail      string
	DatabaseURL    string
}

func LoadConfig() *Config {
	return &Config{
		Port:           getEnv("PORT", "8081"),
		SendGridAPIKey: getEnv("SENDGRID_API_KEY", ""),
		FromEmail:      getEnv("FROM_EMAIL", "noreply@calendar.com"),
		DatabaseURL:    getEnv("DATABASE_URL", "calendar.db"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

