package database

import "hellobox/models"

func GetProductsByCategory(categoryId uint) []models.Product {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var products []models.Product
	connection.Where(&models.Product{CategoryId: categoryId}).Find(&products)
	return products
}

func GetSingleProduct(id uint) models.Product {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var product models.Product
	connection.Where(&models.Product{Id: id}).First(&product)
	return product
}

func GetProducts() []models.Product {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var product []models.Product
	connection.Find(&product)
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
