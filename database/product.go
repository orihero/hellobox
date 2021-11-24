package database

import "hellobox/models"

func GetProductsByCategory(categoryId uint) []models.Product {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var products []models.Product
	connection.Where(&models.Product{CategoryId: categoryId}).Preload("Options").Preload("Category").Find(&products)
	return products
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

func DeleteProduct(id uint) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Delete(&models.Product{Id: id})
}

//  func GetProductsbyPartner (id models.Partner){
// 	connection := GetDatabase()
// 	var partner models.Partner
// 	connection.Where(&models.Partner{Id: id}).First(&partner)
// 	return partner
//  }
