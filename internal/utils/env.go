package utils

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load(".env") // Adjust path as needed
		if err != nil {
			log.Fatalf("❌ Failed to load .env: %v", err)
		}
	}
}

func SetJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
