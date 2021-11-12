package models

import (
	"time"

	"gorm.io/gorm"
)

type News struct {
	Id          uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Description string         `json:"description"`
	ImageUrl    string         `json:"image_url"`
	Partner     Partner        `json:"partner"`
	PartnerId   uint           `json:"partner_id"`
	Product     *Product       `json:"product"`
	ProductId   uint           `json:"product_id"`
}
