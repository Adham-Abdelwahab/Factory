package system

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"Factory/api"
	"Factory/internal/util"

	"github.com/go-chi/chi"
)

/* Parameter Validation */
type validation struct {
	properties []property
	uriParams  []parameter
	headers    []parameter
	query      []parameter
}

/* Method --> Validation */
type validator map[string]validation

func (v validator) validationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if parameters, err := validateRequest(r, v[r.Method]); err != nil {
			util.GetLogger(r).Error(err.Error())
			api.RequestErrorHandler(w, err)
		} else {
			ctx := context.WithValue(r.Context(), "parameters", parameters)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func validateRequest(r *http.Request, v validation) (map[string]string, error) {
	var params = make(map[string]string)

	var uriParams = func(s string) string { return chi.URLParam(r, s) }
	if err := validateParameters(params, v.uriParams, uriParams); err != nil {
		return nil, errors.New("uri parameters: " + err.Error())
	}

	if err := validateParameters(params, v.headers, r.Header.Get); err != nil {
		return nil, errors.New("headers: " + err.Error())
	}

	if err := validateParameters(params, v.query, r.URL.Query().Get); err != nil {
		return nil, errors.New("query parameters: " + err.Error())
	}

	return params, nil
}

func validateParameters(params map[string]string, parameters []parameter, get func(string) string) error {
	var missing []string
	var issues []string

	for _, p := range parameters {
		if v := get(p.name); v == "" {
			missing = append(missing, p.name)
		} else {
			if err := validate(p, v); err != nil {
				issues = append(issues, err.Error())
			} else {
				params[p.name] = v
			}
		}
	}

	if len(missing) != 0 {
		missing := strings.Join(missing, ", ")
		return errors.New(missing + " must be provided")
	}

	if len(issues) != 0 {
		return errors.New(strings.Join(issues, ". "))
	}

	return nil
}

func validate(p parameter, v string) error {
	return nil
}
