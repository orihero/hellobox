package controller

import (
	"encoding/json"
	"fmt"
	"hellobox/database"
	"hellobox/models"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
)

func GetPartner(w http.ResponseWriter, r *http.Request) {
	partner := database.GetPartner()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(partner)
}

func CreatePartner(w http.ResponseWriter, r *http.Request) {
	var partner models.Partner
	err := json.NewDecoder(r.Body).Decode(&partner)
	if err != nil {
		error := models.Error{IsError: true, Message: "Unproccessable entity"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)
		return
	}
	fmt.Println(partner)
	database.CreatePartner(partner)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(partner)
}

func EditPartner(w http.ResponseWriter, r *http.Request) {
	var partner models.Partner
	err := json.NewDecoder(r.Body).Decode(&partner)
	if err != nil {
		error := models.Error{IsError: true, Message: "Unproccessable entity"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)
		return
	}
	database.EditPartner(partner)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(partner)
}

func DeletePartner(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		error := models.Error{IsError: true, Message: "Unproccessable entity"}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)
		return
	}
	database.DeletePartner(uint(id))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}
