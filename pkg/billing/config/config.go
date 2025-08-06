package config

import (
	"fmt"
	"log"
	"os"

	"github.com/doddeeph/billing-engine/internal/billing/model"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Load() {
	currDir, _ := os.Getwd()
	err := godotenv.Load(currDir + "/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func InitBillingDB() *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		dbHost, dbPort, dbName, dbUser, dbPass,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to DB")
	}
	db.AutoMigrate(&model.Billing{}, &model.Payment{})
	return db
}

func GetAppPort() string {
	return os.Getenv("APP_PORT")
}
