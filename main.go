package main

import (
	"hellobox/bot"
	"hellobox/database"
	"hellobox/env"
	"hellobox/router"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	go bot.HandleBot()
	go bot.HandlePartnerBot()
	database.InitialMigration()
	router.CreateRouter()
	http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization", "Token"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(env.Router))

}
