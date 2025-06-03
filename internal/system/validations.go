package system

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"Factory/api"
	"Factory/internal/util"

	"github.com/go-chi/chi"
)

/* Parameter Validation */
type validation struct {
	properties map[int]map[string]string
	uriParams  []parameter
	headers    []parameter
	query      []parameter
}

/* Method --> Validation */
type validator map[string]validation
type entries map[string]string
type resolver func(string) string

func (v validator) validationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if entries, err := validateRequest(r, v[r.Method]); err != nil {
			util.GetLogger(r).Error(err.Error())
			api.RequestErrorHandler(w, err)
		} else {
			ctx := context.WithValue(r.Context(), "entries", entries)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func validateRequest(r *http.Request, v validation) (entries, error) {
	var entries = make(entries)

	var uriParams = func(s string) string { return chi.URLParam(r, s) }
	if err := v.validateParameters(entries, v.uriParams, uriParams); err != nil {
		return nil, errors.New("uri parameters: " + err.Error())
	}

	if err := v.validateParameters(entries, v.headers, r.Header.Get); err != nil {
		return nil, errors.New("headers: " + err.Error())
	}

	if err := v.validateParameters(entries, v.query, r.URL.Query().Get); err != nil {
		return nil, errors.New("query parameters: " + err.Error())
	}

	return entries, nil
}

func (validation validation) validateParameters(entries entries, parameters []parameter, get resolver) error {
	var missing []string
	var issues []string

	for _, p := range parameters {
		if v := get(p.name); v == "" {
			if p.required {
				missing = append(missing, p.name)
			}
		} else {
			props := validation.properties[p.properties]
			if err := validate(props, p, v); err != nil {
				issues = append(issues, err.Error())
			} else {
				entries[strings.ToLower(p.name)] = v
			}
		}
	}

	if len(missing) != 0 {
		missing := strings.Join(missing, ", ")
		return errors.New(missing + " must be provided")
	}

	if len(issues) != 0 {
		issues := strings.Join(issues, ". ")
		return errors.New(issues)
	}

	return nil
}

func validate(props map[string]string, p parameter, v string) error {
	var message = func(s string) error {
		prefix := " (" + p.name + ":" + v + ") "
		m := "failed to convert" + prefix + "to" + s
		return errors.New(m)
	}

	if enum, ok := props["enum"]; ok {
		values := strings.Split(enum, ",")
		if contains := slices.Contains(values, v); !contains {
			return errors.New(v + " must be value " + enum)
		}
	}

	switch p.typ {

	case "array":
		if items, ok := props["items"]; !ok {
			message := "array property 'items' not defined for "
			return errors.New(message + p.name)
		} else {
			p.typ = items
			var issues []string
			values := strings.Split(v, ",")
			for _, v := range values {
				if err := validate(props, p, v); err != nil {
					issues = append(issues, err.Error())
				}
			}
			if len(issues) != 0 {
				message := strings.Join(issues, ". ")
				return errors.New(message)
			}
		}

	case "integer":
		if num, err := strconv.Atoi(v); err != nil {
			return message("int")
		} else {
			return numeric(props, num)
		}

	case "boolean":
		if _, err := strconv.ParseBool(v); err != nil {
			return message("bool")
		}
	}

	return nil
}

func numeric(props map[string]string, number int) error {
	return nil
}
