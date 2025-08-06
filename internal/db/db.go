package db

import (
	"fmt"
	"log"

	"github.com/doddeeph/billing-engine/internal/config"
	"github.com/doddeeph/billing-engine/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.DBConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Name, cfg.User, cfg.Password,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open to DB: %v", err)
	}
	log.Println("Connected to database.")
	db.AutoMigrate(&model.Billing{}, &model.Payment{})
	return db
}
