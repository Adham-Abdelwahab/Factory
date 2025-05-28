package resources

import (
	"encoding/json"
	"net/http"

	"Factory/api"
	"Factory/api/resources"
	"Factory/internal/util"

	"github.com/gorilla/schema"
)

func GetResource(w http.ResponseWriter, r *http.Request) {
	var log = util.GetLogger(r)
	var err error

	var req = resources.ResourceRequest{}
	var d = schema.NewDecoder()

	err = d.Decode(&req, r.URL.Query())
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	if err = json.NewEncoder(w).Encode(req); err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
	}
}
