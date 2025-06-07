package api

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Id      string
	Code    int
	Message string
}

func raise(w http.ResponseWriter, code int, message string) {
	id := w.Header().Get("X-Correlation-ID")
	w.WriteHeader(code)
	resp := Error{
		Id:      id,
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

var (
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		raise(w, http.StatusBadRequest, err.Error())
	}
	NotFoundErrorHandler = func(w http.ResponseWriter, err error) {
		raise(w, http.StatusNotFound, err.Error())
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		raise(w, http.StatusInternalServerError, "Internal Server Error.")
	}
	ConnectionErrorHandler = func(w http.ResponseWriter, err error) {
		raise(w, 503, err.Error())
	}
)
