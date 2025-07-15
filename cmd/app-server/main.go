package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Itaro18/bt_be/internal/db"
	"github.com/Itaro18/bt_be/internal/middleware"
	"github.com/Itaro18/bt_be/internal/routes"
	"github.com/Itaro18/bt_be/internal/utils"
	"github.com/go-chi/chi/v5"
	chi_m "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {

	utils.LoadEnv()
	var (
		Port        = os.Getenv("PORT")
		FrontendURL = os.Getenv("FRONTEND_URL")
	)
	err := db.Init()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	r := chi.NewRouter()
	r.Use(chi_m.RequestID)
	r.Use(chi_m.Logger)

	r.Use(cors.Handler(cors.Options{

		AllowedOrigins:   []string{FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, //
	}))
	r.Use(utils.SetJSONContentType)

	// Mount grouped routes
	r.Route("/users", routes.UserRoutes)
	r.Route("/auth", routes.AuthRoutes)
	r.Route("/bookings", func(r chi.Router) {
		r.Use(middleware.JWTAuth)
		routes.BookingRoutes(r)
	})

	r.Route("/properties", func(r chi.Router) {
		r.Use(middleware.JWTAuth)
		routes.PropertyRoutes(r)
	})
	r.Route("/customers", routes.CustomerRoutes)
	r.Route("/metadata", routes.MetadataRoutes)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})

	// At the bottom of main()
	log.Printf("✅ Server starting on %s", Port)

	if err := http.ListenAndServe(Port, r); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}

}
