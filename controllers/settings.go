package controller

import (
	"encoding/json"
	"hellobox/database"
	"hellobox/models"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
)

func GetSettings(w http.ResponseWriter, r *http.Request) {
	settings := database.GetSettings()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func CreateSettings(w http.ResponseWriter, r *http.Request) {
	var settings models.ContactInfo
	err := json.NewDecoder(r.Body).Decode(&settings)
	if err != nil {
		error := models.Error{IsError: true, Message: "Unproccessable entity"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)
		return
	}
	database.CreateSettings(settings)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func EditSettings(w http.ResponseWriter, r *http.Request) {
	var settings models.ContactInfo
	err := json.NewDecoder(r.Body).Decode(&settings)
	if err != nil {
		error := models.Error{IsError: true, Message: "Unproccessable entity"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)
		return
	}
	database.EditSettings(settings)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func DeleteSettings(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		error := models.Error{IsError: true, Message: "Unproccessable entity"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)
		return
	}
	database.DeleteSettings(uint(id))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}
