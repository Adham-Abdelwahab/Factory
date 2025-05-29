package util

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

func GetLogger(r *http.Request) *logrus.Entry {
	return r.Context().Value("logger").(*logrus.Entry)
}

func SetLogger(r *http.Request, id string) *http.Request {
	entry := logrus.WithField("id", id)
	ctx := context.WithValue(r.Context(), "logger", entry)
	return r.WithContext(ctx)
}
