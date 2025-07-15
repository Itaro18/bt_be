package routes

import (
	"github.com/Itaro18/bt_be/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func MetadataRoutes(r chi.Router) {
	r.Get("/md", handlers.GetMetadata)
}
