package system

import (
	"maps"
	"net/http"
	"slices"
	"strings"

	"Factory/internal/system/rest"

	"github.com/go-chi/chi"
)

func catalog(r *chi.Mux) {
	for _, endpoint := range registry.endpoints {
		r.Route(endpoint.path, func(r chi.Router) {
			r.Use(endpoint.validationHandler)

			verbs := maps.Keys(registry.methods[endpoint.methods])
			methods := slices.Collect(verbs)
			for _, verb := range methods {
				switch verb {
				case "GET":
					r.Get("/", rest.Get)
				case "POST":
					r.Post("/", rest.Post)
				case "PUT":
					r.Put("/", rest.Put)
				case "DELETE":
					r.Delete("/", rest.Delete)
				case "PATCH":
					r.Patch("/", rest.Patch)
				}
			}

			r.Options("/", func(w http.ResponseWriter, _ *http.Request) {
				allowed := strings.Join(methods, ", ")
				w.Header().Set("Allow", allowed)
				w.WriteHeader(204)
			})
		})
	}
}
