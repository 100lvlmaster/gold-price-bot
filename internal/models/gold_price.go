package models

import "time"

type GoldPrice struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}
