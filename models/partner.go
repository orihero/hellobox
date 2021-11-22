package models

import (
	"time"

	"gorm.io/gorm"
)

type Branch struct {
	Id        uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Address   string         `json:"adress"`
	City      string         `json:"city"`
	Region    string         `json:"region"`
	Contact   ContactInfo    `json:"contact"`
	ContactId uint           `json:"contact_id"`
	PartnerId uint           `json:"partner_id"`
}

type Partner struct {
	Id                   uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name                 string         `json:"name"`
	Categories           []Category     `json:"category" gorm:"many2many:partner_categories;"`
	Image                string         `json:"image"`
	DefaultProfitPercent uint           `json:"default_profit_percent"`
	Branches             []Branch       `json:"branches"`
}
