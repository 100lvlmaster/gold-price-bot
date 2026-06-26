package bot

import (
	"errors"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gold-price-bot.com/v2/internal/database"
	"gold-price-bot.com/v2/internal/models"
	"gorm.io/gorm"
)

func StartTelegramBot(bot *tgbotapi.BotAPI, db *gorm.DB) {
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			chatID := update.Message.Chat.ID
			if update.Message.IsCommand() {
				command := update.Message.Command()
				handleCommand(bot, db, chatID, command, update.Message.From.UserName)
			}
		}
	}()
}

func handleCommand(bot *tgbotapi.BotAPI, db *gorm.DB, chatID int64, command string, userName string) {
	switch command {
	case "start":
		user := models.User{
			ChatID:   chatID,
			Name:     userName,
			IsActive: true,
		}
		if err := db.Where(models.User{ChatID: chatID}).Assign(models.User{IsActive: true, Name: user.Name}).FirstOrCreate(&user).Error; err != nil {
			log.Printf("Error during /start upsert: %v", err)
			bot.Send(tgbotapi.NewMessage(chatID, "Error starting the bot. Please try again later."))
		} else {
			bot.Send(tgbotapi.NewMessage(chatID, "Welcome! You will now receive gold price updates every 2 hours."))
		}

	case "stop":
		if err := db.Model(&models.User{}).Where("chat_id = ?", chatID).Update("is_active", false).Error; err != nil {
			log.Printf("Error during /stop: %v", err)
			bot.Send(tgbotapi.NewMessage(chatID, "Error stopping the bot. Please try again later."))
		} else {
			bot.Send(tgbotapi.NewMessage(chatID, "Subscription stopped. You will no longer receive updates."))
		}

	case "recent":
		var user models.User
		if err := db.Where("chat_id = ?", chatID).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				bot.Send(tgbotapi.NewMessage(chatID, "start the bot first"))
			} else {
				log.Printf("Error checking user for /recent: %v", err)
			}
			return
		}

		latest, err := database.GetLatestGoldPrice(db)
		var responseText string
		if err != nil {
			responseText = "Could not fetch latest price: " + err.Error()
		} else {
			ist := time.FixedZone("IST", 5.5*60*60)
			istTime := latest.Timestamp.In(ist)

			stats, err := database.GetGoldPriceStats(db)
			if err == nil {
				responseText = fmt.Sprintf(
					"Latest Gold Price: ₹%.2f\n"+
						"Time: %s (IST)\n\n"+
						"📈 Historical Stats:\n"+
						"• Week Avg:  ₹%.2f\n"+
						"• Week High: ₹%.2f | Low: ₹%.2f\n"+
						"• Month High: ₹%.2f | Low: ₹%.2f\n"+
						"• Year High:  ₹%.2f | Low: ₹%.2f\n"+
						"• All-Time High: ₹%.2f | Low: ₹%.2f",
					latest.Price,
					istTime.Format("2006-01-02 15:04:05"),
					stats.WeekAvg,
					stats.WeekHigh, stats.WeekLow,
					stats.MonthHigh, stats.MonthLow,
					stats.YearHigh, stats.YearLow,
					stats.AllTimeHigh, stats.AllTimeLow,
				)
			} else {
				log.Printf("Error generating historical stats for /recent: %v", err)
				responseText = fmt.Sprintf("Latest Gold Price: ₹%.2f\nTime: %s (IST)", latest.Price, istTime.Format("2006-01-02 15:04:05"))
			}
		}

		msg := tgbotapi.NewMessage(chatID, responseText)
		bot.Send(msg)
	}
}
