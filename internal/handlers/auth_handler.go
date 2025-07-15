package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Itaro18/bt_be/internal/db"
	"github.com/Itaro18/bt_be/internal/models"
	"github.com/Itaro18/bt_be/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type loginInput struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var in loginInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	var u models.User
	if err := db.DB.Where("phone = ?", in.Phone).First(&u).Error; err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(in.Password)); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(u.UserID)
	if err != nil {
		log.Printf("❌ Error generating JWT: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
