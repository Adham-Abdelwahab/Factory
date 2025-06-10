package system

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"Factory/api"
	"Factory/internal/util"

	"github.com/go-chi/chi"
)

func isId(value string) (int, error) {
	if number, err := strconv.Atoi(value); err != nil {
		message := " is not a valid id"
		err = errors.New(value + message)
		return 0, err
	} else {
		return number, nil
	}
}

func notFound(w http.ResponseWriter, r *http.Request, args ...string) {
	err := errors.New(strings.Join(args, " "))
	util.GetLogger(r).Error(err)
	api.NotFoundErrorHandler(w, err)
}

func GetSystemEndpoints(w http.ResponseWriter, r *http.Request) {
	var base = r.URL.Query().Get("basePath")
	var endpoints = make(map[string]int)

	for i, e := range registry.endpoints {
		if strings.HasPrefix(e.path, base) {
			endpoints[e.path] = i
		}
	}

	json.NewEncoder(w).Encode(endpoints)
}

func GetSystemEndpointById(w http.ResponseWriter, r *http.Request) {
	var endpoint = chi.URLParam(r, "endpoint")

	id, err := isId(endpoint)
	if err != nil {
		util.GetLogger(r).Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	route, ok := registry.endpoints[id]
	if !ok {
		message := "no endpoint found with id "
		notFound(w, r, message, endpoint)
		return
	}

	var methods = make(JObject)
	if ms, ok := registry.methods[route.methods]; ok {
		for verb, m := range ms {
			methods[verb] = struct {
				Headers int `json:"headers,omitempty"`
				Query   int `json:"query,omitempty"`
			}{
				m.headers,
				m.query,
			}
		}
	}

	json.NewEncoder(w).Encode(struct {
		Path       string  `json:"path"`
		UriParams  int     `json:"uriParams,omitempty"`
		Methods    int     `json:"methods,omitempty"`
		Configured JObject `json:"configured"`
	}{
		route.path,
		route.uriParams,
		route.methods,
		methods,
	})
}

func GetSystemMethod(w http.ResponseWriter, r *http.Request) {
	var endpoint = chi.URLParam(r, "endpoint")
	var verb = chi.URLParam(r, "method")

	id, err := isId(endpoint)
	if err != nil {
		util.GetLogger(r).Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	route, ok := registry.endpoints[id]
	if !ok {
		message := "no endpoint found with id "
		notFound(w, r, message, endpoint)
		return
	}

	methods, ok := registry.methods[route.methods]
	if !ok {
		message := "no methods defined for "
		notFound(w, r, message, route.path)
		return
	}

	verb = strings.ToUpper(verb)
	method, ok := methods[verb]
	if !ok {
		message := "not defined for " + route.path
		notFound(w, r, "method", verb, message)
		return
	}

	json.NewEncoder(w).Encode(struct {
		Id      int    `json:"id"`
		Method  string `json:"method"`
		Uri     int    `json:"uriParams"`
		Query   int    `json:"query"`
		Headers int    `json:"headers"`
	}{
		method.id,
		method.name,
		route.uriParams,
		method.query,
		method.headers,
	})
}

func GetSystemParameters(w http.ResponseWriter, _ *http.Request) {
	var display = make(map[int]map[string]string)

	for id, params := range registry.parameters {
		if display[id] == nil {
			display[id] = make(map[string]string)
		}
		for name, p := range params {
			display[id][name] = p.typ
		}
	}

	json.NewEncoder(w).Encode(display)
}

func GetSystemParameterById(w http.ResponseWriter, r *http.Request) {
	parameter := chi.URLParam(r, "parameter")

	id, err := isId(parameter)
	if err != nil {
		util.GetLogger(r).Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	params, ok := registry.parameters[id]
	if !ok {
		message := "no parameter group bears the id"
		notFound(w, r, message, parameter)
		return
	}

	var display = make(map[string]any)
	type props map[string]string
	for name, p := range params {
		display[name] = struct {
			Type       string `json:"type"`
			Required   bool   `json:"required"`
			Properties props  `json:"properties,omitempty"`
		}{
			p.typ,
			p.required,
			registry.properties[p.properties],
		}
	}

	json.NewEncoder(w).Encode(display)
}

func GetSystemProperties(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(registry.properties)
}

func GetSystemPropertyById(w http.ResponseWriter, r *http.Request) {
	property := chi.URLParam(r, "property")

	id, err := isId(property)
	if err != nil {
		util.GetLogger(r).Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	if properties, ok := registry.properties[id]; !ok {
		message := "no property group possesses id"
		notFound(w, r, message, property)
	} else {
		json.NewEncoder(w).Encode(properties)
	}
}
