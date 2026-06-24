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

// TelegramBot defines the interface for sending messages to allow mocking in unit tests.
type TelegramBot interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

func StartGoldPriceWorker(db *gorm.DB, intervalHours int, gstMultiplier float64) {
	ticker := time.NewTicker(time.Duration(intervalHours) * time.Hour)
	go func() {
		for {
			log.Println("Fetching gold price...")
			price, err := scraper.ScrapeGoldPrice(gstMultiplier)
			if err != nil {
				log.Printf("Error scraping gold price: %v", err)
			} else {
				goldPrice := models.GoldPrice{
					Price:     price,
					Timestamp: time.Now().UTC(),
				}
				if err := db.Create(&goldPrice).Error; err != nil {
					log.Printf("Error saving gold price to database: %v", err)
				} else {
					log.Printf("Successfully saved gold price: %.2f", price)
				}
			}
			<-ticker.C
		}
	}()
}

func StartTelegramBroadcastWorker(bot TelegramBot, db *gorm.DB, intervalHours int) {
	go func() {
		// Wait a few seconds for the first scrape to complete
		time.Sleep(10 * time.Second)

		log.Println("Performing initial startup broadcast...")
		if latest, err := database.GetLatestGoldPrice(db); err == nil {
			broadcastToAllUsers(bot, db, latest)
		} else {
			log.Printf("Error fetching latest price for startup broadcast: %v", err)
		}

		ticker := time.NewTicker(time.Duration(intervalHours) * time.Hour)
		defer ticker.Stop()

		for {
			<-ticker.C
			log.Println("Checking for price changes for scheduled broadcast...")
			latest, priceChanged, err := database.GetLatestGoldPriceIfChanged(db)
			if err != nil {
				log.Printf("Error fetching latest price for broadcast: %v", err)
				continue
			}

			if priceChanged {
				broadcastToAllUsers(bot, db, latest)
			} else {
				log.Println("Price hasn't changed, skipping scheduled broadcast.")
			}
		}
	}()
}

func broadcastToAllUsers(bot TelegramBot, db *gorm.DB, latest models.GoldPrice) {
	var users []models.User
	if err := db.Where("is_active = ?", true).Find(&users).Error; err != nil {
		log.Printf("Error fetching active users for broadcast: %v", err)
		return
	}

	ist := time.FixedZone("IST", 5.5*60*60)
	istTime := latest.Timestamp.In(ist)
	responseText := fmt.Sprintf("Gold Price Update: ₹%.2f\nTime: %s", latest.Price, istTime.Format("2006-01-02 15:04:05"))

	for _, user := range users {
		msg := tgbotapi.NewMessage(user.ChatID, responseText)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error sending broadcast to user %d: %v", user.ChatID, err)
		}
	}
}
