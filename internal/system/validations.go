package system

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"Factory/internal/util"
	"Factory/models"

	"github.com/go-chi/chi"
)

type entries map[string]map[string]string
type resolver func(string) string

func (e _endpoint) validationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if entries, err := e.validateRequest(r); err != nil {
			models.RequestErrorHandler(w, r, err.Error())
		} else {
			ctx := context.WithValue(r.Context(), "entries", entries)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func (e _endpoint) validateRequest(r *http.Request) (entries, error) {
	var method = registry.methods[e.id][r.Method]
	var entries = make(entries)

	var withPrefix = func(prefix string, err error) error {
		return errors.New(prefix + ": " + err.Error())
	}

	var uriResolver = func(s string) string { return chi.URLParam(r, s) }
	if params, err := validateParameters(uriResolver, e.uriParams); err != nil {
		return nil, withPrefix("uri parameters", err)
	} else {
		entries["uri"] = params
	}

	if params, err := validateParameters(r.Header.Get, method.headers); err != nil {
		return nil, withPrefix("headers", err)
	} else {
		entries["headers"] = params
	}

	if params, err := validateParameters(r.URL.Query().Get, method.query); err != nil {
		return nil, withPrefix("query parameters", err)
	} else {
		entries["query"] = params
	}

	return entries, nil
}

func validateParameters(get resolver, params int) (map[string]string, error) {
	if params == 0 {
		return nil, nil
	}

	var entries = make(map[string]string)
	var missing []string
	var issues []string

	for name, p := range registry.parameters[params] {
		if v := get(name); v == "" {
			if p.required {
				missing = append(missing, name)
			}
		} else {
			if err := validate(p, v); err != nil {
				issues = append(issues, err.Error())
			} else {
				entries[name] = v
			}
		}
	}

	if len(missing) != 0 {
		missing := strings.Join(missing, ", ")
		return nil, errors.New(missing + " must be provided")
	}

	if len(issues) != 0 {
		issues := strings.Join(issues, ". ")
		return nil, errors.New(issues)
	}

	return entries, nil
}

func validate(p _parameter, v string) error {
	var props = registry.properties[p.properties]

	var conversion = func(s string) error {
		message := "failed to convert (%s:%s) to %s"
		message = util.Message(message, p.name, v, s)
		return errors.New(message)
	}

	if enum, ok := props["enum"]; ok {
		values := strings.Split(enum, ",")
		if !slices.Contains(values, v) {
			message := "%s (%s) must be in (%s)"
			message = util.Message(message, p.name, v, enum)
			return errors.New(message)
		}
	}

	switch p.typ {

	case "array":
		if items, ok := props["items"]; ok {
			p.typ = items
		} else {
			message := "array property 'items' not defined for %s"
			message = util.Message(message, p.name)
			return errors.New(message)
		}

		var issues []string
		values := strings.Split(v, ",")
		for _, v := range values {
			if err := validate(p, v); err != nil {
				issues = append(issues, err.Error())
			}
		}

		if len(issues) != 0 {
			message := strings.Join(issues, ". ")
			return errors.New(message)
		}

	case "integer":
		if num, err := strconv.Atoi(v); err != nil {
			return conversion("int")
		} else {
			return numeric(props, num)
		}

	case "boolean":
		if _, err := strconv.ParseBool(v); err != nil {
			return conversion("bool")
		}
	}

	return nil
}

func numeric(props map[string]string, number int) error {
	return nil
}
