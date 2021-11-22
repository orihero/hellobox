package models

import (
	"time"

	"gorm.io/gorm"
)

type ContactInfo struct {
	Id           uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	FacebookUrl  string         `json:"facebook_url"`
	InstagramUrl string         `json:"instagram_url"`
	TelegramUrl  string         `json:"telegram_url"`
	EmailUrl     string         `json:"email_url"`
	ContactInfo  string         `json:"contact_info"`
	PhoneNumber  string         `json:"phone_number"`
	SiteUrl      string         `json:"site_url"`
}
