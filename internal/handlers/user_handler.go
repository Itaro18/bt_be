package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/Itaro18/bt_be/internal/db"
	"github.com/Itaro18/bt_be/internal/models"
	"github.com/Itaro18/bt_be/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// createUserInput for binding only relevant fields from the request
type createUserInput struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	utils.LoadEnv()
	Pass := os.Getenv("PASS")
	var input createUserInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate or hash a default password
	rawPassword := Pass // Can also be random or email-based
	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Password generation failed", http.StatusInternalServerError)
		return
	}

	// Prepare user for DB insertion
	user := models.User{
		Name:     input.Name,
		Phone:    input.Phone,
		Password: string(hashed),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": user.UserID,
		"name":    user.Name,
		"phone":   user.Phone,
		"note":    "Default password set on server",
	})
}
