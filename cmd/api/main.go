package main

import (
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gold-price-bot.com/v2/internal/bot"
	"gold-price-bot.com/v2/internal/config"
	"gold-price-bot.com/v2/internal/database"
	"gold-price-bot.com/v2/internal/handlers"
	"gold-price-bot.com/v2/internal/worker"
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()
	if cfg.TelegramToken == "" {
		log.Fatal("API_KEY not found in environment")
	}

	// 2. Initialize Database
	db := database.InitDB(cfg.DatabasePath)

	// 3. Initialize Telegram Bot
	telegramBot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal("Error initializing Telegram bot:", err)
	}

	// 4. Start Background Workers
	worker.StartGoldPriceWorker(db, cfg.ScrapeInterval, cfg.GSTMultiplier)
	bot.StartTelegramBot(telegramBot, db)
	worker.StartTelegramBroadcastWorker(telegramBot, db, cfg.BroadcastInterval)

	// 5. Setup HTTP Handlers
	h := handlers.NewHandler(db)
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", h.Ping)
	mux.HandleFunc("/recent", h.Recent)

	// 6. Start Server
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
