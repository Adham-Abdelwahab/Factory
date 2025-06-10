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

func isId(w http.ResponseWriter, r *http.Request, id *int, value string) bool {
	if number, err := strconv.Atoi(value); err != nil {
		message := " is not a valid id"
		err = errors.New(value + message)
		util.GetLogger(r).Error(err)
		api.RequestErrorHandler(w, err)
		return false
	} else {
		*id = number
		return true
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

	var id int
	if !isId(w, r, &id, endpoint) {
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
				Headers int `json:",omitempty"`
				Query   int `json:",omitempty"`
			}{
				m.headers,
				m.query,
			}
		}
	}

	json.NewEncoder(w).Encode(struct {
		Path       string
		UriParams  int `json:",omitempty"`
		Methods    int `json:",omitempty"`
		Configured JObject
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

	var id int
	if !isId(w, r, &id, endpoint) {
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
		Id      int
		Name    string
		Uri     int
		Query   int
		Headers int
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

	var id int
	if !isId(w, r, &id, parameter) {
		return
	}

	params, ok := registry.parameters[id]
	if !ok {
		message := "no parameter group bears the id"
		notFound(w, r, message, parameter)
		return
	}

	var display = make(map[string]any)
	for name, p := range params {
		display[name] = struct {
			Type       string
			Required   bool
			Properties map[string]string `json:",omitempty"`
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

	var id int
	if !isId(w, r, &id, property) {
		return
	}

	props, ok := registry.properties[id]
	if !ok {
		message := "no property group possesses id"
		notFound(w, r, message, property)
		return
	}

	json.NewEncoder(w).Encode(props)
}
