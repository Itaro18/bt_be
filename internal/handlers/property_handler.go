package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Itaro18/bt_be/internal/db"
	"github.com/Itaro18/bt_be/internal/models"
	"github.com/go-chi/chi/v5"
)

func CreateProperty(w http.ResponseWriter, r *http.Request) {
	var prop models.Property
	if err := json.NewDecoder(r.Body).Decode(&prop); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := db.DB.Create(&prop).Error; err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(prop)

}

func GetProperties(w http.ResponseWriter, r *http.Request) {
	var properties []models.Property

	// Fetch all properties from the database
	if err := db.DB.Find(&properties).Error; err != nil {
		http.Error(w, "Failed to fetch properties", http.StatusInternalServerError)
		return
	}

	// Return properties as JSON
	if err := json.NewEncoder(w).Encode(properties); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetPropertyByID(w http.ResponseWriter, r *http.Request) {
	var prop models.Property
	id := chi.URLParam(r, "id")

	if err := db.DB.Where("property_id = ?", id).Find(&prop).Error; err != nil {
		http.Error(w, "Failed to fetch properties", http.StatusInternalServerError)
		return
	}

	// Return properties as JSON
	if err := json.NewEncoder(w).Encode(prop); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
