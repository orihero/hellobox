package database

import (
	"hellobox/models"
)

func GetCartProductsById(productId uint) *models.CartProduct {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var product models.CartProduct
	connection.Where(&models.CartProduct{Id: productId}).Preload("Product").Preload("Product.Options").Preload("Product.Category").First(&product)
	return &product
}

func GetProductsByCategory(categoryId uint) []models.Product {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var products []models.Product
	connection.Where(&models.Product{CategoryId: categoryId}).Preload("Options").Preload("Category").Find(&products)
	return products
}

func GetProductsByPartner(partnerId uint) []models.Product {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var products []models.Product
	connection.Where(&models.Product{PartnerId: partnerId}).Preload("Partner").Preload("Options").Preload("Category").Find(&products)
	println("---------------------------")
	println(partnerId)
	println(products)
	println("---------------------------")
	return products
}

func GetTopProducts() []models.Product {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var products []models.Product
	connection.Where(&models.Product{IsTop: 1}).Preload("Partner").Preload("Options").Preload("Category").Find(&products)
	return products
}

func GetCartProductsByToken(token string) *models.CartProduct {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var product models.CartProduct
	connection.Where(&models.CartProduct{Token: token}).Preload("Product").Preload("Product.Options").Preload("Product.Category").First(&product)
	return &product
}

func GetSingleProduct(id uint) models.Product {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var product models.Product
	connection.Where(&models.Product{Id: id}).Preload("Options").Preload("Category").First(&product)
	return product
}

func GetProducts() []models.Product {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var product []models.Product
	connection.Preload("Options").Preload("Category").Find(&product)
	return product
}

func CreateProduct(product models.Product) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Create(&product)
}

func EditProduct(product models.Product) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Save(&product)
}

func EditCartProduct(product models.CartProduct) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	// id := fmt.Sprint("id=%d", product.Id)
	// db := connection.Model(&product).Updates(map[string]interface{}{
	// 	"IsPresent":   product.IsPresent,
	// 	"Product":     product.Product,
	// 	"ProductId":   product.ProductId,
	// 	"Count":       product.Count,
	// 	"Id":          product.Id,
	// 	"OptionIndex": product.OptionIndex,
	// 	"Token":       product.Token,
	// 	"Utilized":    product.Utilized,
	// 	"UpdatedAt":   product.UpdatedAt,
	// 	"DeletedAt":   product.DeletedAt,
	// 	"CartId":      product.CartId,
	// })
	connection.Model(&product).Select("is_present").Updates(map[string]interface{}{"IsPresent": product.IsPresent})
	connection.Save(&product)
}

func DeleteProduct(id uint) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Delete(&models.Product{Id: id})
}

func DeleteCartProduct(id uint) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Delete(&models.CartProduct{Id: id})
}

func UpdatePresentImage(img models.PresentImage) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	if connection.Model(&img).Where("id = ?", img.Id).Updates(&img).RowsAffected == 0 {
		connection.Create(&img)
	}
}

func GetPresentImage() models.PresentImage {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var p models.PresentImage
	connection.First(&p)
	return p
}

func GetProfitPercent() models.ProfitPercent {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var p models.ProfitPercent
	connection.First(&p)
	return p
}

func UpdateProfitPercent(profit models.ProfitPercent) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	if connection.Model(&profit).Where("id = ?", profit.Id).Updates(&profit).RowsAffected == 0 {
		connection.Create(&profit)
	}
}

//  func GetProductsbyPartner (id models.Partner){
// 	connection := GetDatabase()
// 	var partner models.Partner
// 	connection.Where(&models.Partner{Id: id}).First(&partner)
// 	return partner
//  }
