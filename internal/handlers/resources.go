package handlers

import (
	"Factory/internal/handlers/resources"
	"Factory/internal/middleware"

	"github.com/go-chi/chi"
)

func ResourceHandler(r *chi.Mux) {
	r.Route("/resource", func(r chi.Router) {
		r.Use(middleware.Registry)

		r.Get("/res", resources.GetResource)
	})
}
