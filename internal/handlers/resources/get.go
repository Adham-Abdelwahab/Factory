package resources

import (
	"encoding/json"
	"net/http"

	"Factory/api"
	"Factory/api/resources"
	"Factory/internal/util"
)

func GetResource(w http.ResponseWriter, r *http.Request) {
	var log = util.GetLogger(r)
	var err error

	var req = resources.GetResourceRequest{}
	err = util.SafeDecode(&req, r)
	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(req); err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
	}
}
