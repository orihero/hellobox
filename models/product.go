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
	Id          uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	ProductId   uint           `json:"product_id"`
	Product     Product        `json:"product"`
	Count       uint           `json:"count"`
	CartId      uint           `json:"cart_id"`
	Token       string         `json:"token"`
	Utilized    int            `json:"utilized"`
	OptionIndex uint           `json:"option_index"`
	IsPresent   bool           `json:"is_present"`
}

type Option struct {
	Id        uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	ImageUrl  string         `json:"image_url"`
	ProductId uint           `json:"product_id"`
	Name      string         `json:"name"`
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
	Options       []Option       `json:"options"`
	PartnerId     uint           `json:"partner_id"`
}

type Category struct {
	Id        uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name      string         `json:"name"`
}

func (cart *Cart) GetProduct(productId uint) *CartProduct {
	if cart == nil {
		return nil
	}
	if len(cart.Products) > 0 {
		for _, el := range cart.Products {
			if el.ProductId == productId {
				return &el
			}
		}
	}
	return nil
}

func (cart *Cart) SetProduct(product CartProduct) {
	for i, el := range cart.Products {
		if el.Id == product.Id {
			ps := &cart.Products
			(*ps)[i] = product
		}
	}
}

func (cart *Cart) RemoveProduct(productId uint) {
	newProducts := []CartProduct{}
	for _, el := range cart.Products {
		if el.ProductId != productId {
			newProducts = append(newProducts, el)
		}
	}
	cart.Products = newProducts
}

func (cart *Cart) CartTotal() (total uint) {
	if cart == nil {
		return 0
	}
	for _, el := range cart.Products {
		total += el.Count * el.Product.Price
	}
	return total
}
