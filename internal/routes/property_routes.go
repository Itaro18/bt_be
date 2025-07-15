package routes

import (
	"github.com/Itaro18/bt_be/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func PropertyRoutes(r chi.Router) {
	r.Get("/", handlers.GetProperties)
	r.Post("/", handlers.CreateProperty)
	r.Get("/{id}", handlers.GetPropertyByID)
}
