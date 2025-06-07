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

func getId(name string, w http.ResponseWriter, r *http.Request) (int, bool) {
	if id, err := strconv.Atoi(chi.URLParam(r, name)); err != nil {
		err = errors.New("uri value {" + name + "} must be a number")
		util.GetLogger(r).Error(err)
		api.RequestErrorHandler(w, err)
		return -1, true
	} else {
		return id, false
	}
}

func GetSystemEndpoints(w http.ResponseWriter, r *http.Request) {
	var endpoints = make(map[int]string)
	var basePath = r.URL.Query().Get("basePath")

	for _, e := range registry {
		if strings.HasPrefix(e.endpoint.path, basePath) {
			endpoints[e.endpoint.id] = e.endpoint.path
		}
	}

	json.NewEncoder(w).Encode(endpoints)
}

/*
Must index registry by endpoint.id

Must display endpoint methods and parameters
*/
func GetSystemEndpointById(w http.ResponseWriter, r *http.Request) {
	id, err := getId("endpoint", w, r)
	if err {
		return
	}

	type displayParameters map[string]any
	type displayMethod struct {
		Id        int
		Name      string
		UriParams displayParameters `json:",omitempty"`
		Headers   displayParameters `json:",omitempty"`
		Query     displayParameters `json:",omitempty"`
	}

	var displayParams = func(params []parameter) map[string]any {
		toDisplay := make(map[string]any)
		var required []string

		for _, p := range params {
			toDisplay[p.name] = p.typ
			if p.required {
				required = append(required, p.name)
			}
		}

		if len(required) != 0 {
			toDisplay["required"] = required
		}

		return toDisplay
	}

	var methods []displayMethod
	for _, endpoint := range registry {
		if endpoint.endpoint.id == id {
			for _, m := range endpoint.methods {
				v := endpoint.validator[m.name]
				methods = append(methods, displayMethod{
					Id:        m.id,
					Name:      m.name,
					UriParams: displayParams(v.uriParams),
					Headers:   displayParams(v.headers),
					Query:     displayParams(v.query),
				})
			}
			break
		}
	}

	json.NewEncoder(w).Encode(methods)
}

/*
Must index registry by endpoint.id

Must index endpoint by method
*/
func GetSystemMethodById(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(chi.URLParam(r, "method"))
	id, err := getId("endpoint", w, r)
	if err {
		return
	}

	for _, endpoint := range registry {
		if endpoint.endpoint.id == id {
			for _, m := range endpoint.methods {
				if method == m.name {
					json.NewEncoder(w).Encode(struct {
						MethodId  int
						Method    string
						UriParams int
						Headers   int
						Query     int
					}{
						m.id,
						m.name,
						endpoint.endpoint.uriParams,
						m.headers,
						m.parameters})
					return
				}
			}
		}
	}

	message := errors.New(method + " not found for endpoint " + strconv.Itoa(id))
	util.GetLogger(r).Error(message)
	api.NotFoundErrorHandler(w, message)
}

/* Must store all parameters separately */
func GetSystemParameters(w http.ResponseWriter, r *http.Request) {

}

func GetSystemParameterById(w http.ResponseWriter, r *http.Request) {
	_, err := getId("parameter", w, r)
	if err {
		return
	}
}

/* Must store all properties separately */
func GetSystemProperties(w http.ResponseWriter, r *http.Request) {

}

func GetSystemPropertyById(w http.ResponseWriter, r *http.Request) {
	_, err := getId("property", w, r)
	if err {
		return
	}
}
