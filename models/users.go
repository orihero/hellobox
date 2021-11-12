package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id        uint           `gorm:"AUTO_INCREMENT;PRIMARY_KEY;not null" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Phone     string         `json:"phone"`
	Firstname string         `json:"first_name"`
	Lastname  string         `json:"last_name"`
	ImageUrl  string         `json:"image_url"`
	Cart      *Cart          `json:"cart"`
	CartId    uint           `json:"cart_id"`
	ChatId    int64          `json:"chat_id"`
	UserId    int64          `json:"user_id"`
}
