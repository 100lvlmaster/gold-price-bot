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

func GetLatestGoldPriceIfChanged(db *gorm.DB) (models.GoldPrice, bool, error) {
	var recentPrices []models.GoldPrice
	result := db.Order("timestamp desc").Limit(2).Find(&recentPrices)
	if result.Error != nil {
		return models.GoldPrice{}, false, result.Error
	}

	if len(recentPrices) == 0 {
		return models.GoldPrice{}, false, nil
	}

	if len(recentPrices) == 1 {
		return recentPrices[0], true, nil
	}

	if recentPrices[0].Price != recentPrices[1].Price {
		return recentPrices[0], true, nil
	}

	return models.GoldPrice{}, false, nil
}
