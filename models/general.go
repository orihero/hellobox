package models

import (
	"time"

	"gorm.io/gorm"
)

type Error struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
}

type PresentImage struct {
	Id        uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	ImageUrl  string         `json:"image_url"`
}

type ProfitPercent struct {
	Id        uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Percent   uint           `json:"percent"`
}
