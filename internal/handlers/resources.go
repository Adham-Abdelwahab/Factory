package handlers

import (
	"github.com/go-chi/chi"
)

func ResourceHandler(r *chi.Mux) {
	r.Route("/resources", func(r chi.Router) {

	})
}
