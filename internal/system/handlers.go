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
		message := "%s is not a valid id"
		message = util.Message(message, value)
		return 0, errors.New(message)
	} else {
		return number, nil
	}
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
		api.RequestErrorHandler(w, r, err.Error())
		return
	}

	route, ok := registry.endpoints[id]
	if !ok {
		message := "no endpoint found with id %v"
		message = util.Message(message, id)
		api.RequestErrorHandler(w, r, message)
		return
	}

	var methods = make(JObject)
	if ms, ok := registry.methods[route.methods]; ok {
		for verb, m := range ms {
			methods[verb] = api.GetSystemEndpointsMethod{
				Query:   m.query,
				Headers: m.headers,
			}
		}
	}

	json.NewEncoder(w).Encode(api.GetSystemEndpoints{
		Path:       route.path,
		UriParams:  route.uriParams,
		Methods:    route.methods,
		Configured: methods,
	})
}

func GetSystemMethod(w http.ResponseWriter, r *http.Request) {
	var endpoint = chi.URLParam(r, "endpoint")
	var verb = chi.URLParam(r, "method")

	id, err := isId(endpoint)
	if err != nil {
		api.RequestErrorHandler(w, r, err.Error())
		return
	}

	route, ok := registry.endpoints[id]
	if !ok {
		message := "no endpoint found with id %s"
		message = util.Message(message, endpoint)
		api.NotFoundErrorHandler(w, r, message)
		return
	}

	methods, ok := registry.methods[route.methods]
	if !ok {
		message := "no methods defined for %s"
		message = util.Message(message, route.path)
		api.NotFoundErrorHandler(w, r, message)
		return
	}

	verb = strings.ToUpper(verb)
	method, ok := methods[verb]
	if !ok {
		message := "not defined for %s"
		message = util.Message(message, route.path)
		api.NotFoundErrorHandler(w, r, message)
		return
	}

	json.NewEncoder(w).Encode(api.GetSystemMethod{
		Id:      method.id,
		Method:  method.name,
		Uri:     route.uriParams,
		Query:   method.query,
		Headers: method.headers,
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
		api.RequestErrorHandler(w, r, err.Error())
		return
	}

	params, ok := registry.parameters[id]
	if !ok {
		message := "no parameter group bears the id %s"
		message = util.Message(message, parameter)
		api.NotFoundErrorHandler(w, r, message)
		return
	}

	var display = make(map[string]any)
	for name, p := range params {
		display[name] = api.GetSystemParametersById{
			Type:       p.typ,
			Required:   p.required,
			Properties: registry.properties[p.properties],
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
		api.RequestErrorHandler(w, r, err.Error())
		return
	}

	if properties, ok := registry.properties[id]; !ok {
		message := "no property group possesses id %s"
		message = util.Message(message, property)
		api.NotFoundErrorHandler(w, r, message)
	} else {
		json.NewEncoder(w).Encode(properties)
	}
}

func PostSystemEndpoint(w http.ResponseWriter, r *http.Request) {
	var db = util.Database
	var endpoint string

	json.NewDecoder(r.Body).Decode(&endpoint)

	endpoint = strings.Trim(endpoint, " /")
	endpoint = "/" + endpoint
	if endpoint == "/" {
		message := "endpoint path must be provided"
		api.RequestErrorHandler(w, r, message)
		return
	}

	for _, e := range registry.endpoints {
		if e.path == endpoint {
			message := "%s is already registered"
			message = util.Message(message, e.path)
			api.RequestErrorHandler(w, r, message)
			return
		}
	}

	var id int
	sql := `INSERT INTO endpoint(path, "uriParams", methods)
			VALUES ($1, 0, 0) RETURNING id`
	_ = db.QueryRow(&id, sql, endpoint)

	registry.endpoints[id] = _endpoint{
		id,
		endpoint,
		0,
		0,
	}

	message := "Successfully registered endpoint %s"
	message = util.Message(message, endpoint)
	api.SuccessfulSystemPost(w, r, message)
}

func PostSystemMethod(w http.ResponseWriter, r *http.Request) {
	var db = util.Database

	var method api.PostSystemMethodRequest
	json.NewDecoder(r.Body).Decode(&method)

	method.Name = chi.URLParam(r, "method")
	method.Name = strings.ToUpper(method.Name)

	endpoint := chi.URLParam(r, "endpoint")
	id, err := strconv.Atoi(endpoint)
	if err != nil {
		message := "endpoint id (%s) must be a number"
		message = util.Message(message, endpoint)
		api.RequestErrorHandler(w, r, message)
		return
	}

	e, ok := registry.endpoints[id]
	if !ok {
		message := "endpoint %s does not exist"
		message = util.Message(message, endpoint)
		api.RequestErrorHandler(w, r, message)
		return
	}

	if _, ok = registry.methods[e.methods][method.Name]; ok {
		message := "%s is already registered for %s"
		message = util.Message(message, method.Name, e.path)
		api.RequestErrorHandler(w, r, message)
		return
	}

	if e.methods == 0 {
		sql := `INSERT INTO method (id, name, headers, query) VALUES (DEFAULT, $1, $2, $3) RETURNING id`
		_ = db.QueryRow(&e.methods, sql, method.Name, method.Headers, method.Query)

		sql = `UPDATE endpoint SET methods = $1 WHERE id = $2`
		_ = db.QueryRow(nil, sql, e.methods, e.id)

		registry.endpoints[id] = e
		registry.methods[e.methods] = make(map[string]_method)
	} else {
		sql := `INSERT INTO method (id, name, headers, query) VALUES ($1, $2, $3, $4)`
		_ = db.QueryRow(nil, sql, e.methods, method.Name, method.Headers, method.Query)
	}

	registry.methods[e.methods][method.Name] = _method{
		id:      e.methods,
		name:    method.Name,
		query:   method.Query,
		headers: method.Headers,
	}

	message := "Successfully registered method %s for %s"
	message = util.Message(message, method.Name, e.path)
	api.SuccessfulSystemPost(w, r, message)
}
