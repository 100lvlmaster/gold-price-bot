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

func StartTelegramBroadcastWorker(bot *tgbotapi.BotAPI, db *gorm.DB, intervalHours int) {
	ticker := time.NewTicker(time.Duration(intervalHours) * time.Hour)
	go func() {
		for {
			<-ticker.C
			log.Println("Broadcasting latest gold price to active users...")
			latest, err := database.GetLatestGoldPrice(db)
			if err != nil {
				log.Printf("Error fetching latest price for broadcast: %v", err)
				continue
			}

			var users []models.User
			if err := db.Where("is_active = ?", true).Find(&users).Error; err != nil {
				log.Printf("Error fetching active users for broadcast: %v", err)
				continue
			}

			ist := time.FixedZone("IST", 5.5*60*60)
			istTime := latest.Timestamp.In(ist)
			responseText := fmt.Sprintf("Scheduled Gold Price Update: ₹%.2f\nTime: %s", latest.Price, istTime.Format("2006-01-02 15:04:05"))

			for _, user := range users {
				msg := tgbotapi.NewMessage(user.ChatID, responseText)
				if _, err := bot.Send(msg); err != nil {
					log.Printf("Error sending broadcast to user %d: %v", user.ChatID, err)
				}
			}
		}
	}()
}
