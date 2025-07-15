package routes

import (
	"github.com/Itaro18/bt_be/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router) {
	r.Post("/login", handlers.Login)
}
