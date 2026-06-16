package database

import (
	"log"

	"gold-price-bot.com/v2/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run Migrations
	if err := db.AutoMigrate(&models.User{}, &models.GoldPrice{}); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	return db
}

func GetLatestGoldPrice(db *gorm.DB) (models.GoldPrice, error) {
	var latestPrice models.GoldPrice
	result := db.Order("timestamp desc").First(&latestPrice)
	return latestPrice, result.Error
}
