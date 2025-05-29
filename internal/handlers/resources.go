package handlers

import (
	"Factory/internal/handlers/resources"
	"Factory/internal/middleware"

	"github.com/go-chi/chi"
)

func ResourceHandler(r *chi.Mux) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.Registry)

		r.Get("/resource", resources.GetResource)
		r.Post("/resource", resources.GetResource)
	})
}
