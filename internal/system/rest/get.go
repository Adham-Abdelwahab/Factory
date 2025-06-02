package rest

import (
	"encoding/json"
	"net/http"
)

func Get(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("get request")
}
