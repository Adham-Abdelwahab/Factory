package util

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func GetLogger(r *http.Request) *logrus.Entry {
	return r.Context().Value("logger").(*logrus.Entry)
}
