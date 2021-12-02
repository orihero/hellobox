package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	Id        uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	CartId    uint           `json:"cart_id"`
	Cart      Cart           `json:"cart"`
}
