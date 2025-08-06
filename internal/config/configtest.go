package config

import (
	"log"

	"github.com/joho/godotenv"
)

type DBTestConfig struct {
	Image    string
	Host     string
	Post     string
	Port     string
	Name     string
	User     string
	Password string
}

type AppTestConfig struct {
	DB      DBTestConfig
	AppPort string
}

func LoadTestConfig() *AppTestConfig {
	err := godotenv.Load("../.env.test")
	if err != nil {
		log.Fatalf("Error loading .env.test file: %v", err)
	}
	return &AppTestConfig{
		DB: DBTestConfig{
			Image:    getEnv("DB_IMAGE", "postgres:15"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "billing_db"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
		},
		AppPort: getEnv("APP_PORT", "8080"),
	}
}
