package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	ChatID    int64     `gorm:"unique" json:"chat_id"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	IsAdmin   bool      `gorm:"default:false" json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}
