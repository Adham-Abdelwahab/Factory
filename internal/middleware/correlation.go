package middleware

import (
	"net/http"

	"Factory/internal/util"

	"github.com/google/uuid"
)

func Correlation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id string
		if id = r.Header.Get("X-Correlation-ID"); id == "" {
			id = uuid.New().String()
		}

		w.Header().Set("X-Correlation-ID", id)
		r = util.SetLogger(r, id)
		next.ServeHTTP(w, r)
	})
}
