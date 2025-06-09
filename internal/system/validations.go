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

type entries map[string]string
type resolver func(string) string

func (e _endpoint) validationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if entries, err := e.validateRequest(r); err != nil {
			util.GetLogger(r).Error(err.Error())
			api.RequestErrorHandler(w, err)
		} else {
			ctx := context.WithValue(r.Context(), "entries", entries)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func (e _endpoint) validateRequest(r *http.Request) (entries, error) {
	var entries = make(entries)
	var method = _method{query: -1, headers: -1}

	if m, ok := registry.methods[e.id]; ok {
		method = m[r.Method]
	}

	var withPrefix = func(prefix string, err error) error {
		return errors.New(prefix + ": " + err.Error())
	}

	var uriResolver = func(s string) string { return chi.URLParam(r, s) }
	if err := entries.validateParameters(uriResolver, e.uriParams); err != nil {
		return nil, withPrefix("uri query", err)
	}

	if err := entries.validateParameters(r.Header.Get, method.headers); err != nil {
		return nil, withPrefix("headers", err)
	}

	if err := entries.validateParameters(r.URL.Query().Get, method.query); err != nil {
		return nil, withPrefix("query query", err)
	}

	return entries, nil
}

func (e entries) validateParameters(get resolver, params int) error {
	if params == -1 {
		return nil
	}

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
				e[name] = v
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

func validate(p _parameter, v string) error {
	var props = registry.properties[p.properties]

	var message = func(s string) error {
		prefix := " (" + p.name + ":" + v + ") "
		m := "failed to convert" + prefix + "to " + s
		return errors.New(m)
	}

	if enum, ok := props["enum"]; ok {
		values := strings.Split(enum, ",")
		if contains := slices.Contains(values, v); !contains {
			return errors.New(v + " must be in (" + enum + ")")
		}
	}

	switch p.typ {

	case "array":
		if items, ok := props["items"]; !ok {
			message := "array _property 'items' not defined for "
			return errors.New(message + p.name)
		} else {
			p.typ = items
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
