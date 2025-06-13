package api

import (
	"Factory/internal/util"
	"encoding/json"
	"net/http"
	"strings"
)

/*  **************************
           GET REQUESTS
	************************** */

type GetSystemEndpointsMethod struct {
	Query   int `json:"query,omitempty"`
	Headers int `json:"headers,omitempty"`
}

type GetSystemEndpoints struct {
	Path       string `json:"path"`
	UriParams  int    `json:"uriParams,omitempty"`
	Methods    int    `json:"methods,omitempty"`
	Configured any    `json:"configured"`
}

type GetSystemMethod struct {
	Id      int    `json:"id"`
	Method  string `json:"method"`
	Uri     int    `json:"uriParams"`
	Query   int    `json:"query"`
	Headers int    `json:"headers"`
}

type GetSystemParametersById struct {
	Type       string            `json:"type"`
	Required   bool              `json:"required"`
	Properties map[string]string `json:"properties,omitempty"`
}

/*  **************************
          POST REQUESTS
	************************** */

func SuccessfulSystemPost(w http.ResponseWriter, r *http.Request, args ...string) {
	message := strings.Join(args, " ")
	util.GetLogger(r).Info(message)

	json.NewEncoder(w).Encode(struct {
		Id      string
		Code    int
		Message string
	}{
		Id:      w.Header().Get("X-Correlation-ID"),
		Code:    200,
		Message: message,
	})
}

type PostSystemEndpointRequest struct {
	Path      string
	Methods   int
	UriParams int
}

type PostSystemMethodRequest struct {
	Name    string
	Query   int
	Headers int
}
