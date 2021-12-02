package database

import (
	"fmt"
	"hellobox/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDatabase() *gorm.DB {
	connection, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln("Invalid database url")
	}
	sql, err := connection.DB()

	err = sql.Ping()
	if err != nil {
		log.Fatal("Database connected")
	}
	fmt.Println("Database connection successuful.")
	return connection
}

func InitialMigration() {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	_ = connection.AutoMigrate(models.User{})
	_ = connection.AutoMigrate(models.Product{})
	_ = connection.AutoMigrate(models.Category{})
	_ = connection.AutoMigrate(models.Cart{})
	_ = connection.AutoMigrate(models.CartProduct{})
	_ = connection.AutoMigrate(models.News{})
	_ = connection.AutoMigrate(models.Partner{})
	_ = connection.AutoMigrate(models.ContactInfo{})
	_ = connection.AutoMigrate(models.Branch{})
	_ = connection.AutoMigrate(models.Option{})
	_ = connection.AutoMigrate(models.Order{})
}

//closes database connection
func CloseDatabase(connection *gorm.DB) {
	sqldb, _ := connection.DB()
	sqldb.Close()
}
