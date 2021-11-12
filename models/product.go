package models

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	Id        uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Products  []CartProduct  `json:"products"`
}

type CartProduct struct {
	Id        uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	ProductId uint           `json:"product_id"`
	Product   Product        `json:"product"`
	Count     uint           `json:"count"`
	CartId    uint           `json:"cart_id"`
	Token     string
}

type Product struct {
	Id            uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Category      Category       `json:"category"`
	CategoryId    uint           `json:"category_id"`
	Price         uint           `json:"price"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	ImageUrl      string         `json:"image_url"`
	ExpiresIn     uint           `json:"expires_in"` //hours
	ProfitPercent uint           `json:"profit_percent"`
}

type Category struct {
	Id        uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name      string         `json:"name"`
}
