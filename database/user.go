package database

import (
	"hellobox/models"
)

func GetUsers() []models.User {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var users []models.User
	connection.Find(&users)
	return users
}

func GetUser(userId uint) models.User {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	user := models.User{Id: userId}
	connection.Where(&user).Preload("Cart").Find(&user)
	return user
}

func FilterUser(usr models.User) models.User {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Where(&usr).Preload("Cart").Preload("Cart.Products").Preload("Cart.Products.Product").Preload("Cart.Products.Product.Options").Find(&usr)
	return usr
}
func CreateUser(user models.User) *models.Error {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	var usr models.User
	usr.Phone = user.Phone
	connection.Where(usr).Find(&usr)
	if usr.TgId != 0 {
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
	connection.Save(&user.Cart)
	if user.Cart != nil && len(user.Cart.Products) > 0 {
		for _, el := range user.Cart.Products {
			connection.Save(&el)
		}
	}
}

func ClearUserCart(user models.User) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Model(&user).Association("Cart").Delete(user.Cart)
}

func DeleteUser(id uint) {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	connection.Delete(&models.User{Id: id})
}
