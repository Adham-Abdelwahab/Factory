package api

import (
	"encoding/json"
	"net/http"

	"Factory/internal/util"
)

type response struct {
	Id      string
	Code    int
	Message string
}

func raise(w http.ResponseWriter, r *http.Request, code int, err string) {
	id := w.Header().Get("X-Correlation-ID")
	util.GetLogger(r).Error(err)

	w.WriteHeader(code)
	response := response{
		Id:      id,
		Code:    code,
		Message: err,
	}

	json.NewEncoder(w).Encode(response)
}

var (
	RequestErrorHandler = func(w http.ResponseWriter, r *http.Request, err string) {
		raise(w, r, http.StatusBadRequest, err)
	}
	NotFoundErrorHandler = func(w http.ResponseWriter, r *http.Request, err string) {
		raise(w, r, http.StatusNotFound, err)
	}
	ConnectionErrorHandler = func(w http.ResponseWriter, r *http.Request, err string) {
		raise(w, r, http.StatusServiceUnavailable, err)
	}
	InternalErrorHandler = func(w http.ResponseWriter, r *http.Request) {
		raise(w, r, http.StatusInternalServerError, "Internal Server Error")
	}
)
