package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Correlation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id = r.Header.Get("X-Correlation-ID")
		if id == "" {
			id = uuid.New().String()
		}
		w.Header().Set("X-Correlation-ID", id)

		entry := logrus.WithField("id", id)
		ctx := context.WithValue(r.Context(), "logger", entry)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
