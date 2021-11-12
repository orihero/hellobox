package controller

import (
	"encoding/json"
	"hellobox/database"
	"hellobox/models"
	"log"
	"strconv"

	"hellobox/bot"
	"net/http"

	"github.com/gorilla/mux"
)

func GetNews(w http.ResponseWriter, r *http.Request) {
	news := database.GetNews()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

func CreateNews(w http.ResponseWriter, r *http.Request) {
	var news models.News
	err := json.NewDecoder(r.Body).Decode(&news)
	if err != nil {
		error := models.Error{IsError: true, Message: "Unproccessable entity"}
		log.Fatal(err)
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)
		return
	}
	bot.SendNews(news)
	database.CreateNews(news)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

func EditNews(w http.ResponseWriter, r *http.Request) {
	var news models.News
	err := json.NewDecoder(r.Body).Decode(&news)
	if err != nil {
		error := models.Error{IsError: true, Message: "Unproccessable entity"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)
		return
	}
	database.EditNews(news)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

func DeleteNews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		error := models.Error{IsError: true, Message: "Unproccessable entity"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)
		return
	}
	database.DeleteNews(uint(id))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}
