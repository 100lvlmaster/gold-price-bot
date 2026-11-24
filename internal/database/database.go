package database

import (
	"log"
	"sync"
	"time"

	"github.com/glebarez/sqlite"
	"gold-price-bot.com/v2/internal/config"
	"gold-price-bot.com/v2/internal/models"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func InitDB(path string) {
	once.Do(func() {
		conn, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}

		// Restrict SQLite to a single underlying connection to avoid "database is locked" errors.
		sqlDB, err := conn.DB()
		if err != nil {
			log.Fatal("Failed to get underlying database handle:", err)
		}
		sqlDB.SetMaxOpenConns(1)

		// Run Migrations
		if err := conn.AutoMigrate(&models.User{}, &models.GoldPrice{}, &models.Setting{}); err != nil {
			log.Fatal("Failed to run migrations:", err)
		}

		db = conn
	})
}

func SaveGoldPrice(goldPrice models.GoldPrice) error {
	return db.Create(&goldPrice).Error
}

func GetLatestGoldPrice() (models.GoldPrice, error) {
	var latestPrice models.GoldPrice
	result := db.Order("timestamp desc").First(&latestPrice)
	return latestPrice, result.Error
}

func GetLatestGoldPriceIfChanged() (models.GoldPrice, bool, error) {
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

type GoldPriceStats struct {
	WeekAvg     float64
	WeekLow     float64
	WeekHigh    float64
	MonthLow    float64
	MonthHigh   float64
	YearLow     float64
	YearHigh    float64
	AllTimeLow  float64
	AllTimeHigh float64
}

func GetGoldPriceStats() (GoldPriceStats, error) {
	now := time.Now().UTC()
	oneWeekAgo := now.AddDate(0, 0, -7)
	oneMonthAgo := now.AddDate(0, 0, -30)
	oneYearAgo := now.AddDate(0, 0, -365)

	var stats GoldPriceStats

	// 1. Week stats
	var week struct {
		Avg float64 `gorm:"column:avg"`
		Min float64 `gorm:"column:min"`
		Max float64 `gorm:"column:max"`
	}
	if err := db.Model(&models.GoldPrice{}).Where("timestamp >= ?", oneWeekAgo).Select("COALESCE(AVG(price), 0) as avg, COALESCE(MIN(price), 0) as min, COALESCE(MAX(price), 0) as max").Scan(&week).Error; err != nil {
		return stats, err
	}
	stats.WeekAvg = week.Avg
	stats.WeekLow = week.Min
	stats.WeekHigh = week.Max

	// 2. Month stats
	var month struct {
		Min float64 `gorm:"column:min"`
		Max float64 `gorm:"column:max"`
	}
	if err := db.Model(&models.GoldPrice{}).Where("timestamp >= ?", oneMonthAgo).Select("COALESCE(MIN(price), 0) as min, COALESCE(MAX(price), 0) as max").Scan(&month).Error; err != nil {
		return stats, err
	}
	stats.MonthLow = month.Min
	stats.MonthHigh = month.Max

	// 3. Year stats
	var year struct {
		Min float64 `gorm:"column:min"`
		Max float64 `gorm:"column:max"`
	}
	if err := db.Model(&models.GoldPrice{}).Where("timestamp >= ?", oneYearAgo).Select("COALESCE(MIN(price), 0) as min, COALESCE(MAX(price), 0) as max").Scan(&year).Error; err != nil {
		return stats, err
	}
	stats.YearLow = year.Min
	stats.YearHigh = year.Max

	// 4. All time stats
	var allTime struct {
		Min float64 `gorm:"column:min"`
		Max float64 `gorm:"column:max"`
	}
	if err := db.Model(&models.GoldPrice{}).Select("COALESCE(MIN(price), 0) as min, COALESCE(MAX(price), 0) as max").Scan(&allTime).Error; err != nil {
		return stats, err
	}
	stats.AllTimeLow = allTime.Min
	stats.AllTimeHigh = allTime.Max

	return stats, nil
}

func UpsertUser(chatID int64, name string) (models.User, error) {
	user := models.User{
		ChatID:   chatID,
		Name:     name,
		IsActive: true,
	}
	err := db.Where(models.User{ChatID: chatID}).Assign(models.User{IsActive: true, Name: name}).FirstOrCreate(&user).Error
	return user, err
}

func GetUser(chatID int64) (models.User, error) {
	var user models.User
	err := db.Where("chat_id = ?", chatID).First(&user).Error
	return user, err
}

func SetUserActive(chatID int64, active bool) error {
	return db.Model(&models.User{}).Where("chat_id = ?", chatID).Update("is_active", active).Error
}

func IsAdmin(chatID int64) bool {
	user, err := GetUser(chatID)
	if err != nil {
		return false
	}
	return user.IsAdmin
}

// GetBroadcastUsers returns the users that should receive price broadcasts.
// Admins are always included. When adminOnly is true, only admins are returned.
func GetBroadcastUsers() ([]models.User, error) {
	var query *gorm.DB
	if config.IsDevMode {
		query = db.Where("is_active = ? AND is_admin = ?", true, true)
	} else {
		query = db.Where("is_active = ?", true)
	}

	var users []models.User
	err := query.Find(&users).Error
	return users, err
}
