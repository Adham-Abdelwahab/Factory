package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

func RailwayHandler(r *chi.Mux) {
	r.Route("/railway", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(struct {
				Response string
			}{
				Response: "Railway request",
			})
		})
	})
}
