package database

import (
	"fmt"
	"hellobox/models"
)

func GetUsers() []models.User {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var user []models.User
	connection.Find(&user)
	return user
}

func CreateUser(user models.User) *models.Error {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var usr models.User
	usr.Phone = user.Phone
	connection.Where(usr).Find(&usr)
	fmt.Println(usr)
	
	if usr.UserId != 0 {
		var err models.Error = models.Error{IsError: true, Message: "User already exist"}
		return &err
	}
	connection.Create(&user)
	return nil
}

func EditUser(user models.User) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Save(&user)
}

func DeleteUser(id uint) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Delete(&models.User{Id: id})
}
