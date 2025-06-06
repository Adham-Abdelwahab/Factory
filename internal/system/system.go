package system

import (
	"github.com/go-chi/chi"
)

func system(r *chi.Mux) {
	r.Get("/system/endpoints", GetSystemEndpoints)
	r.Get("/system/endpoints/{endpoint}", GetSystemEndpointById)
	r.Get("/system/endpoints/{endpoint}/{method}", GetSystemMethodById)

	r.Get("/system/parameters", GetSystemParameters)
	r.Get("/system/parameters/{parameter}", GetSystemParameterById)

	r.Get("/system/properties", GetSystemProperties)
	r.Get("/system/properties/{property}", GetSystemPropertyById)
}
