package database

import "hellobox/models"

func GetPartner() []models.Partner {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var partner []models.Partner
	connection.Find(&partner)
	return partner
}

func CreatePartner(partner models.Partner) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Create(&partner)
}

func EditPartner(partner models.Partner) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Save(&partner)
}

func DeletePartner(id uint) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Delete(&models.Partner{Id: id})
}
