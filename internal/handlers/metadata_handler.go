package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Itaro18/bt_be/internal/db"
	"github.com/Itaro18/bt_be/internal/models"
)

func GetMetadata(w http.ResponseWriter, r *http.Request) {
	var properties []models.Property
	if err := db.DB.Find(&properties).Error; err != nil {
		http.Error(w, "Failed to fetch properties", http.StatusInternalServerError)
		return
	}

	citySet := make(map[string]bool)
	cityPropertyMap := make(map[string][]string)

	for _, p := range properties {
		citySet[p.City] = true
		cityPropertyMap[p.City] = append(cityPropertyMap[p.City], p.Name)
	}

	cities := make([]string, 0, len(citySet))
	for city := range citySet {
		cities = append(cities, city)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"cities":          cities,
		"cityPropertyMap": cityPropertyMap,
	})
}
