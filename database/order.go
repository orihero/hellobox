package database

import "hellobox/models"

func GetOrderHistory(userId uint) []models.Order {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var orders []models.Order
	connection.Where(models.Order{UserId: userId}).Preload("Cart").Preload("Cart.Products").Preload("Cart.Products.Product").Preload("Cart.Products.Product.Options").Preload("Cart.Products.Product.Category").Find(&orders)
	return orders
}

func GetOrders() []models.Order {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var orders []models.Order
	connection.Preload("Cart").Preload("Cart.Products").Preload("Cart.Products.Product").Preload("Cart.Products.Product.Options").Preload("Cart.Products.Product.Category").Find(&orders)
	return orders
}

func CreateOrder(orders models.Order) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Create(&orders)
}

func EditOrder(orders models.Order) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Save(&orders)
}

func DeleteOrder(id uint) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Delete(&models.Order{Id: id})
}
