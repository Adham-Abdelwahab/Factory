package middleware

import (
	"Factory/api"
	"Factory/internal/util"
	"errors"
	"fmt"
	"net/http"
)

func Registry(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var log = util.GetLogger(r)
		var resource = r.URL.Query().Get("resource")
		var err error

		fields := []string{"resource", "notresource"}
		err = util.ValidateParameters(r, util.Query, fields)
		if err != nil {
			log.Error(err)
			api.RequestErrorHandler(w, err)
			return
		}

		var database *util.DatabaseInterface
		database, err = util.NewDatabase()
		if err != nil {
			err = errors.New("failed to connect to the database")
			log.Error(err)
			api.ConnectionErrorHandler(w, err)
			return
		}

		resourceDetails := (*database).GetResourceDetails(resource)
		if resourceDetails == nil {
			message := fmt.Sprintf("no information found for resource '%s'", resource)
			err = errors.New(message)
			log.Error(err)
			api.RequestErrorHandler(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
