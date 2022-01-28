package database

import "hellobox/models"

func GetNews() []models.News {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var news []models.News
	connection.Preload("Partner").Preload("Product").Find(&news)
	return news
}

func CreateNews(news models.News) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Create(&news)
}

func EditNews(news models.News) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Save(&news)
}

func DeleteNews(id uint) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Delete(&models.News{Id: id})
}

func GetNewsByid(id uint) models.News {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var news models.News
	connection.Where(models.News{Id: id}).Preload("Partner").Preload("Product").First(&news)
	return news
}
