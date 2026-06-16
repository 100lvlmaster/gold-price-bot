package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken  string
	DatabasePath   string
	GSTMultiplier  float64
	ScrapeInterval int
	BroadcastInterval int
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		TelegramToken:     getEnv("API_KEY", ""),
		DatabasePath:      getEnv("DB_PATH", "database.db"),
		GSTMultiplier:     getEnvAsFloat("GST_MULTIPLIER", 1.03),
		ScrapeInterval:    getEnvAsInt("SCRAPE_INTERVAL", 1),    // hours
		BroadcastInterval: getEnvAsInt("BROADCAST_INTERVAL", 2), // hours
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

func getEnvAsFloat(key string, fallback float64) float64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return fallback
}
