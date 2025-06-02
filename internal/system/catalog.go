package system

import (
	"net/http"
	"strings"

	"Factory/internal/system/rest"

	"github.com/go-chi/chi"
)

func catalog(r *chi.Mux, registry []route) {
	for _, route := range registry {
		r.Route(route.endpoint.path, func(r chi.Router) {
			r.Use(route.validator.validationHandler)

			var methods []string
			for _, method := range route.methods {
				methods = append(methods, method.method)
				switch method.method {
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
				methods := strings.Join(methods, ", ")
				w.Header().Set("Allow", methods)
				w.WriteHeader(204)
			})
		})
	}
}
