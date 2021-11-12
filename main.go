package main

import (
	"hellobox/bot"
	"hellobox/database"
	"hellobox/env"
	"hellobox/router"
	"net/http"
)

func main() {
	go bot.HandleBot()
	database.InitialMigration()
	router.CreateRouter()
	http.ListenAndServe(":8081", env.Router)
}
