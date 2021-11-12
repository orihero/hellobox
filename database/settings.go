package database

import "hellobox/models"

func GetSettings() []models.ContactInfo {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var settings []models.ContactInfo
	connection.Find(&settings)
	return settings
}

func CreateSettings(settings models.ContactInfo) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Create(&settings)
}

func EditSettings(settings models.ContactInfo) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Save(&settings)
}

func DeleteSettings(id uint) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Delete(&models.ContactInfo{Id: id})
}
