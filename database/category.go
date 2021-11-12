package database

import "hellobox/models"

func GetCategories() []models.Category {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var categories []models.Category
	connection.Find(&categories)
	return categories
}

func CreateCategory(category models.Category) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Create(&category)
}

func EditCategory(category models.Category) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Save(&category)
}

func DeleteCategory(id uint) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Delete(&models.Category{Id: id})
}
