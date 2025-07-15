package routes

import (
	"github.com/Itaro18/bt_be/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func UserRoutes(r chi.Router ) {
	
	//r.Get("/{id}", handlers.GetUser)
	r.Post("/", handlers.CreateUser)
	//r.Put("/{id}", handlers.UpdateUser)



}