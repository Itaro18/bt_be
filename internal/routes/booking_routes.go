package routes

import (
	"github.com/Itaro18/bt_be/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func BookingRoutes(r chi.Router) {
	r.Get("/", handlers.GetBookings)
	r.Post("/createbooking", handlers.CreateBooking)
	r.Get("/property/{id}", handlers.GetBookingsByPropertyID)
	r.Get("/properties/today-status", handlers.GetPropertiesTodayStatus)
	r.Get("/{id}", handlers.GetBookingsByID)
	r.Put("/update-booking/{id}", handlers.UpdateBooking)
	r.Delete("/delete/{id}", handlers.DeleteBooking)

}
