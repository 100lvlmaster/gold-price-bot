package worker

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gold-price-bot.com/v2/internal/database"
	"gold-price-bot.com/v2/internal/models"
	"gold-price-bot.com/v2/internal/scraper"
	"gorm.io/gorm"
)

// const timeformat = "2006-01-02 03:04 PM"

// TelegramBot defines the interface for sending messages to allow mocking in unit tests.
type TelegramBot interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

func StartGoldPriceWorker(bot TelegramBot, db *gorm.DB, intervalHours int, gstMultiplier float64) {
	go func() {
		// Run first execution immediately on startup
		log.Println("Executing initial startup scrape and broadcast...")
		runScrapeAndBroadcast(bot, db, gstMultiplier, true)

		ticker := time.NewTicker(time.Duration(intervalHours) * time.Hour)
		defer ticker.Stop()

		for {
			<-ticker.C
			log.Println("Executing scheduled scrape and broadcast check...")
			runScrapeAndBroadcast(bot, db, gstMultiplier, false)
		}
	}()
}

// Scrape gold price
// 1. Broadcast on first run
// 2. Skip if price hasn't changed since last time.
// 3. Broadcast and save to DB once price changes.
func runScrapeAndBroadcast(bot TelegramBot, db *gorm.DB, gstMultiplier float64, isStartup bool) {
	log.Println("Fetching gold price...")
	price, err := scraper.ScrapeGoldPrice(gstMultiplier)
	latest, latestError := database.GetLatestGoldPrice(db)
	if err != nil {
		log.Printf("Error scraping gold price: %v", err)
		// If scrape fails during startup, we still try to broadcast the last known price from DB
		if isStartup && latestError == nil {
			log.Println("Scrape failed on startup. Broadcasting last known gold price from database...")
			broadcastToAllUsers(bot, db, latest)
		}
		return
	}

	goldPrice := models.GoldPrice{
		Price:     price,
		Timestamp: time.Now().UTC(),
	}

	if isStartup {
		// Always broadcast on startup
		broadcastToAllUsers(bot, db, goldPrice)
	}

	if latest.Price == price {
		log.Printf("Price has not changed since last time. Not saving.")
		return
	}

	if err := db.Create(&goldPrice).Error; err != nil {
		log.Printf("Error saving gold price to database: %v", err)
		return
	}

	log.Printf("Successfully saved gold price: %.2f", price)
	broadcastToAllUsers(bot, db, goldPrice)
}

func broadcastToAllUsers(bot TelegramBot, db *gorm.DB, latest models.GoldPrice) {
	var users []models.User
	if err := db.Where("is_active = ?", true).Find(&users).Error; err != nil {
		log.Printf("Error fetching active users for broadcast: %v", err)
		return
	}

	// stats, err := database.GetGoldPriceStats(db)
	var responseText string
	// if err == nil {
	responseText = fmt.Sprintf(
		"Gold Price Update: ₹%.2f",
		latest.Price,
	)
	// }

	for _, user := range users {
		msg := tgbotapi.NewMessage(user.ChatID, responseText)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error sending broadcast to user %d: %v", user.ChatID, err)
		}
	}
}
