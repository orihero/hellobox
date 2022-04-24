package controller

import (
	"encoding/json"
	"hellobox/database"
	"net/http"
)

func GetOrders(w http.ResponseWriter, r *http.Request) {
	category := database.GetOrders()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}
