package controller

import (
	"encoding/json"
	"hellobox/database"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	settings := database.GetSettings()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}
