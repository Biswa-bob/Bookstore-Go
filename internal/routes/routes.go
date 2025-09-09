package routes

import (
	"net/http"

	"github.com/Biswa-bob/bookstore/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)

		r.Post("/books", app.Middleware.RequireUser(app.BookHandler.HandleCreateBook))
		r.Get("/books", app.Middleware.RequireUser(app.BookHandler.HandleGetBooks))
		r.Get("/books/{id}", app.Middleware.RequireUser(app.BookHandler.HandleGetBookById))
		r.Patch("/books/{id}", app.Middleware.RequireUser(app.BookHandler.HandleUpdateBookById))
		r.Delete("/books/{id}", app.Middleware.RequireUser(app.BookHandler.HandleDeleteBookById))
		r.Get("/auth/profile", app.Middleware.RequireUser(app.UserHandler.HandleGetUserProfile))
	})

	r.Get("/health", app.HealthCheck)

	r.Post("/auth/register", app.UserHandler.HandleRegisterUser)
	r.Post("/auth/login", app.TokenHandler.HandleCreateToken)
	r.Post("/auth/refresh", func(w http.ResponseWriter, r *http.Request) {})
	r.Post("/auth/logout", func(w http.ResponseWriter, r *http.Request) {})

	return r
}
